import { packagesApiClient } from "./client";

// Dashboard API types
export interface RecentActivityItem {
	id: number;
	exam_title: string;
	package_name: string;
	date: string;
	score: number;
	total_questions: number;
	is_corrects: number;
	time_taken: string;
	status: string;
}

export interface DashboardSummaryResponse {
	success: boolean;
	message: string;
	data: {
		user_stats: {
			total_attempts: number;
			is_corrects: number;
			accuracy_rate: number;
		};
		recent_activity: RecentActivityItem[];
	};
}

// Dashboard API functions
export const DashboardAPI = {
	// Get dashboard summary (includes both stats and recent activity)
	async getDashboardSummary(): Promise<DashboardSummaryResponse> {
		const response = await packagesApiClient.get("/dashboard/summary");
		return response.data;
	},
};
