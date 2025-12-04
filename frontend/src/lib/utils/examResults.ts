import type {
	ExamResultOption,
	ExamResultQuestion,
	RawAnswersData,
	RawQuestionData,
} from "@/lib/api/types";

/**
 * Transform raw backend question data to format expected by UI components
 */
export function transformRawQuestionData(
	rawQuestion: RawQuestionData
): ExamResultQuestion {
	// Transform options from backend format to UI format
	const options: ExamResultOption[] = rawQuestion.options.map(
		(option, index) => ({
			id: index, // Use index as id since backend doesn't provide option IDs
			option_text: option.text,
			is_correct: option.is_correct,
		})
	);

	// Transform correct_answer from backend format to UI format
	let correct_answer: number | boolean[];
	if (rawQuestion.question_type === "SBA") {
		// For SBA, correct_answer is a string (e.g., "c")
		const correctKey = rawQuestion.correct_answer as string;
		correct_answer = rawQuestion.options.findIndex(
			(option) => option.key === correctKey
		);
	} else {
		// For TRUE_FALSE, correct_answer is an object
		const correctAnswerObj = rawQuestion.correct_answer as Record<
			string,
			boolean
		>;
		correct_answer = rawQuestion.options.map(
			(option) => correctAnswerObj[option.key] || false
		);
	}

	// Transform user_answer from backend format to UI format
	let user_answer: number | boolean[] | (boolean | undefined)[] | null = null;
	if (rawQuestion.user_answer && rawQuestion.user_answer.length > 0) {
		if (rawQuestion.question_type === "SBA") {
			// For SBA, user_answer is a single key (e.g., ["c"])
			const selectedKey = rawQuestion.user_answer[0];
			if (selectedKey) {
				user_answer = rawQuestion.options.findIndex(
					(option) => option.key === selectedKey
				);
			}
		} else {
			// For TRUE_FALSE, parse the user answers
			const userAnswerMap: Record<string, boolean | undefined> = {};
			rawQuestion.user_answer.forEach((answer) => {
				const [key, value] = answer.split(":");
				userAnswerMap[key] = value === "true";
			});
			user_answer = rawQuestion.options.map(
				(option) => userAnswerMap[option.key] // undefined for unanswered, boolean for answered
			);
		}
	}

	return {
		id: rawQuestion.question_id,
		question: rawQuestion.question_text,
		question_type: rawQuestion.question_type,
		options,
		correct_answer,
		user_answer,
		is_correct: rawQuestion.is_correct,
		points: rawQuestion.points_earned,
		max_points: rawQuestion.max_points,
		explanation: rawQuestion.explanation,
	};
}

/**
 * Transform raw backend data to format expected by UI components
 */
export function transformRawAnswersData(rawData: RawAnswersData) {
	// Use backend-provided data from exam_snapshot instead of calculating
	const examSnapshot = rawData.exam_snapshot;

	// Fallback calculations if backend data is not available (for backwards compatibility)
	const totalQuestionsFromAnswers = rawData.answers.length;
	const correctAnswersFromAnswers = rawData.answers.filter(
		(q) => q.is_correct
	).length;
	const totalScoreFromAnswers = rawData.answers.reduce(
		(sum, q) => sum + q.points_earned,
		0
	);
	const maxScoreFromAnswers = rawData.answers.reduce(
		(sum, q) => sum + q.max_points,
		0
	);
	const scorePercentageFromAnswers =
		maxScoreFromAnswers > 0
			? (totalScoreFromAnswers / maxScoreFromAnswers) * 100
			: 0;

	// Use backend data when available, otherwise fall back to calculated values
	const totalQuestions =
		examSnapshot.total_questions || totalQuestionsFromAnswers;
	const correctAnswers =
		examSnapshot.correct_answers ?? correctAnswersFromAnswers;
	const totalScore = examSnapshot.score ?? totalScoreFromAnswers;
	const maxScore = maxScoreFromAnswers; // This still needs to be calculated from questions
	const scorePercentage =
		examSnapshot.score !== null && examSnapshot.score !== undefined
			? maxScore > 0
				? (totalScore / maxScore) * 100
				: 0
			: scorePercentageFromAnswers;
	const isPassed =
		examSnapshot.is_passed ?? totalScore >= examSnapshot.passing_score;
	const timeSpent =
		examSnapshot.actual_time_spent || examSnapshot.duration_minutes * 60;

	return {
		exam: {
			id: 1, // Default value since not provided
			title: "Exam Results", // Default value since not provided
			description: "Exam Results", // Default value since not provided
			total_questions: totalQuestions,
			duration_minutes: examSnapshot.duration_minutes,
			passing_score: examSnapshot.passing_score,
			exam_type: "EXAM", // Default value since not provided
		},
		attempt: {
			id: 1, // Default value since not provided
			status: "COMPLETED", // Default value since not provided
			started_at: new Date(Date.now() - timeSpent * 1000).toISOString(), // Estimate based on actual time spent
			completed_at: rawData.submission_timestamp,
			score: totalScore,
			max_score: maxScore,
			score_percentage: scorePercentage,
			correct_answers: correctAnswers,
			wrong_answers: totalQuestions - correctAnswers,
			is_passed: isPassed,
			time_spent: timeSpent,
		},
		package: {
			name: "Package Name", // Default value since not provided
			slug: "package-slug", // Default value since not provided
		},
		questions: rawData.answers.map(transformRawQuestionData),
	};
}
