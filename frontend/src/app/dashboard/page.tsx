"use client";

import { withAuth } from "@/contexts/AuthContext";
import { layouts } from "@/styles/design-tokens";

// New organized imports
import { BottomNav, StickyHeader } from "@/components/layout";
import {
	MyEnrollments,
	QuickStats,
	RecentActivitySection,
	WelcomeSection,
	useDashboardData,
} from ".";

function DashboardPage() {
	const {
		user,
		myEnrollments,
		enrollmentsLoading,
		enrollmentsError,
		recentActivity,
		userStats,
		dashboardLoading,
		dashboardError,
	} = useDashboardData();

	// Transform userStats to quickStatsData format
	const quickStatsData = userStats
		? [
				{
					label: "Total Attempts",
					value: userStats.total_attempts.toString(),
					color: "bg-blue-500",
					trend: "+3 this week",
				},
				{
					label: "Correct Answers",
					value: userStats.is_corrects.toString(),
					color: "bg-green-500",
					trend: "+15 this week",
				},
				{
					label: "Accuracy",
					value: `${userStats.accuracy_rate.toFixed(1)}%`,
					color: "bg-yellow-500",
					trend: "+2% improvement",
				},
		  ]
		: [
				{
					label: "Total Attempts",
					value: "0",
					color: "bg-blue-500",
					trend: "Start practicing!",
				},
				{
					label: "Correct Answers",
					value: "0",
					color: "bg-green-500",
					trend: "Start practicing!",
				},
				{
					label: "Accuracy",
					value: "0%",
					color: "bg-yellow-500",
					trend: "Start practicing!",
				},
		  ];

	return (
		<main className="bg-blue-50 min-h-screen relative">
			{/* Sticky Header */}
			<div className={layouts.stickyHeader}>
				<div className={`${layouts.container} bg-white shadow`}>
					<StickyHeader />
				</div>
			</div>

			<div className={`${layouts.container} ${layouts.pageContent}`}>
				{/* Welcome Section */}
				<WelcomeSection user={user} />

				{/* Quick Stats */}
				<QuickStats stats={quickStatsData} />

				{/* My Enrollments */}
				<MyEnrollments
					enrollments={myEnrollments}
					loading={enrollmentsLoading}
					error={enrollmentsError}
				/>

				{/* Recent Activity */}
				<RecentActivitySection
					activities={recentActivity}
					loading={dashboardLoading}
					error={dashboardError}
				/>
			</div>

			{/* Blue separator above BottomNav */}
			<div className={`${layouts.container} h-3 bg-blue-100`} />

			<BottomNav />
		</main>
	);
}
export default withAuth(DashboardPage);
