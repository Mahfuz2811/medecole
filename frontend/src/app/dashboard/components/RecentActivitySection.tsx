import { designTokens } from "@/styles/design-tokens";
import Link from "next/link";

interface RecentActivity {
	id: number;
	packageName: string;
	examTitle: string;
	date: string;
	score: number;
	totalQuestions: number;
	correctAnswers: number;
	timeTaken: string;
	status: string;
}

interface RecentActivitySectionProps {
	activities: RecentActivity[];
	loading?: boolean;
	error?: string | null;
}

export function RecentActivitySection({
	activities,
	loading = false,
	error = null,
}: RecentActivitySectionProps) {
	return (
		<div className="bg-white rounded-lg shadow-sm p-6 mb-6">
			<div className="flex justify-between items-center mb-4">
				<h2 className="text-lg font-semibold text-gray-900">
					Recent Activity
				</h2>
				<span className="text-sm text-gray-500">Exam History</span>
			</div>

			{loading ? (
				<div className="space-y-4">
					{/* Loading skeleton */}
					{[1, 2, 3].map((i) => (
						<div
							key={i}
							className="flex items-center justify-between p-4 border border-gray-100 rounded-lg animate-pulse"
						>
							<div className="flex items-center flex-1">
								<div className="w-10 h-10 bg-gray-200 rounded-lg mr-4"></div>
								<div className="flex-1">
									<div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
									<div className="h-3 bg-gray-200 rounded w-1/2"></div>
								</div>
							</div>
							<div className="text-right ml-4">
								<div className="h-5 bg-gray-200 rounded w-12 mb-1"></div>
								<div className="h-3 bg-gray-200 rounded w-16"></div>
							</div>
						</div>
					))}
				</div>
			) : error ? (
				<div className="text-center py-8">
					<div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg
							className="w-8 h-8 text-red-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
					</div>
					<h3 className="text-lg font-medium text-gray-900 mb-2">
						Failed to load activity
					</h3>
					<p className="text-gray-600 mb-4">{error}</p>
					<button
						onClick={() => window.location.reload()}
						className="text-blue-600 hover:text-blue-800 font-medium"
					>
						Try again
					</button>
				</div>
			) : activities.length > 0 ? (
				<div className="space-y-4">
					{activities.map((activity) => (
						<div
							key={activity.id}
							className="flex items-center justify-between p-4 border border-gray-100 rounded-lg"
						>
							<div className="flex items-center flex-1">
								<div
									className={`w-10 h-10 rounded-lg flex items-center justify-center mr-4 ${
										activity.score >= 80
											? "bg-green-100"
											: activity.score >= 60
											? "bg-yellow-100"
											: "bg-red-100"
									}`}
								>
									<svg
										className={`w-5 h-5 ${
											activity.score >= 80
												? "text-green-600"
												: activity.score >= 60
												? "text-yellow-600"
												: "text-red-600"
										}`}
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											strokeLinecap="round"
											strokeLinejoin="round"
											strokeWidth={2}
											d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
										/>
									</svg>
								</div>
								<div className="flex-1">
									<div className="flex items-center gap-2 mb-1">
										<h4 className="font-medium text-gray-900">
											{activity.examTitle}
										</h4>
									</div>
									<div className="flex items-center gap-4 text-sm text-gray-600">
										<span>{activity.date}</span>
										<span className="text-gray-400">|</span>
										<span>{activity.timeTaken}</span>
									</div>
								</div>
							</div>
							<div className="text-right ml-4">
								<div>
									<div className="text-lg font-bold text-gray-900 mb-1">
										{activity.score}
									</div>
									<div className="text-xs text-gray-500">
										{activity.correctAnswers}/
										{activity.totalQuestions} correct
									</div>
								</div>
							</div>
						</div>
					))}
				</div>
			) : (
				<div className="text-center py-8">
					<div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg
							className="w-8 h-8 text-gray-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
							/>
						</svg>
					</div>
					<h3 className="text-lg font-medium text-gray-900 mb-2">
						No activity yet
					</h3>
					<p className="text-gray-600 mb-4">
						Start practicing to see your activity here
					</p>
					<Link
						href="/packages"
						className={`${designTokens.components.button.primary} inline-flex`}
					>
						Start Practicing
					</Link>
				</div>
			)}
		</div>
	);
}
