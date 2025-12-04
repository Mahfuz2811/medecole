import { cn, designTokens } from "@/styles/design-tokens";
import { BookOpen, Calendar, CheckCircle, Users } from "lucide-react";
import { PackageDetailProps } from "./types";

export function PackageStats({ package: pkg }: PackageDetailProps) {
	const stats = [
		{
			label: "Total Exams",
			value: pkg.total_exams,
			icon: BookOpen,
		},
		{
			label: "Students Enrolled",
			value: (pkg.enrollment_count || 0).toLocaleString(),
			icon: Users,
		},
	];

	return (
		<div className="bg-white shadow-sm p-6">
			<h3
				className={cn(
					designTokens.typography.base,
					"font-semibold text-gray-900 mb-4"
				)}
			>
				Package Statistics
			</h3>
			<div className="space-y-4">
				{stats.map((stat, index) => {
					const Icon = stat.icon;
					return (
						<div
							key={index}
							className="flex items-center justify-between"
						>
							<div className="flex items-center gap-2">
								<Icon size={16} className="text-gray-500" />
								<span className="text-sm text-gray-600">
									{stat.label}
								</span>
							</div>
							<span className="font-semibold text-gray-900">
								{stat.value}
							</span>
						</div>
					);
				})}
			</div>
		</div>
	);
}

export function PackageValidityInfo({ package: pkg }: PackageDetailProps) {
	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString("en-US", {
			year: "numeric",
			month: "long",
			day: "numeric",
		});
	};

	return (
		<div className="bg-white shadow-sm p-6">
			<h3
				className={cn(
					designTokens.typography.base,
					"font-semibold text-gray-900 mb-4"
				)}
			>
				Access Information
			</h3>
			<div className="space-y-3">
				<div className="flex items-center gap-2">
					<Calendar size={16} className="text-blue-500" />
					<div>
						<p className="text-sm font-medium text-gray-900">
							Validity
						</p>
						<p className="text-sm text-gray-600">
							{pkg.validity_type === "FIXED" && pkg.validity_date
								? `Until ${formatDate(pkg.validity_date)}`
								: pkg.validity_days
								? `${pkg.validity_days} days from activation`
								: "Unlimited access"}
						</p>
					</div>
				</div>

				<div className="flex items-center gap-2">
					<CheckCircle size={16} className="text-green-500" />
					<div>
						<p className="text-sm font-medium text-gray-900">
							Status
						</p>
						<p className="text-sm text-gray-600">
							{pkg.is_active ? "Active Package" : "Inactive"}
						</p>
					</div>
				</div>
			</div>
		</div>
	);
}
