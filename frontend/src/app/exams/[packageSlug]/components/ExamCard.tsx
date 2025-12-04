import { ExamResponse } from "@/lib/api/types";
import Link from "next/link";
import { useParams } from "next/navigation";

// Helper function to generate the correct results URL
function getResultsUrl(
	packageSlug: string,
	examSlug: string,
	sessionId?: string
): string {
	if (sessionId) {
		// Use session-based results URL if session_id is available
		return `/exams/${packageSlug}/${examSlug}/results?session=${sessionId}`;
	}
	// Fallback to the old URL format
	return `/exams/${packageSlug}/${examSlug}/results`;
}

// Using backend API data directly - no need for extended interface
interface ExamCardProps {
	exam: ExamResponse;
}

export function ExamCard({ exam }: ExamCardProps) {
	const params = useParams();
	const packageSlug = params.packageSlug as string;

	const getStatusIcon = (exam: ExamResponse) => {
		// Priority 1: Check computed status from backend
		switch (exam.computed_status) {
			case "UPCOMING":
				return (
					<div className="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-yellow-100 flex items-center justify-center">
						<svg
							className="w-4 h-4 sm:w-5 sm:h-5 text-yellow-600"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
					</div>
				);
			case "LIVE":
				return (
					<div className="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-red-100 flex items-center justify-center animate-pulse">
						<svg
							className="w-4 h-4 sm:w-5 sm:h-5 text-red-600"
							fill="currentColor"
							viewBox="0 0 24 24"
						>
							<circle cx="12" cy="12" r="3" />
						</svg>
					</div>
				);
			case "COMPLETED":
				return (
					<div className="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-gray-100 flex items-center justify-center">
						<svg
							className="w-4 h-4 sm:w-5 sm:h-5 text-gray-600"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
					</div>
				);
			case "AVAILABLE":
			default:
				// Priority 2: Check user attempt status for available exams
				if (exam.user_attempt?.status === "COMPLETED") {
					const isPassed = exam.user_attempt.is_passed;
					return (
						<div
							className={`w-8 h-8 sm:w-10 sm:h-10 rounded-full flex items-center justify-center ${
								isPassed ? "bg-green-100" : "bg-red-100"
							}`}
						>
							<svg
								className={`w-4 h-4 sm:w-5 sm:h-5 ${
									isPassed ? "text-green-600" : "text-red-600"
								}`}
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								{isPassed ? (
									<path
										strokeLinecap="round"
										strokeLinejoin="round"
										strokeWidth={2}
										d="M5 13l4 4L19 7"
									/>
								) : (
									<path
										strokeLinecap="round"
										strokeLinejoin="round"
										strokeWidth={2}
										d="M6 18L18 6M6 6l12 12"
									/>
								)}
							</svg>
						</div>
					);
				}

				if (exam.user_attempt?.status === "STARTED") {
					return (
						<div className="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-orange-100 flex items-center justify-center">
							<svg
								className="w-4 h-4 sm:w-5 sm:h-5 text-orange-600"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
						</div>
					);
				}

				// Default: Available to start
				return (
					<div className="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-blue-100 flex items-center justify-center group-hover:bg-blue-200 transition-colors">
						<svg
							className="w-4 h-4 sm:w-5 sm:h-5 text-blue-600"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M9 5l7 7-7 7"
							/>
						</svg>
					</div>
				);
		}
	};

	const formatDuration = (minutes: number) => {
		if (minutes < 60) {
			return `${minutes} min`;
		}
		const hours = Math.floor(minutes / 60);
		const remainingMinutes = minutes % 60;
		return remainingMinutes > 0
			? `${hours}h ${remainingMinutes}m`
			: `${hours}h`;
	};

	const formatScheduledTime = (dateString?: string) => {
		if (!dateString) return "";
		const date = new Date(dateString);
		return date.toLocaleDateString("en-US", {
			month: "short",
			day: "numeric",
			hour: "numeric",
			minute: "2-digit",
			hour12: true,
			timeZone: "UTC",
		});
	};

	const getStatusMessage = (exam: ExamResponse) => {
		switch (exam.computed_status) {
			case "UPCOMING":
				return exam.scheduled_start_date
					? `Starts ${formatScheduledTime(exam.scheduled_start_date)}`
					: "Scheduled for later";
			case "LIVE":
				return exam.scheduled_end_date
					? `Ends ${formatScheduledTime(exam.scheduled_end_date)}`
					: "Live now";
			case "COMPLETED":
				return exam.scheduled_end_date
					? `Ended ${formatScheduledTime(exam.scheduled_end_date)}`
					: "Exam completed";
			default:
				return "";
		}
	};

	return (
		<div className="bg-white shadow-sm border border-gray-200 hover:shadow-lg hover:border-blue-300 transition-all duration-200 group">
			{/* Header with Status Indicator */}
			<div className="relative">
				{exam.computed_status === "LIVE" && (
					<div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-red-500 to-red-600 rounded-t-xl animate-pulse"></div>
				)}
				{exam.computed_status === "UPCOMING" && (
					<div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-yellow-400 to-yellow-500 rounded-t-xl"></div>
				)}
			</div>

			<div className="p-5">
				{/* Top Section: Title + Action */}
				<div className="flex items-start justify-between mb-4">
					<div className="flex-1 pr-4">
						<div className="flex items-center gap-3 mb-2">
							{getStatusIcon(exam)}
							<div>
								<h3 className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors leading-tight">
									{exam.title}
								</h3>
								{getStatusMessage(exam) && (
									<p className="text-sm text-blue-600 font-medium mt-1">
										{getStatusMessage(exam)}
									</p>
								)}
							</div>
						</div>
					</div>

					{/* Action Button - Desktop */}
					<div className="hidden sm:block flex-shrink-0">
						{(() => {
							// Priority 1: Check user attempt status first
							if (exam.user_attempt) {
								switch (exam.user_attempt.status) {
									case "STARTED":
										return (
											<button
												onClick={() =>
													(window.location.href = `/exams/${packageSlug}/${exam.slug}`)
												}
												className="px-4 py-2 bg-orange-500 text-white rounded-lg text-sm font-semibold hover:bg-orange-600 transition-colors"
											>
												Resume Exam
											</button>
										);
									case "COMPLETED":
									case "AUTO_SUBMITTED":
										return (
											<Link
												href={getResultsUrl(
													packageSlug,
													exam.slug,
													exam.user_attempt
														?.session_id
												)}
											>
												<button className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-200 transition-colors">
													View Results
												</button>
											</Link>
										);
									case "ABANDONED":
										return (
											<Link
												href={getResultsUrl(
													packageSlug,
													exam.slug,
													exam.user_attempt
														?.session_id
												)}
											>
												<button className="px-4 py-2 bg-blue-500 text-white rounded-lg text-sm font-semibold hover:bg-blue-600 transition-colors">
													Review Answers
												</button>
											</Link>
										);
								}
							}

							// Priority 2: Check computed status for non-attempted exams
							switch (exam.computed_status) {
								case "UPCOMING":
									return (
										<div className="px-4 py-2 bg-yellow-50 text-yellow-700 rounded-lg text-sm font-medium border border-yellow-200">
											Upcoming
										</div>
									);
								case "LIVE":
									return (
										<button
											onClick={() =>
												(window.location.href = `/exams/${packageSlug}/${exam.slug}`)
											}
											className="px-6 py-2 bg-red-500 text-white rounded-lg text-sm font-semibold hover:bg-red-600 transition-colors shadow-sm"
										>
											Join Live
										</button>
									);
								case "COMPLETED":
									return (
										<button
											onClick={() =>
												(window.location.href = `/exams/${packageSlug}/${exam.slug}`)
											}
											className="px-4 py-2 bg-blue-500 text-white rounded-lg text-sm font-semibold hover:bg-blue-600 transition-colors"
										>
											Start Exam
										</button>
									);
								case "AVAILABLE":
								default:
									return (
										<Link
											href={`/exams/${packageSlug}/${exam.slug}`}
											className="px-6 py-2 bg-blue-500 text-white rounded-lg text-sm font-semibold hover:bg-blue-600 transition-colors shadow-sm"
										>
											Start
										</Link>
									);
							}
						})()}
					</div>
				</div>

				{/* Middle Section: Details */}
				<div className="mb-4">
					<div className="flex items-center gap-6 text-sm text-gray-600">
						<div className="flex items-center gap-1">
							<svg
								className="w-4 h-4"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
							<span>{exam.total_questions} questions</span>
						</div>
						<div className="flex items-center gap-1">
							<svg
								className="w-4 h-4"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
							<span>{formatDuration(exam.duration_minutes)}</span>
						</div>
						<div className="flex items-center gap-1">
							<svg
								className="w-4 h-4"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
							<span>{exam.passing_score} to pass</span>
						</div>
					</div>
				</div>

				{/* Description */}
				{exam.description && (
					<div className="mb-4">
						<p className="text-sm text-gray-600 leading-relaxed line-clamp-2">
							{exam.description}
						</p>
					</div>
				)}

				{/* User Progress - Only show if user has completed the exam */}
				{exam.user_attempt &&
					(exam.user_attempt.status === "COMPLETED" ||
						exam.user_attempt.status === "AUTO_SUBMITTED" ||
						exam.user_attempt.status === "ABANDONED") && (
						<div className="bg-gradient-to-r from-gray-50 to-gray-100 rounded-lg p-4 mb-4 border border-gray-200">
							<div className="flex items-center justify-between">
								<div className="flex items-center gap-8">
									<div className="text-center">
										<p className="text-xs font-medium text-gray-500 mb-1">
											Your Score
										</p>
										<p
											className={`text-xl font-bold ${
												exam.user_attempt.is_passed
													? "text-green-600"
													: "text-red-600"
											}`}
										>
											{exam.user_attempt.score?.toFixed(
												1
											)}
										</p>
									</div>
									<div className="text-center">
										<p className="text-xs font-medium text-gray-500 mb-1">
											Accuracy
										</p>
										<p className="text-xl font-bold text-gray-900">
											{exam.user_attempt.correct_answers}/
											{exam.total_questions}
										</p>
									</div>
								</div>
								<div className="flex items-center gap-3">
									<span
										className={`px-4 py-2 rounded-full text-sm font-semibold ${
											exam.user_attempt.is_passed
												? "bg-green-100 text-green-700 border border-green-200"
												: "bg-red-100 text-red-700 border border-red-200"
										}`}
									>
										{exam.user_attempt.is_passed
											? "PASSED"
											: "FAILED"}
									</span>
								</div>
							</div>
						</div>
					)}

				{/* Mobile Action Button */}
				<div className="sm:hidden mb-4">
					{(() => {
						// Priority 1: Check user attempt status first
						if (exam.user_attempt) {
							switch (exam.user_attempt.status) {
								case "STARTED":
									return (
										<button
											onClick={() =>
												(window.location.href = `/exams/${packageSlug}/${exam.slug}`)
											}
											className="w-full px-4 py-3 bg-orange-500 text-white rounded-lg font-semibold hover:bg-orange-600 transition-colors"
										>
											Resume Exam
										</button>
									);
								case "COMPLETED":
								case "AUTO_SUBMITTED":
									return (
										<Link
											href={getResultsUrl(
												packageSlug,
												exam.slug,
												exam.user_attempt?.session_id
											)}
											className="block w-full"
											onClick={(e) => e.stopPropagation()}
										>
											<button className="w-full px-4 py-3 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 transition-colors">
												View Exam Results
											</button>
										</Link>
									);
								case "ABANDONED":
									return (
										<Link
											href={getResultsUrl(
												packageSlug,
												exam.slug,
												exam.user_attempt?.session_id
											)}
											className="block w-full"
											onClick={(e) => e.stopPropagation()}
										>
											<button className="w-full px-4 py-3 bg-blue-500 text-white rounded-lg font-semibold hover:bg-blue-600 transition-colors">
												Review Answers
											</button>
										</Link>
									);
							}
						}

						// Priority 2: Check computed status for non-attempted exams
						switch (exam.computed_status) {
							case "UPCOMING":
								return (
									<div className="w-full px-4 py-3 bg-yellow-50 text-yellow-700 rounded-lg text-center font-medium border border-yellow-200">
										Exam Not Started Yet
									</div>
								);
							case "LIVE":
								return (
									<button
										onClick={() =>
											(window.location.href = `/exams/${packageSlug}/${exam.slug}`)
										}
										className="w-full px-4 py-3 bg-red-500 text-white rounded-lg font-semibold hover:bg-red-600 transition-colors shadow-sm"
									>
										Join Live Exam Now
									</button>
								);
							case "COMPLETED":
								return (
									<button
										onClick={() =>
											(window.location.href = `/exams/${packageSlug}/${exam.slug}`)
										}
										className="w-full px-4 py-3 bg-blue-500 text-white rounded-lg font-semibold hover:bg-blue-600 transition-colors"
									>
										Start Exam
									</button>
								);
							case "AVAILABLE":
							default:
								return (
									<Link
										href={`/exams/${packageSlug}/${exam.slug}`}
										className="w-full px-4 py-3 bg-blue-500 text-white rounded-lg font-semibold hover:bg-blue-600 transition-colors shadow-sm"
									>
										Start Exam
									</Link>
								);
						}
					})()}
				</div>

				{/* Footer Stats */}
				<div className="pt-3 border-t border-gray-200">
					<div className="flex items-center justify-between text-xs text-gray-500">
						<div className="flex items-center gap-4">
							<span className="flex items-center gap-1">
								<svg
									className="w-3 h-3"
									fill="currentColor"
									viewBox="0 0 20 20"
								>
									<path
										fillRule="evenodd"
										d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
										clipRule="evenodd"
									/>
								</svg>
								{exam.attempt_count} attempts
							</span>
						</div>
						<div className="flex items-center gap-3">
							<span>Avg: {exam.average_score?.toFixed(1)}%</span>
							<span>
								Pass Rate: {exam.pass_rate?.toFixed(1)}%
							</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
