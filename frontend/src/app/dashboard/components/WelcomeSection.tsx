interface User {
	name?: string;
	msisdn?: string;
}

interface WelcomeSectionProps {
	user: User | null;
}

export function WelcomeSection({ user }: WelcomeSectionProps) {
	return (
		<div className="bg-white rounded-lg shadow-sm p-6 mb-6">
			<div className="flex justify-between items-start">
				<div>
					<h1 className="text-2xl font-bold text-gray-900 mb-2">
						Welcome back, {user?.name}! 👋
					</h1>
					<p className="text-gray-600">
						Ready to continue your learning journey?
					</p>
					<p className="text-sm text-gray-500 mt-1">
						Phone: {user?.msisdn}
					</p>
				</div>
			</div>
		</div>
	);
}
