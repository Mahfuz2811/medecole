import {
	Award,
	BookOpen,
	CheckCircle,
	Clock,
	TrendingUp,
	XCircle,
} from "lucide-react";
import { ExamData, Question, QuestionOption, UserAnswer } from "../types";

interface ExamResultsProps {
	exam: ExamData;
	userAnswers: Record<number, UserAnswer>;
	score: number;
	timeSpent: number;
	onRetakeExam: () => void;
	onViewSolutions?: () => void;
}

export function ExamResults({
	exam,
	userAnswers,
	score,
	timeSpent,
	onRetakeExam,
	onViewSolutions,
}: ExamResultsProps) {
	const totalQuestions = exam.total_questions;
	const correctAnswers = Math.round((score / 100) * totalQuestions);
	const incorrectAnswers = totalQuestions - correctAnswers;
	const percentage = Math.round(score);

	// Format time spent
	const formatTime = (seconds: number) => {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		const secs = seconds % 60;

		if (hours > 0) {
			return `${hours}h ${minutes}m ${secs}s`;
		} else if (minutes > 0) {
			return `${minutes}m ${secs}s`;
		} else {
			return `${secs}s`;
		}
	};

	// Get grade based on score
	const getGrade = (percentage: number) => {
		if (percentage >= 90)
			return {
				grade: "A+",
				color: "text-green-600",
				bgColor: "bg-green-50",
			};
		if (percentage >= 80)
			return {
				grade: "A",
				color: "text-green-600",
				bgColor: "bg-green-50",
			};
		if (percentage >= 70)
			return {
				grade: "B",
				color: "text-blue-600",
				bgColor: "bg-blue-50",
			};
		if (percentage >= 60)
			return {
				grade: "C",
				color: "text-yellow-600",
				bgColor: "bg-yellow-50",
			};
		if (percentage >= 50)
			return {
				grade: "D",
				color: "text-orange-600",
				bgColor: "bg-orange-50",
			};
		return { grade: "F", color: "text-red-600", bgColor: "bg-red-50" };
	};

	const gradeInfo = getGrade(percentage);

	// Get performance message
	const getPerformanceMessage = (percentage: number) => {
		if (percentage >= 90)
			return "Excellent work! You&apos;ve demonstrated outstanding knowledge.";
		if (percentage >= 80)
			return "Great job! You have a solid understanding of the material.";
		if (percentage >= 70)
			return "Good work! You&apos;re on the right track.";
		if (percentage >= 60)
			return "Fair performance. Consider reviewing the material again.";
		if (percentage >= 50)
			return "You&apos;re getting there. More practice will help improve your score.";
		return "Keep studying! Consider reviewing the fundamentals and retaking the exam.";
	};

	// Helper function to get correct answer for a question
	const getCorrectAnswer = (question: Question): string => {
		if (question.question_type === "SBA") {
			// Find the option marked as correct
			const correctOption = Object.entries(question.options).find(
				([, option]) => (option as QuestionOption).is_correct
			);
			return correctOption ? correctOption[0] : "";
		} else {
			// For TRUE_FALSE, find the option marked as correct
			const correctOption = Object.entries(question.options).find(
				([, option]) => (option as QuestionOption).is_correct
			);
			return correctOption ? correctOption[0] : "";
		}
	};

	// Helper function to check if user answer is correct
	const isAnswerCorrect = (
		question: Question,
		userAnswer: UserAnswer
	): boolean => {
		if (userAnswer.isSkipped) return false;

		const correctAnswer = getCorrectAnswer(question);

		if (question.question_type === "SBA") {
			return userAnswer.selectedOptions.includes(correctAnswer);
		} else {
			return userAnswer.selectedOptions[0] === correctAnswer;
		}
	};

	// Analyze question types performance
	const analyzePerformance = () => {
		const sbaQuestions = exam.questions.filter(
			(q) => q.question_type === "SBA"
		);
		const truefalseQuestions = exam.questions.filter(
			(q) => q.question_type === "TRUE_FALSE"
		);

		let sbaCorrect = 0;
		let truefalseCorrect = 0;

		exam.questions.forEach((question) => {
			const answer = userAnswers[question.id];
			if (answer && !answer.isSkipped) {
				const isCorrect = isAnswerCorrect(question, answer);

				if (isCorrect) {
					if (question.question_type === "SBA") {
						sbaCorrect++;
					} else {
						truefalseCorrect++;
					}
				}
			}
		});

		return {
			sba: {
				correct: sbaCorrect,
				total: sbaQuestions.length,
				percentage:
					sbaQuestions.length > 0
						? Math.round((sbaCorrect / sbaQuestions.length) * 100)
						: 0,
			},
			truefalse: {
				correct: truefalseCorrect,
				total: truefalseQuestions.length,
				percentage:
					truefalseQuestions.length > 0
						? Math.round(
								(truefalseCorrect / truefalseQuestions.length) *
									100
						  )
						: 0,
			},
		};
	};

	const performance = analyzePerformance();

	return (
		<div className="space-y-6">
			{/* Results Header */}
			<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-8 text-center">
				<div
					className={`inline-flex items-center justify-center w-20 h-20 rounded-full mb-6 ${gradeInfo.bgColor}`}
				>
					<Award className={`w-10 h-10 ${gradeInfo.color}`} />
				</div>

				<h2 className="text-3xl font-bold text-gray-900 mb-2">
					Exam Complete!
				</h2>

				<div className="mb-6">
					<div
						className={`inline-block text-5xl font-bold mb-2 ${gradeInfo.color}`}
					>
						{percentage}%
					</div>
					<div
						className={`inline-block ml-4 text-2xl font-semibold px-4 py-2 rounded-lg ${gradeInfo.bgColor} ${gradeInfo.color}`}
					>
						Grade: {gradeInfo.grade}
					</div>
				</div>

				<p className="text-lg text-gray-600 max-w-2xl mx-auto">
					{getPerformanceMessage(percentage)}
				</p>
			</div>

			{/* Score Breakdown */}
			<div className="grid grid-cols-1 md:grid-cols-4 gap-4">
				<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 text-center">
					<CheckCircle className="w-8 h-8 text-green-500 mx-auto mb-3" />
					<p className="text-2xl font-bold text-green-600">
						{correctAnswers}
					</p>
					<p className="text-sm text-gray-600">Correct</p>
				</div>

				<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 text-center">
					<XCircle className="w-8 h-8 text-red-500 mx-auto mb-3" />
					<p className="text-2xl font-bold text-red-600">
						{incorrectAnswers}
					</p>
					<p className="text-sm text-gray-600">Incorrect</p>
				</div>

				<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 text-center">
					<Clock className="w-8 h-8 text-blue-500 mx-auto mb-3" />
					<p className="text-2xl font-bold text-blue-600">
						{formatTime(timeSpent)}
					</p>
					<p className="text-sm text-gray-600">Time Spent</p>
				</div>

				<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 text-center">
					<BookOpen className="w-8 h-8 text-purple-500 mx-auto mb-3" />
					<p className="text-2xl font-bold text-purple-600">
						{totalQuestions}
					</p>
					<p className="text-sm text-gray-600">Total Questions</p>
				</div>
			</div>

			{/* Question Type Performance */}
			<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
				<h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
					<TrendingUp className="w-5 h-5" />
					Performance by Question Type
				</h3>

				<div className="grid grid-cols-1 md:grid-cols-2 gap-6">
					{performance.sba.total > 0 && (
						<div className="space-y-3">
							<div className="flex justify-between items-center">
								<span className="font-medium text-gray-900">
									Single Best Answer (SBA)
								</span>
								<span className="text-sm text-gray-600">
									{performance.sba.correct}/
									{performance.sba.total}
								</span>
							</div>
							<div className="w-full bg-gray-200 rounded-full h-3">
								<div
									className="bg-green-500 h-3 rounded-full transition-all duration-500"
									style={{
										width: `${performance.sba.percentage}%`,
									}}
								></div>
							</div>
							<p className="text-sm text-gray-600 text-right">
								{performance.sba.percentage}%
							</p>
						</div>
					)}

					{performance.truefalse.total > 0 && (
						<div className="space-y-3">
							<div className="flex justify-between items-center">
								<span className="font-medium text-gray-900">
									True/False
								</span>
								<span className="text-sm text-gray-600">
									{performance.truefalse.correct}/
									{performance.truefalse.total}
								</span>
							</div>
							<div className="w-full bg-gray-200 rounded-full h-3">
								<div
									className="bg-purple-500 h-3 rounded-full transition-all duration-500"
									style={{
										width: `${performance.truefalse.percentage}%`,
									}}
								></div>
							</div>
							<p className="text-sm text-gray-600 text-right">
								{performance.truefalse.percentage}%
							</p>
						</div>
					)}
				</div>
			</div>

			{/* Detailed Results */}
			<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
				<h3 className="text-lg font-semibold text-gray-900 mb-4">
					Question by Question Results
				</h3>

				<div className="space-y-3">
					{exam.questions.map((question, index) => {
						const answer = userAnswers[question.id];
						let isCorrect = false;
						let userAnswerText = "Not answered";

						if (answer && !answer.isSkipped) {
							isCorrect = isAnswerCorrect(question, answer);
							userAnswerText = answer.selectedOptions
								.join(", ")
								.toUpperCase();
						}

						const correctAnswerText = getCorrectAnswer(question);

						return (
							<div
								key={question.id}
								className={`p-4 rounded-lg border-l-4 ${
									answer && !answer.isSkipped
										? isCorrect
											? "border-green-500 bg-green-50"
											: "border-red-500 bg-red-50"
										: "border-gray-300 bg-gray-50"
								}`}
							>
								<div className="flex items-start justify-between">
									<div className="flex-1">
										<div className="flex items-center gap-3 mb-2">
											<span className="font-semibold text-gray-900">
												Question {index + 1}
											</span>
											<span
												className={`text-xs px-2 py-1 rounded-full font-medium ${
													question.question_type ===
													"SBA"
														? "bg-green-100 text-green-700"
														: "bg-purple-100 text-purple-700"
												}`}
											>
												{question.question_type ===
												"SBA"
													? "SBA"
													: "T/F"}
											</span>
											{answer && !answer.isSkipped ? (
												isCorrect ? (
													<CheckCircle className="w-4 h-4 text-green-500" />
												) : (
													<XCircle className="w-4 h-4 text-red-500" />
												)
											) : (
												<XCircle className="w-4 h-4 text-gray-400" />
											)}
										</div>

										<p className="text-sm text-gray-700 mb-3 line-clamp-2">
											{question.question_text}
										</p>

										<div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
											<div>
												<span className="font-medium text-gray-900">
													Your answer:{" "}
												</span>
												<span
													className={
														answer &&
														!answer.isSkipped
															? isCorrect
																? "text-green-600"
																: "text-red-600"
															: "text-gray-500"
													}
												>
													{userAnswerText}
												</span>
											</div>
											<div>
												<span className="font-medium text-gray-900">
													Correct answer:{" "}
												</span>
												<span className="text-green-600">
													{correctAnswerText.toUpperCase()}
												</span>
											</div>
										</div>
									</div>
								</div>
							</div>
						);
					})}
				</div>
			</div>

			{/* Action Buttons */}
			<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
				<div className="flex flex-col sm:flex-row gap-4 justify-center">
					<button
						onClick={onRetakeExam}
						className="px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-colors"
					>
						Retake Exam
					</button>

					{onViewSolutions && (
						<button
							onClick={onViewSolutions}
							className="px-6 py-3 bg-gray-600 text-white rounded-lg font-semibold hover:bg-gray-700 transition-colors"
						>
							View Solutions
						</button>
					)}
				</div>
			</div>
		</div>
	);
}
