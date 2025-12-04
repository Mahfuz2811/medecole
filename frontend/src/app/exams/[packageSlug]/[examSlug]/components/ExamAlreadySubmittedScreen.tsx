import { ArrowLeft, CheckCircle, Clock, TrendingUp } from "lucide-react";
import Link from "next/link";

interface PreviousAttempt {
	attempt_id: number;
	submitted_at: string;
	score: number;
	status: string;
}

interface ExamAlreadySubmittedData {
	previous_attempt: PreviousAttempt;
	can_retry: boolean;
}

interface ExamConflictError {
	error: string;
	message: string;
	data: ExamAlreadySubmittedData;
}

interface ExamAlreadySubmittedProps {
	error: ExamConflictError;
	examTitle?: string;
	packageSlug: string;
}

export function ExamAlreadySubmittedScreen({
	error,
	examTitle,
	packageSlug,
}: ExamAlreadySubmittedProps) {
	// Handle the nested data structure safely
	const previousAttempt = error?.data?.previous_attempt;
	
	if (!previousAttempt) {
		return (
			<div className="min-h-screen bg-gray-50 flex items-center justify-center">
				<div className="text-center p-6">
					<h1 className="text-xl font-semibold text-gray-900 mb-2">
						Exam Already Completed
					</h1>
					<p className="text-gray-600 mb-4">
						You have already submitted this exam.
					</p>
					<Link
						href={`/exams/${packageSlug}`}
						className="bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors"
					>
						Back to Exams
					</Link>
				</div>
			</div>
		);
	}
	
	const { previous_attempt } = error.data;

	// Format the submission date and time
	const submissionDate = new Date(previous_attempt.submitted_at);
	const formattedDate = submissionDate.toLocaleDateString("en-US", {
		year: "numeric",
		month: "long",
		day: "numeric",
	});
	const formattedTime = submissionDate.toLocaleTimeString("en-US", {
		hour: "2-digit",
		minute: "2-digit",
	});

	// Calculate time ago
	const timeAgo = (() => {
		const now = new Date();
		const diffInHours = Math.floor(
			(now.getTime() - submissionDate.getTime()) / (1000 * 60 * 60)
		);

		if (diffInHours < 1) {
			const diffInMinutes = Math.floor(
				(now.getTime() - submissionDate.getTime()) / (1000 * 60)
			);
			return `${diffInMinutes} minutes ago`;
		} else if (diffInHours < 24) {
			return `${diffInHours} hours ago`;
		} else {
			const diffInDays = Math.floor(diffInHours / 24);
			return `${diffInDays} days ago`;
		}
	})();

	const isPassed = previous_attempt.score >= 60; // Assuming 60% is passing

	return (
		<div className="min-h-screen bg-gray-50">
			{/* Header */}
			<div className="bg-white border-b border-gray-200">
				<div className="max-w-4xl mx-auto px-4 py-4">
					<Link
						href={`/exams/${packageSlug}`}
						className="inline-flex items-center text-gray-600 hover:text-gray-900 transition-colors"
					>
						<ArrowLeft className="w-4 h-4 mr-2" />
						Back to Exams
					</Link>
				</div>
			</div>

			{/* Main Content */}
			<div className="max-w-2xl mx-auto px-4 py-8">
				<div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
					{/* Status Header */}
					<div className="bg-green-50 border-b border-green-100 px-6 py-4">
						<div className="flex items-center">
							<div className="flex-shrink-0">
								<CheckCircle className="w-8 h-8 text-green-600" />
							</div>
							<div className="ml-3">
								<h1 className="text-lg font-semibold text-green-900">
									Exam Already Completed
								</h1>
								<p className="text-sm text-green-700 mt-1">
									You have successfully submitted this exam
								</p>
							</div>
						</div>
					</div>

					{/* Exam Details */}
					<div className="px-6 py-6">
						<h2 className="text-xl font-semibold text-gray-900 mb-4">
							{examTitle}
						</h2>

						{/* Score Card */}
						<div className="bg-gray-50 rounded-lg p-4 mb-6">
							<div className="flex items-center justify-between">
								<div className="flex items-center">
									<TrendingUp
										className={`w-5 h-5 mr-2 ${
											isPassed
												? "text-green-600"
												: "text-red-600"
										}`}
									/>
									<span className="text-sm font-medium text-gray-700">
										Your Score
									</span>
								</div>
								<div className="text-right">
									<div
										className={`text-2xl font-bold ${
											isPassed
												? "text-green-600"
												: "text-red-600"
										}`}
									>
										{previous_attempt.score.toFixed(1)}%
									</div>
									<div
										className={`text-xs font-medium ${
											isPassed
												? "text-green-600"
												: "text-red-600"
										}`}
									>
										{isPassed ? "Passed" : "Failed"}
									</div>
								</div>
							</div>
						</div>

						{/* Submission Details */}
						<div className="space-y-4">
							<div className="flex items-center text-sm text-gray-600">
								<Clock className="w-4 h-4 mr-2" />
								<span>
									Submitted {timeAgo} on {formattedDate} at{" "}
									{formattedTime}
								</span>
							</div>

							<div className="text-sm text-gray-600">
								<span className="font-medium">Status:</span>{" "}
								{previous_attempt.status}
							</div>

							<div className="text-sm text-gray-600">
								<span className="font-medium">Attempt ID:</span>{" "}
								#{previous_attempt.attempt_id}
							</div>
						</div>

						{/* Info Message */}
						<div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
							<p className="text-sm text-blue-800">
								<strong>Note:</strong> {error.message}
							</p>
						</div>
					</div>

					{/* Actions */}
					<div className="px-6 py-4 bg-gray-50 border-t border-gray-200">
						<div className="flex flex-col sm:flex-row gap-3">
							<Link
								href={`/exams/${packageSlug}`}
								className="flex-1 bg-blue-600 text-white text-center py-3 px-4 rounded-lg font-medium hover:bg-blue-700 transition-colors"
							>
								View Other Exams
							</Link>
							<Link
								href="/dashboard"
								className="flex-1 bg-white border border-gray-300 text-gray-700 text-center py-3 px-4 rounded-lg font-medium hover:bg-gray-50 transition-colors"
							>
								Go to Dashboard
							</Link>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
