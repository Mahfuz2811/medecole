import { ExamResponse } from "@/lib/api/types";
import { cn } from "@/styles/design-tokens";
import { Calendar, ChevronDown, Clock, FileText, Timer } from "lucide-react";
import Link from "next/link";
import { ExamAnalytics } from "./ExamAnalytics";
import { isExamAvailable, shouldShowAnalytics } from "./examUtils";
import {
	formatDate,
	formatTime,
	getExamTypeColor,
	getExamTypeIcon,
} from "./utils";

interface ExamCardProps {
	exam: ExamResponse;
	isExpanded: boolean;
	onToggle: () => void;
	packageSlug: string;
}

export function ExamCard({
	exam,
	isExpanded,
	onToggle,
	packageSlug,
}: ExamCardProps) {
	const IconComponent = getExamTypeIcon(exam.exam_type);
	const colorClass = getExamTypeColor(exam.exam_type);
	const available = isExamAvailable(exam);

	return (
		<div
			className={cn(
				"border border-gray-200 rounded-lg overflow-hidden transition-all duration-200 hover:border-gray-300"
			)}
		>
			{/* Exam Header - Always Visible */}
			<button
				onClick={onToggle}
				className="w-full p-4 text-left hover:bg-gray-50/50 transition-colors duration-150"
			>
				<div className="flex items-center gap-4">
					{/* Exam Icon & Type */}
					<div className="flex-shrink-0">
						<div
							className={cn(
								"w-12 h-12 rounded-lg flex items-center justify-center",
								colorClass
							)}
						>
							<IconComponent size={20} />
						</div>
					</div>

					{/* Exam Summary */}
					<div className="flex-1 min-w-0">
						<div className="flex items-center justify-between">
							<div className="flex-1">
								<h3 className="font-semibold text-gray-900 mb-1">
									{exam.title}
								</h3>
								<div className="flex items-center gap-4 text-sm text-gray-600">
									<span className="flex items-center gap-1">
										<FileText size={14} />
										{exam.total_questions} Questions
									</span>
									<span className="flex items-center gap-1">
										<Timer size={14} />
										{exam.duration_minutes} Min
									</span>
									<span className="flex items-center gap-1">
										<Calendar size={14} />
										{exam.scheduled_start_date
											? formatDate(
													exam.scheduled_start_date
											  )
											: "Not scheduled"}
									</span>
								</div>
							</div>
							<div className="flex items-center gap-3">
								<ChevronDown
									size={20}
									className={cn(
										"text-gray-400 transition-transform duration-200",
										isExpanded ? "rotate-180" : ""
									)}
								/>
							</div>
						</div>
					</div>
				</div>
			</button>

			{/* Expanded Content */}
			{isExpanded && (
				<div className="border-t border-gray-100 p-4 bg-gray-50/30">
					<ExamExpandedContent
						exam={exam}
						available={available}
						packageSlug={packageSlug}
					/>
				</div>
			)}
		</div>
	);
}

function ExamExpandedContent({
	exam,
	available,
	packageSlug,
}: {
	exam: ExamResponse;
	available: boolean;
	packageSlug: string;
}) {
	return (
		<div className="space-y-4">
			{/* Description */}
			<p className="text-sm text-gray-600">
				{exam.description || "No description available"}
			</p>

			{/* Schedule Information */}
			<div className="bg-white rounded-lg p-4 border">
				<h4 className="font-medium text-gray-900 mb-2">Schedule</h4>
				<div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
					<div className="flex items-center gap-2">
						<Calendar size={16} className="text-blue-500" />
						<span className="text-gray-600">
							Start:{" "}
							{exam.scheduled_start_date
								? `${formatDate(
										exam.scheduled_start_date
								  )} at ${formatTime(
										exam.scheduled_start_date
								  )}`
								: "Not scheduled"}
						</span>
					</div>
					<div className="flex items-center gap-2">
						<Clock size={16} className="text-orange-500" />
						<span className="text-gray-600">
							End:{" "}
							{exam.scheduled_end_date
								? `${formatDate(
										exam.scheduled_end_date
								  )} at ${formatTime(exam.scheduled_end_date)}`
								: "Not scheduled"}
						</span>
					</div>
				</div>
			</div>

			{/* Analytics - Only show for available exams */}
			{shouldShowAnalytics(exam) && <ExamAnalytics exam={exam} />}

			{/* Action Button */}
			<div className="flex justify-end">
				{available ? (
					(() => {
						const userAttempt = exam.user_attempt;

						if (!userAttempt) {
							// No attempt exists - show Start Exam
							return (
								<button className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors">
									Start Exam
								</button>
							);
						}

						// User has an attempt - show button based on status
						switch (userAttempt.status) {
							case "STARTED":
								return (
									<button className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition-colors">
										Resume Exam
									</button>
								);
							case "COMPLETED":
							case "AUTO_SUBMITTED":
								return (
									<Link
										href={`/exams/${packageSlug}/${exam.slug}/results`}
									>
										<button className="bg-purple-600 text-white px-4 py-2 rounded-lg hover:bg-purple-700 transition-colors">
											View Results
										</button>
									</Link>
								);
							case "ABANDONED":
								return (
									<button className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors">
										Exam Questions
									</button>
								);
							default:
								return (
									<button className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors">
										Start Exam
									</button>
								);
						}
					})()
				) : (
					<button
						disabled
						className="bg-gray-100 text-gray-500 px-4 py-2 rounded-lg cursor-not-allowed"
					>
						Scheduled
					</button>
				)}
			</div>
		</div>
	);
}
