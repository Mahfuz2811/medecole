import { cn, designTokens } from "@/styles/design-tokens";
import {
	BookOpen,
	CheckCircle,
	Clock,
	Star,
	Trophy,
	Users,
} from "lucide-react";
import { PackageDetailProps } from "./types";

export function PackageFeatures({ package: pkg }: PackageDetailProps) {
	const features = [
		{
			icon: BookOpen,
			title: "Comprehensive Exams",
			description: `${pkg.total_exams} carefully designed practice exams`,
		},
		{
			icon: Trophy,
			title: "Detailed Analytics",
			description:
				"Track your progress with detailed performance analytics",
		},
		{
			icon: Clock,
			title: "Flexible Schedule",
			description: "Study at your own pace with flexible exam scheduling",
		},
		{
			icon: CheckCircle,
			title: "Instant Results",
			description: "Get immediate feedback with detailed explanations",
		},
		{
			icon: Users,
			title: "Expert Support",
			description: "Access to expert instructors and community support",
		},
		{
			icon: Star,
			title: "Quality Content",
			description: "High-quality questions reviewed by subject experts",
		},
	];

	return (
		<div className="bg-white shadow-sm p-6">
			<h2
				className={cn(
					designTokens.typography.lg,
					"font-semibold text-gray-900 mb-6"
				)}
			>
				What&apos;s Included
			</h2>
			<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
				{features.map((feature, index) => {
					const Icon = feature.icon;
					return (
						<div key={index} className="flex gap-3">
							<div className="flex-shrink-0">
								<Icon
									size={20}
									className="text-blue-500 mt-1"
								/>
							</div>
							<div>
								<h3 className="font-medium text-gray-900 mb-1">
									{feature.title}
								</h3>
								<p className="text-sm text-gray-600">
									{feature.description}
								</p>
							</div>
						</div>
					);
				})}
			</div>
		</div>
	);
}
