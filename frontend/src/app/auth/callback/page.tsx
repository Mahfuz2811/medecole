"use client";

import { authAPI } from "@/lib/api/auth";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";

function OAuthCallbackContent() {
	const router = useRouter();
	const searchParams = useSearchParams();
	const [error, setError] = useState<string>("");
	const [isProcessing, setIsProcessing] = useState(true);

	useEffect(() => {
		const processCallback = async () => {
			try {
				const token = searchParams.get("token");
				const errorParam = searchParams.get("error");

				if (errorParam) {
					// Handle error from backend OAuth flow
					const errorMessages: Record<string, string> = {
						missing_parameters:
							"Authentication failed: Missing parameters",
						invalid_state:
							"Authentication failed: Invalid or expired request",
						google_auth_failed:
							"Google authentication failed. Please try again.",
						facebook_auth_failed:
							"Facebook authentication failed. Please try again.",
						"email already registered with different provider":
							"This email is already registered with a different login method. Please use your original sign-in method.",
					};

					setError(
						errorMessages[errorParam] ||
							`Authentication failed: ${errorParam}`
					);
					setIsProcessing(false);

					// Redirect back to auth page after 3 seconds
					setTimeout(() => {
						router.push("/auth");
					}, 3000);
					return;
				}

				if (token) {
					// Store token from server-side OAuth flow
					localStorage.setItem("authToken", token);

					// Fetch user profile
					try {
						const user = await authAPI.getProfile();
						localStorage.setItem("user", JSON.stringify(user));

						// Set cache timestamp
						localStorage.setItem(
							"profile_cache_timestamp",
							Date.now().toString()
						);

						// Redirect to dashboard
						router.push("/dashboard");
					} catch (err) {
						console.error("Failed to fetch profile:", err);
						setError(
							"Authentication successful but failed to load profile. Redirecting..."
						);

						// Still redirect to dashboard, the profile will be fetched there
						setTimeout(() => {
							router.push("/dashboard");
						}, 2000);
					}
				} else {
					// No token or error - something went wrong
					setError(
						"Authentication failed: No response received. Please try again."
					);
					setTimeout(() => {
						router.push("/auth");
					}, 3000);
				}
			} catch (err) {
				console.error("OAuth callback error:", err);
				setError("An unexpected error occurred. Redirecting...");
				setTimeout(() => {
					router.push("/auth");
				}, 3000);
			} finally {
				setIsProcessing(false);
			}
		};

		processCallback();
	}, [searchParams, router]);

	return (
		<div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-purple-50">
			<div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
				{isProcessing ? (
					<div className="text-center">
						<div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-600 mx-auto mb-4"></div>
						<h2 className="text-xl font-semibold text-gray-900 mb-2">
							Processing authentication...
						</h2>
						<p className="text-gray-600">Please wait a moment</p>
					</div>
				) : error ? (
					<div className="text-center">
						<div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-red-100 mb-4">
							<svg
								className="h-8 w-8 text-red-600"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
						</div>
						<h2 className="text-xl font-semibold text-gray-900 mb-2">
							Authentication Failed
						</h2>
						<p className="text-red-600 mb-4">{error}</p>
						<p className="text-sm text-gray-500">
							Redirecting you back to the login page...
						</p>
					</div>
				) : (
					<div className="text-center">
						<div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-green-100 mb-4">
							<svg
								className="h-8 w-8 text-green-600"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M5 13l4 4L19 7"
								/>
							</svg>
						</div>
						<h2 className="text-xl font-semibold text-gray-900 mb-2">
							Authentication Successful!
						</h2>
						<p className="text-gray-600">
							Redirecting you to your dashboard...
						</p>
					</div>
				)}
			</div>
		</div>
	);
}

export default function OAuthCallback() {
	return (
		<Suspense
			fallback={
				<div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-purple-50">
					<div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
						<div className="text-center">
							<div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-600 mx-auto mb-4"></div>
							<h2 className="text-xl font-semibold text-gray-900 mb-2">
								Processing authentication...
							</h2>
							<p className="text-gray-600">
								Please wait a moment
							</p>
						</div>
					</div>
				</div>
			}
		>
			<OAuthCallbackContent />
		</Suspense>
	);
}
