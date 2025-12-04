import { Eye, Send } from "lucide-react";

interface ExamNavigationListProps {
	canSubmit: boolean;
	isReviewMode: boolean;
	onReview: () => void;
	onSubmit: () => void;
	answeredCount: number;
}

export function ExamNavigationList({
	canSubmit,
	isReviewMode,
	onReview,
	onSubmit,
	answeredCount,
}: ExamNavigationListProps) {
	if (isReviewMode) {
		return (
			<div className="bg-white border-t border-gray-200 p-4">
				<div className="max-w-4xl mx-auto">
					<div className="flex items-center justify-center">
						<button
							onClick={onSubmit}
							disabled={!canSubmit}
							className={`flex items-center gap-2 px-8 py-3 rounded-lg font-semibold transition-colors ${
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
			</div>
		);
	}

	return (
		<div className="bg-white border-t border-gray-200 p-4">
			<div className="max-w-4xl mx-auto">
				<div className="flex items-center justify-center">
					{/* Action Buttons */}
					<div className="flex items-center gap-3">
						{/* Review Button (show when at least some questions answered) */}
						{answeredCount > 0 && (
							<button
								onClick={onReview}
								className="flex items-center gap-2 px-4 py-2 bg-orange-100 text-orange-700 rounded-lg font-medium hover:bg-orange-200 transition-colors"
							>
								<Eye className="w-4 h-4" />
								Review & Submit ({answeredCount})
							</button>
						)}
					</div>
				</div>
			</div>
		</div>
	);
}
