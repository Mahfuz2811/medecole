import { AlertTriangle, Award, Clock, FileText, Target } from "lucide-react";
import { ExamData } from "../types";

interface ExamInstructionsProps {
	exam: ExamData;
	onStartExam: () => void;
}

export function ExamInstructions({ exam, onStartExam }: ExamInstructionsProps) {
	// Calculate pass percentage and format passing score
	const passPercentage = exam.total_marks
		? Math.round((exam.passing_score / exam.total_marks) * 100)
		: Math.round(exam.passing_score);

	return (
		<div className="min-h-screen bg-blue-50">
			{/* Header */}
			<div className="bg-white shadow-sm border-b">
				<div className="max-w-4xl mx-auto px-4 py-6">
					<h1 className="text-2xl font-bold text-gray-900 mb-2">
						{exam.title}
					</h1>
					<p className="text-gray-600">
						{exam.description ||
							"Test your knowledge with this comprehensive exam."}
					</p>
				</div>
			</div>

			{/* Main Content */}
			<div className="max-w-4xl mx-auto px-4 py-8">
				<div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
					{/* Exam Overview */}
					<div className="p-6 border-b border-gray-200">
						<h2 className="text-lg font-semibold text-gray-900 mb-4">
							Exam Overview
						</h2>
						<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
							<div className="flex items-center gap-3 p-4 bg-blue-50 rounded-lg">
								<FileText className="w-5 h-5 text-blue-600" />
								<div>
									<p className="text-sm font-medium text-gray-900">
										Questions
									</p>
									<p className="text-lg font-bold text-blue-600">
										{exam.total_questions}
									</p>
								</div>
							</div>
							<div className="flex items-center gap-3 p-4 bg-orange-50 rounded-lg">
								<Clock className="w-5 h-5 text-orange-600" />
								<div>
									<p className="text-sm font-medium text-gray-900">
										Duration
									</p>
									<p className="text-lg font-bold text-orange-600">
										{exam.duration_minutes} min
									</p>
								</div>
							</div>
							<div className="flex items-center gap-3 p-4 bg-purple-50 rounded-lg">
								<Award className="w-5 h-5 text-purple-600" />
								<div>
									<p className="text-sm font-medium text-gray-900">
										Total Marks
									</p>
									<p className="text-lg font-bold text-purple-600">
										{exam.total_marks || "N/A"}
									</p>
								</div>
							</div>
							<div className="flex items-center gap-3 p-4 bg-green-50 rounded-lg">
								<Target className="w-5 h-5 text-green-600" />
								<div>
									<p className="text-sm font-medium text-gray-900">
										Passing Score
									</p>
									<p className="text-lg font-bold text-green-600">
										{exam.total_marks
											? `${passPercentage}% to pass`
											: `${passPercentage} to pass`}
									</p>
								</div>
							</div>
						</div>
					</div>

					{/* Instructions */}
					<div className="p-6">
						<h3 className="text-lg font-semibold text-gray-900 mb-4">
							Instructions
						</h3>
						<div className="prose prose-sm max-w-none text-gray-700">
							<div className="space-y-4">
								<p>
									<strong>General Instructions:</strong>
								</p>
								<ul className="list-disc pl-6 space-y-2">
									<li>
										This exam contains{" "}
										{exam.total_questions} questions
									</li>
									<li>
										You have {exam.duration_minutes} minutes
										to complete the exam
									</li>
									<li>
										{exam.total_marks ? (
											<span>
												You need to score{" "}
												{passPercentage}% or higher to
												pass this exam (
												{Math.round(exam.passing_score)}{" "}
												out of {exam.total_marks} marks)
											</span>
										) : (
											<span>
												You need {passPercentage} or
												higher to pass
											</span>
										)}
									</li>
									<li>
										You can navigate between questions and
										change your answers
									</li>
									<li>
										Click &quot;Submit Exam&quot; when
										you&apos;re ready to finish
									</li>
									<li>
										Make sure to review your answers before
										submitting
									</li>
								</ul>
								<p>
									<strong>Question Types:</strong>
								</p>
								<ul className="list-disc pl-6 space-y-2">
									<li>
										<strong>TRUE/FALSE:</strong> Select T
										for True and F for False
									</li>
									<li>
										<strong>
											SBA (Single Best Answer):
										</strong>{" "}
										Select the ONE best answer
									</li>
								</ul>

								{/* Custom Instructions Section */}
								{exam.instructions &&
									exam.instructions.trim() !== "" && (
										<div className="mt-6">
											<p>
												<strong>
													Additional Instructions:
												</strong>
											</p>
											<div className="whitespace-pre-line mt-2 p-4 bg-blue-50 rounded-lg border-l-4 border-blue-400">
												{exam.instructions}
											</div>
										</div>
									)}
							</div>
						</div>
					</div>

					{/* Important Notice */}
					<div className="p-6 bg-yellow-50 border-t border-yellow-200">
						<div className="flex items-start gap-3">
							<AlertTriangle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
							<div>
								<h4 className="font-medium text-yellow-800 mb-1">
									Important Notice
								</h4>
								<p className="text-sm text-yellow-700">
									Once you start the exam, the timer will
									begin automatically. Make sure you have a
									stable internet connection and won&apos;t be
									interrupted. The exam will auto-submit when
									time runs out.
								</p>
							</div>
						</div>
					</div>

					{/* Start Button */}
					<div className="p-6 bg-gray-50 border-t border-gray-200">
						<div className="flex justify-center">
							<button
								onClick={onStartExam}
								className="bg-blue-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors shadow-sm"
							>
								Start Exam
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
