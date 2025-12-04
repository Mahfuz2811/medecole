import { ExamResponse } from "@/lib/api/types";

export const isExamAvailable = (exam: ExamResponse) => {
	if (!exam.scheduled_start_date) {
		return false;
	}
	const now = new Date();
	const startDate = new Date(exam.scheduled_start_date);
	// Check if exam is currently available (either SCHEDULED or ACTIVE status)
	return (
		now >= startDate &&
		(exam.computed_status === "AVAILABLE" ||
			exam.computed_status === "LIVE")
	);
};

export const shouldShowAnalytics = (exam: ExamResponse) => {
	// Show analytics if the exam has participation data OR if it's currently available
	return hasParticipation(exam) || isExamAvailable(exam);
};

export const hasParticipation = (exam: ExamResponse) => {
	return exam.attempt_count > 0;
};
