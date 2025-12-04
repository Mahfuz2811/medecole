"use client";

import { useAuth, withAuth } from "@/contexts/AuthContext";
import { useProfileEnrollments } from "@/hooks/useProfileEnrollments";
import { layouts } from "@/styles/design-tokens";
import Image from "next/image";
import Link from "next/link";

// New organized imports
import { BottomNav, StickyHeader } from "@/components/layout";

function ProfilePage() {
	const { user } = useAuth();
	const {
		enrollments,
		activeEnrollments,
		loading: enrollmentsLoading,
		error: enrollmentsError,
	} = useProfileEnrollments();

	if (!user) {
		return (
			<div className="min-h-screen flex items-center justify-center bg-blue-50">
				<div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
			</div>
		);
	}

	return (
		<main className="bg-blue-50 min-h-screen relative">
			{/* Sticky Header */}
			<div className={layouts.stickyHeader}>
				<div className={`${layouts.container} bg-white shadow`}>
					<StickyHeader />
				</div>
			</div>

			<div className={`${layouts.container} ${layouts.pageContent}`}>
				{/* Profile Header */}
				<div className="bg-gradient-to-r from-blue-500 to-purple-600 p-6 text-white rounded-lg mb-6">
					<div className="flex items-center space-x-4">
						<div className="w-16 h-16 rounded-full overflow-hidden bg-white/20 flex items-center justify-center">
							<Image
								src="/avatar.png"
								alt={user.name}
								width={64}
								height={64}
								className="object-cover w-full h-full"
							/>
						</div>
						<div>
							<h1 className="text-2xl font-bold mb-1">
								{user.name}
							</h1>
							<p className="text-blue-100 text-sm">
								Member since{" "}
								{new Date(user.createdAt).toLocaleDateString(
									"en-US",
									{
										year: "numeric",
										month: "long",
										day: "numeric",
									}
								)}
							</p>
						</div>
					</div>
				</div>

				{/* Subscription Status */}
				<div className="bg-white rounded-lg shadow-sm p-6 mb-6">
					<h2 className="text-lg font-semibold text-gray-900 mb-4">
						Subscription Status
					</h2>

					{enrollmentsLoading ? (
						<div className="flex items-center justify-center py-8">
							<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
						</div>
					) : enrollmentsError ? (
						<div className="bg-red-50 border border-red-200 rounded-lg p-4 text-center">
							<p className="text-red-600 text-sm">
								{enrollmentsError}
							</p>
						</div>
					) : activeEnrollments === 0 ? (
						<div className="bg-gray-50 border border-gray-200 rounded-lg p-4 text-center">
							<div className="w-12 h-12 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-3">
								<svg
									className="w-6 h-6 text-gray-400"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										strokeLinecap="round"
										strokeLinejoin="round"
										strokeWidth={2}
										d="M20 12H4"
									/>
								</svg>
							</div>
							<h3 className="text-sm font-medium text-gray-900 mb-1">
								No Active Subscriptions
							</h3>
							<p className="text-xs text-gray-600 mb-3">
								You don&apos;t have any active package
								subscriptions at the moment.
							</p>
							<Link
								href="/"
								className="inline-flex items-center bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors"
							>
								Browse Packages
							</Link>
						</div>
					) : (
						<div className="space-y-3">
							<div className="flex items-center justify-between mb-4">
								<span className="text-sm font-medium text-gray-600">
									Active Subscriptions
								</span>
								<span className="bg-blue-100 text-blue-800 text-xs font-medium px-2 py-1 rounded-full">
									{activeEnrollments} Active
								</span>
							</div>
							{enrollments.map((enrollment) => (
								<div
									key={enrollment.id}
									className="border border-gray-200 rounded-lg p-4 hover:border-blue-200 transition-colors"
								>
									<div className="flex items-start justify-between">
										<div className="flex-1">
											<h4 className="font-medium text-gray-900 text-sm mb-1">
												{enrollment.package_name}
											</h4>
											<div className="flex items-center space-x-3 text-xs text-gray-500">
												<span
													className={`px-2 py-1 rounded-full text-xs font-medium ${
														enrollment.package_type ===
														"FREE"
															? "bg-green-100 text-green-800"
															: "bg-purple-100 text-purple-800"
													}`}
												>
													{enrollment.package_type}
												</span>
												<span>
													{enrollment.completed_exams}
													/{enrollment.total_exams}{" "}
													exams completed
												</span>
											</div>
											{enrollment.expiry_date && (
												<p className="text-xs text-gray-500 mt-1">
													Expires:{" "}
													{new Date(
														enrollment.expiry_date
													).toLocaleDateString(
														"en-US",
														{
															year: "numeric",
															month: "short",
															day: "numeric",
														}
													)}
												</p>
											)}
										</div>
										<div className="flex flex-col items-end">
											<div className="text-right mb-2">
												<div className="text-xs text-gray-500 mb-1">
													Progress
												</div>
												<div className="w-16 bg-gray-200 rounded-full h-2">
													<div
														className="bg-blue-600 h-2 rounded-full transition-all"
														style={{
															width: `${enrollment.progress}%`,
														}}
													></div>
												</div>
												<div className="text-xs text-gray-600 mt-1">
													{enrollment.progress}%
												</div>
											</div>
										</div>
									</div>
								</div>
							))}
						</div>
					)}
				</div>
			</div>

			<BottomNav />
		</main>
	);
}

export default withAuth(ProfilePage);
