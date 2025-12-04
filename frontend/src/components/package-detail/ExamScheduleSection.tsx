import { PackageResponse } from "@/lib/api/types";
import { cn, designTokens } from "@/styles/design-tokens";
import { ChevronDown, ChevronUp } from "lucide-react";
import { useState } from "react";
import { ExamCard } from "./ExamCard";

interface ExamScheduleSectionProps {
	package: PackageResponse;
}

export function ExamScheduleSection({
	package: pkg,
}: ExamScheduleSectionProps) {
	// Get exams from the package data, sorted by sortOrder
	const packageExams =
		pkg.exams
			?.sort((a, b) => a.exam.sort_order - b.exam.sort_order)
			.map((pkgExam) => pkgExam.exam) || [];

	// Get a meaningful description from the first exam or use default
	const getScheduleDescription = () => {
		if (packageExams.length === 0) {
			return "No exams available for this package yet.";
		}

		// Try to find an exam with a description
		const examWithDescription = packageExams.find(
			(exam) => exam.description && exam.description.trim()
		);

		if (examWithDescription?.description) {
			return examWithDescription.description;
		}

		// Fallback to a simple count-based description
		return `${packageExams.length} exams available in this package with detailed schedules and information`;
	};

	// State for managing expanded exams
	const [expandedExams, setExpandedExams] = useState<Set<number>>(
		new Set(packageExams.slice(0, 2).map((exam) => exam.id))
	); // First 2 exams expanded by default
	const [showAll, setShowAll] = useState(false);

	// If no exams available, show empty state
	if (packageExams.length === 0) {
		return (
			<div className="bg-white shadow-sm">
				<div className="p-4 sm:p-6 border-b border-gray-200">
					<h2
						className={cn(
							designTokens.typography.lg,
							"font-semibold text-gray-900"
						)}
					>
						Exam Schedule & Syllabus
					</h2>
				</div>
				<div className="p-4 sm:p-6 text-center py-12">
					<p className="text-gray-500">
						No exams available for this package yet.
					</p>
				</div>
			</div>
		);
	}

	// Helper functions for exam expansion
	const toggleExam = (examId: number) => {
		const newExpanded = new Set(expandedExams);
		if (newExpanded.has(examId)) {
			newExpanded.delete(examId);
		} else {
			newExpanded.add(examId);
		}
		setExpandedExams(newExpanded);
	};

	const toggleAllExams = () => {
		if (expandedExams.size === packageExams.length) {
			// All expanded, collapse all
			setExpandedExams(new Set());
		} else {
			// Some or none expanded, expand all
			setExpandedExams(new Set(packageExams.map((exam) => exam.id)));
		}
	};

	const displayedExams = showAll ? packageExams : packageExams.slice(0, 6); // Show first 6 by default

	return (
		<div className="bg-white shadow-sm">
			<div className="p-4 sm:p-6 border-b border-gray-200">
				{/* Mobile: Stacked Layout */}
				<div className="block sm:hidden">
					<h2
						className={cn(
							designTokens.typography.lg,
							"font-semibold text-gray-900 mb-1"
						)}
					>
						Exam Schedule & Syllabus
					</h2>
					<p className="text-sm text-gray-600 mb-3">
						{packageExams.length} exams available
					</p>
					<div className="flex gap-2">
						<button
							onClick={toggleAllExams}
							className="flex-1 px-4 py-2 text-sm bg-blue-50 text-blue-600 hover:bg-blue-100 rounded-lg font-medium transition-colors"
						>
							{expandedExams.size === packageExams.length
								? "Collapse All"
								: "Expand All"}
						</button>
					</div>
					<p className="text-gray-600 text-sm mt-3">
						{getScheduleDescription()}
					</p>
				</div>

				{/* Desktop: Side-by-side Layout */}
				<div className="hidden sm:block">
					<div className="flex items-center justify-between mb-2">
						<h2
							className={cn(
								designTokens.typography.lg,
								"font-semibold text-gray-900"
							)}
						>
							Exam Schedule & Syllabus
						</h2>
						<div className="flex items-center gap-2">
							<button
								onClick={toggleAllExams}
								className="text-sm text-blue-600 hover:text-blue-700 font-medium"
							>
								{expandedExams.size === packageExams.length
									? "Collapse All"
									: "Expand All"}
							</button>
							<span className="text-gray-400">|</span>
							<span className="text-sm text-gray-600">
								{packageExams.length} Exams
							</span>
						</div>
					</div>
					<p className="text-gray-600 text-sm">
						{getScheduleDescription()}
					</p>
				</div>
			</div>

			<div className="p-4 sm:p-6">
				<div className="space-y-3">
					{displayedExams.map((exam) => (
						<ExamCard
							key={exam.id}
							exam={exam}
							isExpanded={expandedExams.has(exam.id)}
							onToggle={() => toggleExam(exam.id)}
							packageSlug={pkg.slug}
						/>
					))}
				</div>

				{/* Show More/Less Button */}
				{packageExams.length > 6 && (
					<div className="mt-6 text-center">
						<button
							onClick={() => setShowAll(!showAll)}
							className="text-blue-600 hover:text-blue-700 font-medium flex items-center gap-2 mx-auto px-4 py-2 rounded-lg hover:bg-blue-50 transition-colors"
						>
							{showAll ? (
								<>
									<ChevronUp size={16} />
									Show Less ({packageExams.length - 6} hidden)
								</>
							) : (
								<>
									<ChevronDown size={16} />
									Show More ({packageExams.length - 6} more
									exams)
								</>
							)}
						</button>
					</div>
				)}
			</div>
		</div>
	);
}
