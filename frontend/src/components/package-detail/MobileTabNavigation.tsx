import { cn } from "@/styles/design-tokens";
import { CheckCircle, Info, Trophy } from "lucide-react";
import { MobileTabNavigationProps } from "./types";

export function MobileTabNavigation({
	activeTab,
	onTabChange,
	package: pkg,
}: MobileTabNavigationProps) {
	const tabs = [
		{
			id: "overview" as const,
			label: "Overview",
			icon: Info,
			count: null,
		},
		{
			id: "exams" as const,
			label: "Exams",
			icon: Trophy,
			count: pkg.total_exams,
		},
		{
			id: "enroll" as const,
			label: pkg.package_type === "FREE" ? "Free Enroll" : "Enroll",
			icon: pkg.package_type === "FREE" ? CheckCircle : Trophy,
			count: null,
		},
	];

	return (
		<div className="bg-white shadow-sm rounded-lg overflow-hidden">
			{/* Quick Package Info Header */}
			<div className="p-4 bg-gradient-to-r from-blue-50 to-purple-50 border-b border-gray-200">
				<h1 className="font-bold text-lg text-gray-900 mb-1">
					{pkg.name}
				</h1>
				<div className="flex items-center justify-between">
					<span
						className={cn(
							"px-2 py-1 rounded-full text-xs font-medium",
							pkg.package_type === "FREE"
								? "bg-green-100 text-green-800"
								: "bg-blue-100 text-blue-800"
						)}
					>
						{pkg.package_type === "FREE"
							? "Free Package"
							: "Premium Package"}
					</span>
					{pkg.package_type === "PREMIUM" && (
						<span className="font-bold text-lg text-gray-900">
							৳{pkg.price.toLocaleString()}
						</span>
					)}
				</div>
			</div>

			{/* Tab Navigation */}
			<div className="flex">
				{tabs.map((tab) => {
					const Icon = tab.icon;
					const isActive = activeTab === tab.id;

					return (
						<button
							key={tab.id}
							onClick={() => onTabChange(tab.id)}
							className={cn(
								"flex-1 flex flex-col items-center justify-center py-4 px-2 border-b-2 transition-all duration-200",
								isActive
									? "border-blue-500 bg-blue-50 text-blue-700"
									: "border-transparent text-gray-600 hover:text-gray-900 hover:bg-gray-50"
							)}
						>
							<div className="flex items-center gap-1">
								<Icon size={18} />
								{tab.count !== null && (
									<span
										className={cn(
											"ml-1 px-1.5 py-0.5 rounded-full text-xs font-medium",
											isActive
												? "bg-blue-200 text-blue-800"
												: "bg-gray-200 text-gray-700"
										)}
									>
										{tab.count}
									</span>
								)}
							</div>
							<span className="text-xs font-medium mt-1">
								{tab.label}
							</span>
						</button>
					);
				})}
			</div>
		</div>
	);
}
