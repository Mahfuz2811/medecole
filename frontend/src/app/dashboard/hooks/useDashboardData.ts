import { useAuth } from "@/contexts/AuthContext";
import { DashboardAPI, type RecentActivityItem } from "@/lib/api/dashboard";
import { EnrollmentAPI } from "@/lib/api/enrollment";
import type { DashboardEnrollmentsResponse } from "@/lib/api/enrollment-types";
import { useEffect, useState } from "react";

export interface DashboardEnrollment {
	id: number;
	packageId: number;
	packageName: string;
	packageSlug: string;
	packageType: string;
	progress: number;
	totalExams: number;
	completedExams: number;
	expiryDate?: string;
	status: string;
}

// Transform API data to match component interface
export interface RecentActivity {
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

// User stats interface to match API response
export interface UserStats {
	total_attempts: number;
	is_corrects: number;
	accuracy_rate: number;
}

export function useDashboardData() {
	const { user } = useAuth();

	// State for enrollments
	const [enrollments, setEnrollments] =
		useState<DashboardEnrollmentsResponse | null>(null);
	const [enrollmentsLoading, setEnrollmentsLoading] = useState(true);
	const [enrollmentsError, setEnrollmentsError] = useState<string | null>(
		null
	);

	// State for dashboard summary (user stats + recent activity)
	const [userStats, setUserStats] = useState<UserStats | null>(null);
	const [recentActivity, setRecentActivity] = useState<RecentActivity[]>([]);
	const [dashboardLoading, setDashboardLoading] = useState(true);
	const [dashboardError, setDashboardError] = useState<string | null>(null);

	// Fetch enrollments on component mount
	useEffect(() => {
		const fetchEnrollments = async () => {
			try {
				setEnrollmentsLoading(true);
				setEnrollmentsError(null);
				const data = await EnrollmentAPI.getDashboardEnrollments();
				setEnrollments(data);
			} catch (error) {
				console.error("Failed to fetch enrollments:", error);
				setEnrollmentsError("Failed to load enrollments");
			} finally {
				setEnrollmentsLoading(false);
			}
		};

		// Only fetch if user is authenticated
		if (user) {
			fetchEnrollments();
		}
	}, [user]);

	// Fetch dashboard summary (user stats + recent activity)
	useEffect(() => {
		const fetchDashboardSummary = async () => {
			try {
				setDashboardLoading(true);
				setDashboardError(null);
				const response = await DashboardAPI.getDashboardSummary();

				// Transform API data to match component interface
				const transformedActivity: RecentActivity[] =
					response.data.recent_activity.map(
						(item: RecentActivityItem) => ({
							id: item.id,
							packageName: item.package_name,
							examTitle: item.exam_title,
							date: item.date,
							score: item.score,
							totalQuestions: item.total_questions,
							correctAnswers: item.is_corrects,
							timeTaken: item.time_taken,
							status: item.status,
						})
					);

				setRecentActivity(transformedActivity);
				setUserStats(response.data.user_stats);
			} catch (error) {
				console.error("Failed to fetch dashboard summary:", error);
				setDashboardError("Failed to load dashboard data");
			} finally {
				setDashboardLoading(false);
			}
		};

		// Only fetch if user is authenticated
		if (user) {
			fetchDashboardSummary();
		}
	}, [user]);

	// Transform enrollment data for UI compatibility
	const myEnrollments: DashboardEnrollment[] =
		enrollments?.enrollments.map((enrollment) => ({
			id: enrollment.id,
			packageId: enrollment.package_id,
			packageName: enrollment.package_name,
			packageSlug: enrollment.package_slug,
			packageType: enrollment.package_type,
			progress: enrollment.progress,
			totalExams: enrollment.total_exams,
			completedExams: enrollment.completed_exams,
			expiryDate: enrollment.expiry_date,
			status: enrollment.status,
		})) || [];

	return {
		user,
		myEnrollments,
		enrollmentsLoading,
		enrollmentsError,
		recentActivity,
		userStats,
		dashboardLoading,
		dashboardError,
	};
}
