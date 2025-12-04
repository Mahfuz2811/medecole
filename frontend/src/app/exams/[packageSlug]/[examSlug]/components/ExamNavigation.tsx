import { ChevronLeft, ChevronRight, Eye, Send } from "lucide-react";

interface ExamNavigationProps {
	canGoNext: boolean;
	canGoPrevious: boolean;
	canSubmit: boolean;
	isReviewMode: boolean;
	onPrevious: () => void;
	onNext: () => void;
	onReview: () => void;
	onSubmit: () => void;
	currentQuestionIndex: number;
	totalQuestions: number;
	answeredCount: number;
}

export function ExamNavigation({
	canGoNext,
	canGoPrevious,
	canSubmit,
	isReviewMode,
	onPrevious,
	onNext,
	onReview,
	onSubmit,
	currentQuestionIndex,
	totalQuestions,
	answeredCount,
}: ExamNavigationProps) {
	const isLastQuestion = currentQuestionIndex === totalQuestions - 1;

	if (isReviewMode) {
		return (
			<div className="p-4">
				<div className="flex items-center justify-center">
					<button
						onClick={onSubmit}
						disabled={!canSubmit}
						className={`flex items-center gap-2 px-6 py-3 rounded-lg font-semibold transition-colors ${
							canSubmit
								? "bg-green-600 text-white hover:bg-green-700"
								: "bg-gray-100 text-gray-400 cursor-not-allowed"
						}`}
					>
						<Send className="w-4 h-4" />
						Submit Exam
					</button>
				</div>
				{!canSubmit && (
					<p className="text-center text-sm text-gray-500 mt-2">
						Answer at least one question to submit
					</p>
				)}
			</div>
		);
	}

	return (
		<div className="p-4">
			<div className="flex items-center justify-between">
				{/* Previous Button */}
				<button
					onClick={onPrevious}
					disabled={!canGoPrevious}
					className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${
						canGoPrevious
							? "bg-gray-100 text-gray-700 hover:bg-gray-200"
							: "bg-gray-50 text-gray-400 cursor-not-allowed"
					}`}
				>
					<ChevronLeft className="w-4 h-4" />
					Previous
				</button>

				{/* Center Actions */}
				<div className="flex items-center gap-3">
					<span className="text-sm text-gray-600">
						{currentQuestionIndex + 1} of {totalQuestions}
					</span>

					{/* Review Button (show when at least some questions answered) */}
					{answeredCount > 0 && (
						<button
							onClick={onReview}
							className="flex items-center gap-2 px-4 py-2 bg-orange-100 text-orange-700 rounded-lg font-medium hover:bg-orange-200 transition-colors"
						>
							<Eye className="w-4 h-4" />
							Review ({answeredCount})
						</button>
					)}
				</div>

				{/* Next/Submit Button */}
				{isLastQuestion ? (
					<button
						onClick={onReview}
						className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
					>
						<Eye className="w-4 h-4" />
						Review & Submit
					</button>
				) : (
					<button
						onClick={onNext}
						disabled={!canGoNext}
						className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${
							canGoNext
								? "bg-blue-600 text-white hover:bg-blue-700"
								: "bg-gray-100 text-gray-400 cursor-not-allowed"
						}`}
					>
						Next
						<ChevronRight className="w-4 h-4" />
					</button>
				)}
			</div>
		</div>
	);
}
