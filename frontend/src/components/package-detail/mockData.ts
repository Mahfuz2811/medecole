import { ExamType } from "@/lib/api/types";

export interface MockExam {
	id: number;
	title: string;
	slug: string;
	description: string;
	exam_type: ExamType;
	total_questions: number;
	duration_minutes: number;
	passing_score: number;
	max_attempts: number;
	scheduled_start_date: string;
	scheduled_end_date: string;
	status: "SCHEDULED";
	is_active: boolean;
	sort_order: number;
	created_at: string;
	updated_at: string;
	// Analytics data
	attempt_count: number;
	completed_attempt_count: number;
	average_score: number;
	pass_rate: number;
	last_attempt_at: string | null;
}

export const mockExams: MockExam[] = [
	{
		id: 1,
		title: "Daily Practice Test #1",
		slug: "daily-practice-1",
		description: "Foundation level questions covering basic concepts",
		exam_type: "DAILY" as ExamType,
		total_questions: 50,
		duration_minutes: 60,
		passing_score: 60,
		max_attempts: 3,
		scheduled_start_date: "2025-08-01T09:00:00Z",
		scheduled_end_date: "2025-08-01T23:59:59Z",
		status: "SCHEDULED",
		is_active: true,
		sort_order: 1,
		created_at: "2025-07-30T00:00:00Z",
		updated_at: "2025-07-30T00:00:00Z",
		attempt_count: 1247,
		completed_attempt_count: 1158,
		average_score: 78.5,
		pass_rate: 85.2,
		last_attempt_at: "2025-07-30T10:30:00Z",
	},
	{
		id: 3,
		title: "Mock Examination #1",
		slug: "mock-exam-1",
		description:
			"Full-length mock examination simulating real exam conditions",
		exam_type: "MOCK" as ExamType,
		total_questions: 200,
		duration_minutes: 180,
		passing_score: 75,
		max_attempts: 1,
		scheduled_start_date: "2025-08-15T09:00:00Z",
		scheduled_end_date: "2025-08-15T14:00:00Z",
		status: "SCHEDULED",
		is_active: true,
		sort_order: 3,
		created_at: "2025-07-30T00:00:00Z",
		updated_at: "2025-07-30T00:00:00Z",
		attempt_count: 567,
		completed_attempt_count: 521,
		average_score: 68.9,
		pass_rate: 71.4,
		last_attempt_at: "2025-07-30T08:45:00Z",
	},
	{
		id: 4,
		title: "Final Assessment",
		slug: "final-assessment",
		description:
			"Comprehensive final assessment covering all course materials",
		exam_type: "FINAL" as ExamType,
		total_questions: 150,
		duration_minutes: 150,
		passing_score: 80,
		max_attempts: 1,
		scheduled_start_date: "2025-07-25T09:00:00Z", // Available now
		scheduled_end_date: "2025-08-30T12:30:00Z",
		status: "SCHEDULED",
		is_active: true,
		sort_order: 4,
		created_at: "2025-07-30T00:00:00Z",
		updated_at: "2025-07-30T00:00:00Z",
		attempt_count: 0,
		completed_attempt_count: 0,
		average_score: 0,
		pass_rate: 0,
		last_attempt_at: null,
	},
	{
		id: 5,
		title: "Advanced Review Test",
		slug: "advanced-review",
		description: "Advanced level review test for experienced candidates",
		exam_type: "REVIEW" as ExamType,
		total_questions: 75,
		duration_minutes: 90,
		passing_score: 70,
		max_attempts: 3,
		scheduled_start_date: "2025-08-15T09:00:00Z", // Future exam
		scheduled_end_date: "2025-08-15T12:30:00Z",
		status: "SCHEDULED",
		is_active: true,
		sort_order: 5,
		created_at: "2025-07-30T00:00:00Z",
		updated_at: "2025-07-30T00:00:00Z",
		attempt_count: 0,
		completed_attempt_count: 0,
		average_score: 0,
		pass_rate: 0,
		last_attempt_at: null,
	},
];

// Helper functions for exam availability and analytics
export const isExamAvailable = (exam: MockExam) => {
	const now = new Date();
	const startDate = new Date(exam.scheduled_start_date);
	return now >= startDate && exam.status === "SCHEDULED";
};

export const shouldShowAnalytics = (exam: MockExam) => {
	if (!isExamAvailable(exam)) {
		return false;
	}
	return true;
};

export const hasParticipation = (exam: MockExam) => {
	return exam.attempt_count > 0;
};
