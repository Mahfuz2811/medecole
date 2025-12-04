interface AuthHeaderProps {
	formType: "login" | "register";
}

export default function AuthHeader({ formType }: AuthHeaderProps) {
	return (
		<div>
			<h1 className="text-3xl sm:text-4xl font-bold text-gray-900 tracking-tight">
				{formType === "login" ? "Welcome back" : "Get started"}
			</h1>
			<p className="mt-2 text-base sm:text-lg text-gray-600">
				{formType === "login"
					? "Sign in to continue your learning journey"
					: "Create your account and start learning"}
			</p>
		</div>
	);
}
