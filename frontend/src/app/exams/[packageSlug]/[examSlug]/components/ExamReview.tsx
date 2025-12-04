import { AlertCircle, ArrowLeft, CheckCircle, XCircle } from "lucide-react";
import { ExamData, UserAnswer } from "../types";

interface ExamReviewProps {
	exam: ExamData;
	userAnswers: Record<number, UserAnswer>;
	onBackToExam: () => void;
}

export function ExamReview({
	exam,
	userAnswers,
	onBackToExam,
}: ExamReviewProps) {
	const answeredQuestions = Object.keys(userAnswers).length;
	const unansweredQuestions = exam.total_questions - answeredQuestions;

	const getQuestionStatus = (questionId: number) => {
		const answer = userAnswers[questionId];
		if (!answer || answer.isSkipped) return "unanswered";
		return "answered";
	};

	const getStatusIcon = (status: string) => {
		switch (status) {
			case "answered":
				return <CheckCircle className="w-5 h-5 text-green-500" />;
			case "unanswered":
				return <XCircle className="w-5 h-5 text-red-500" />;
			default:
				return <AlertCircle className="w-5 h-5 text-yellow-500" />;
		}
	};

	const getStatusColor = (status: string) => {
		switch (status) {
			case "answered":
				return "border-green-200 bg-green-50";
			case "unanswered":
				return "border-red-200 bg-red-50";
			default:
				return "border-yellow-200 bg-yellow-50";
		}
	};

	return (
		<div className="space-y-4">
			{/* Review Header */}
			<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
				<div className="flex items-center justify-between mb-4">
					<h2 className="text-lg font-semibold text-gray-900">
						Review Answers
					</h2>
					<button
						onClick={onBackToExam}
						className="flex items-center gap-1 px-2 py-1.5 sm:px-4 sm:py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors text-xs sm:text-sm"
					>
						<ArrowLeft className="w-3 h-3 sm:w-4 sm:h-4" />
						<span className="whitespace-nowrap">Back to Exam</span>
					</button>
				</div>

				{/* Summary Stats */}
				<div className="grid grid-cols-3 gap-3 mb-4">
					<div className="text-center p-3 bg-green-50 rounded-lg">
						<p className="text-xl font-bold text-green-600">
							{answeredQuestions}
						</p>
						<p className="text-xs text-green-700">Answered</p>
					</div>
					<div className="text-center p-3 bg-red-50 rounded-lg">
						<p className="text-xl font-bold text-red-600">
							{unansweredQuestions}
						</p>
						<p className="text-xs text-red-700">Unanswered</p>
					</div>
					<div className="text-center p-3 bg-blue-50 rounded-lg">
						<p className="text-xl font-bold text-blue-600">
							{Math.round(
								(answeredQuestions / exam.total_questions) * 100
							)}
							%
						</p>
						<p className="text-xs text-blue-700">Complete</p>
					</div>
				</div>
			</div>

			{/* Question Details */}
			<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
				<h3 className="text-base font-semibold text-gray-900 mb-3">
					Answer Details
				</h3>
				<div className="space-y-3">
					{exam.questions.map((question, index) => {
						const status = getQuestionStatus(question.id);
						const answer = userAnswers[question.id];

						return (
							<div
								key={question.id}
								className={`p-3 rounded-lg border ${getStatusColor(
									status
								)}`}
							>
								<div className="flex items-start">
									<div className="flex-1">
										<div className="flex items-center gap-2 mb-2">
											<span className="font-semibold text-gray-900 text-sm">
												Q{index + 1}
											</span>
											<span
												className={`text-xs px-2 py-0.5 rounded-full font-medium ${
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
											{getStatusIcon(status)}
										</div>
										<p className="text-xs text-gray-700 mb-2 line-clamp-2">
											{question.question_text}
										</p>
										{answer && !answer.isSkipped ? (
											<div className="text-xs">
												<span className="font-medium text-gray-900">
													Answer:
												</span>
												<span className="ml-2 px-2 py-1 bg-blue-100 text-blue-800 rounded font-medium">
													{answer.selectedOptions
														.join(", ")
														.toUpperCase()}
												</span>
											</div>
										) : (
											<div className="text-xs">
												<span className="px-2 py-1 bg-red-100 text-red-800 rounded font-medium">
													Not answered
												</span>
											</div>
										)}
									</div>
								</div>
							</div>
						);
					})}
				</div>
			</div>
		</div>
	);
}
