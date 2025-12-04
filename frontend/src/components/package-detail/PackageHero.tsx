import { OptimizedImage } from "@/components/ui";
import { ExamType } from "@/lib/api/types";
import { cn, designTokens } from "@/styles/design-tokens";
import { BookOpen, BookOpenCheck, NotebookPen } from "lucide-react";
import { PackageDetailProps } from "./types";

export function PackageHero({ package: pkg }: PackageDetailProps) {
	// Calculate exam counts by type
	const examCounts = pkg.exams.reduce((acc, { exam }) => {
		acc[exam.exam_type] = (acc[exam.exam_type] || 0) + 1;
		return acc;
	}, {} as Record<ExamType, number>);

	return (
		<div className="bg-white shadow-sm overflow-hidden">
			{/* Package Image */}
			<div className="relative h-64 sm:h-80 lg:h-96">
				<OptimizedImage
					images={pkg.images}
					fill
					variant="desktop"
					objectFit="cover"
					priority
					sizes="(max-width: 768px) 100vw, (max-width: 1200px) 80vw, 60vw"
				/>
				{/* Package Type Badge */}
				<div className="absolute top-4 right-4 z-10">
					<span
						className={cn(
							"px-3 py-1 rounded-full text-sm font-medium",
							pkg.package_type === "FREE"
								? "bg-green-100 text-green-800"
								: "bg-blue-100 text-blue-800"
						)}
					>
						{pkg.package_type === "FREE"
							? "Free Package"
							: "Premium Package"}
					</span>
				</div>
			</div>

			{/* Package Info */}
			<div className="p-6">
				<h1
					className={cn(
						designTokens.typography.xl,
						"font-bold text-gray-900 mb-3"
					)}
				>
					{pkg.name}
				</h1>

				{/* Quick Stats */}
				<div className="flex flex-wrap gap-4 text-sm text-gray-600">
					{examCounts.DAILY && (
						<div className="flex items-center gap-1">
							<BookOpen size={16} className="text-blue-500" />
							<span>{examCounts.DAILY} Daily Practice</span>
						</div>
					)}
					{examCounts.REVIEW && (
						<div className="flex items-center gap-1">
							<BookOpenCheck
								size={16}
								className="text-blue-500"
							/>
							<span>{examCounts.REVIEW} Review Exam</span>
						</div>
					)}
					{examCounts.MOCK && (
						<div className="flex items-center gap-1">
							<NotebookPen size={16} className="text-blue-500" />
							<span>{examCounts.MOCK} Mock Exam</span>
						</div>
					)}
					{examCounts.FINAL && (
						<div className="flex items-center gap-1">
							<BookOpenCheck
								size={16}
								className="text-orange-500"
							/>
							<span>{examCounts.FINAL} Final Exam</span>
						</div>
					)}
				</div>
			</div>
		</div>
	);
}
