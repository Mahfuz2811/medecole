"use client";

import {
	AuthHeader,
	BackToHome,
	ErrorMessage,
	FormToggle,
	LoginForm,
	RegisterForm,
	SocialLoginOptions,
	useAuthRedirect,
	useFormType,
	type LoginFormData,
	type RegisterFormData,
} from "@/components/auth";
import { useAuth } from "@/contexts/AuthContext";
import { Suspense, useState } from "react";

function AuthPageContent() {
	const {
		login,
		register,
		isAuthenticated,
		isLoading: authLoading,
	} = useAuth();
	const { formType, setFormType } = useFormType();

	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string>("");

	// Handle redirect when authenticated
	useAuthRedirect({ isAuthenticated });

	// Show loading while auth context is initializing
	if (authLoading) {
		return (
			<div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-purple-50">
				<div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
			</div>
		);
	}

	const handleLogin = async (data: LoginFormData) => {
		setIsLoading(true);
		setError("");

		try {
			await login(data.msisdn, data.password);
			// Redirect will happen via useAuthRedirect hook
		} catch (err) {
			setError(err instanceof Error ? err.message : "Login failed");
		} finally {
			setIsLoading(false);
		}
	};

	const handleRegister = async (data: RegisterFormData) => {
		setIsLoading(true);
		setError("");

		try {
			await register(data.name, data.msisdn, data.password);
			// Redirect will happen via useAuthRedirect hook
		} catch (err) {
			setError(
				err instanceof Error ? err.message : "Registration failed"
			);
		} finally {
			setIsLoading(false);
		}
	};

	const enableTraditionalAuth =
		process.env.NEXT_PUBLIC_ENABLE_TRADITIONAL_AUTH === "true";

	return (
		<main className="min-h-screen bg-gradient-to-br from-indigo-50 via-purple-50 to-pink-50 relative overflow-hidden">
			{/* Decorative background elements */}
			<div className="absolute inset-0 overflow-hidden pointer-events-none">
				<div className="absolute -top-40 -right-40 w-80 h-80 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob"></div>
				<div className="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-2000"></div>
				<div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-pink-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-4000"></div>
			</div>

			<div className="relative flex items-center justify-center min-h-screen px-4 py-8 sm:px-6 lg:px-8">
				<div className="w-full max-w-md space-y-6">
					{/* Logo and Header */}
					<div className="text-center">
						<div className="flex justify-center mb-6">
							<div className="w-16 h-16 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg">
								<svg
									className="w-10 h-10 text-white"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
								>
									<path
										strokeLinecap="round"
										strokeLinejoin="round"
										strokeWidth={2}
										d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
									/>
								</svg>
							</div>
						</div>

						{enableTraditionalAuth ? (
							<AuthHeader formType={formType} />
						) : (
							<>
								<h1 className="text-4xl sm:text-5xl font-extrabold text-gray-900 tracking-tight">
									Welcome to
									<span className="block text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600 mt-1">
										Medecole
									</span>
								</h1>
								<p className="mt-3 text-base sm:text-lg text-gray-600 max-w-sm mx-auto">
									Sign in to access your exams and track your
									learning progress
								</p>
							</>
						)}
					</div>

					{/* Auth Card */}
					<div className="bg-white/80 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-6 sm:p-8 space-y-6">
						<ErrorMessage error={error} />

						{enableTraditionalAuth ? (
							<>
								<FormToggle
									formType={formType}
									onToggle={setFormType}
									onError={setError}
								/>

								{formType === "login" ? (
									<LoginForm
										isLoading={isLoading}
										onSubmit={handleLogin}
										onError={setError}
									/>
								) : (
									<RegisterForm
										isLoading={isLoading}
										onSubmit={handleRegister}
										onError={setError}
									/>
								)}
							</>
						) : null}

						<SocialLoginOptions />
					</div>

					<BackToHome />
				</div>
			</div>
		</main>
	);
}

export default function AuthPage() {
	return (
		<Suspense
			fallback={
				<div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-purple-50">
					<div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
				</div>
			}
		>
			<AuthPageContent />
		</Suspense>
	);
}
