"use client";

import { authAPI } from "@/lib/api/auth";
import { auth } from "@/lib/api/utils";
import FacebookLogin from "@greatsumini/react-facebook-login";
import { GoogleLogin, GoogleOAuthProvider } from "@react-oauth/google";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function SocialLoginOptions() {
	const router = useRouter();
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string>("");

	const googleClientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || "";
	const facebookAppId = process.env.NEXT_PUBLIC_FACEBOOK_APP_ID || "";

	const enableTraditionalAuth =
		process.env.NEXT_PUBLIC_ENABLE_TRADITIONAL_AUTH === "true";

	const handleGoogleSuccess = async (
		credentialResponse: { credential?: string } | undefined
	) => {
		if (!credentialResponse?.credential) {
			setError("No credential received from Google");
			return;
		}

		setIsLoading(true);
		setError("");

		try {
			const authResponse = await authAPI.googleAuth(
				credentialResponse.credential
			);
			auth.setAuthData(authResponse);
			router.push("/dashboard");
		} catch (err) {
			setError(
				err instanceof Error
					? err.message
					: "Google authentication failed"
			);
		} finally {
			setIsLoading(false);
		}
	};

	const handleGoogleError = () => {
		setError("Google sign-in was cancelled or failed");
	};

	const handleFacebookSuccess = async (response: {
		accessToken?: string;
	}) => {
		if (!response.accessToken) {
			setError("No access token received from Facebook");
			return;
		}

		setIsLoading(true);
		setError("");

		try {
			const authResponse = await authAPI.facebookAuth(
				response.accessToken
			);
			auth.setAuthData(authResponse);
			router.push("/dashboard");
		} catch (err) {
			setError(
				err instanceof Error
					? err.message
					: "Facebook authentication failed"
			);
		} finally {
			setIsLoading(false);
		}
	};

	const handleFacebookError = (error: unknown) => {
		console.error("Facebook login error:", error);
		setError("Facebook sign-in was cancelled or failed");
	};

	if (!googleClientId && !facebookAppId) {
		return (
			<div className="mt-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
				<p className="text-sm text-yellow-800">
					OAuth credentials not configured. Please set
					NEXT_PUBLIC_GOOGLE_CLIENT_ID and/or
					NEXT_PUBLIC_FACEBOOK_APP_ID in your environment variables.
				</p>
			</div>
		);
	}

	return (
		<div className={enableTraditionalAuth ? "" : ""}>
			{enableTraditionalAuth && (
				<div className="relative my-6">
					<div className="absolute inset-0 flex items-center">
						<div className="w-full border-t border-gray-200" />
					</div>
					<div className="relative flex justify-center text-sm">
						<span className="px-4 bg-white text-gray-500 font-medium">
							Or continue with
						</span>
					</div>
				</div>
			)}

			{!enableTraditionalAuth && (
				<div className="text-center mb-6">
					<p className="text-sm text-gray-600 font-medium">
						sign-in with Google
					</p>
				</div>
			)}

			{error && (
				<div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-xl flex items-start gap-3">
					<svg
						className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5"
						fill="currentColor"
						viewBox="0 0 20 20"
					>
						<path
							fillRule="evenodd"
							d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
							clipRule="evenodd"
						/>
					</svg>
					<p className="text-sm text-red-800 flex-1">{error}</p>
				</div>
			)}

			<div className="space-y-3">
				{googleClientId && (
					<GoogleOAuthProvider clientId={googleClientId}>
						<div className="relative">
							{isLoading && (
								<div className="absolute inset-0 bg-white/50 backdrop-blur-[2px] flex items-center justify-center z-10 rounded-xl">
									<div className="animate-spin rounded-full h-5 w-5 border-b-2 border-indigo-600"></div>
								</div>
							)}
							<div className="[&_div[role=button]]:!w-full [&_div[role=button]]:!rounded-xl [&_div[role=button]]:!shadow-sm [&_div[role=button]]:!border-gray-300 [&_div[role=button]]:hover:!bg-gray-50 [&_div[role=button]]:!transition-all [&_div[role=button]]:!duration-200">
								<GoogleLogin
									onSuccess={handleGoogleSuccess}
									onError={handleGoogleError}
									useOneTap={false}
									theme="outline"
									size="large"
									text="continue_with"
									shape="rectangular"
									logo_alignment="center"
									width="100%"
								/>
							</div>
						</div>
					</GoogleOAuthProvider>
				)}

				{facebookAppId && (
					<div className="relative">
						{isLoading && (
							<div className="absolute inset-0 bg-white/50 backdrop-blur-[2px] flex items-center justify-center z-10 rounded-xl">
								<div className="animate-spin rounded-full h-5 w-5 border-b-2 border-indigo-600"></div>
							</div>
						)}
						<FacebookLogin
							appId={facebookAppId}
							onSuccess={handleFacebookSuccess}
							onFail={handleFacebookError}
							onProfileSuccess={(response: unknown) => {
								console.log("Facebook profile:", response);
							}}
							className="w-full inline-flex justify-center items-center h-[40px] border border-gray-300 rounded-xl shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 hover:border-gray-400 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
							style={{
								backgroundColor: "#fff",
								cursor: isLoading ? "not-allowed" : "pointer",
								pointerEvents: isLoading ? "none" : "auto",
								paddingLeft: "12px",
								paddingRight: "12px",
							}}
						>
							<div className="flex items-center gap-2">
								<svg
									className="w-[18px] h-[18px] text-[#1877F2] flex-shrink-0"
									fill="currentColor"
									viewBox="0 0 24 24"
								>
									<path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z" />
								</svg>
								<span className="leading-none">
									Continue with Facebook
								</span>
							</div>
						</FacebookLogin>
					</div>
				)}
			</div>
		</div>
	);
}
