import { ExamResponse } from "@/lib/api/types";
import { cn } from "@/styles/design-tokens";
import { Star, Trophy, Users } from "lucide-react";
import { hasParticipation } from "./examUtils";

interface ExamAnalyticsProps {
	exam: ExamResponse;
}

export function ExamAnalytics({ exam }: ExamAnalyticsProps) {
	if (!hasParticipation(exam)) {
		return <NoParticipationMessage />;
	}

	return (
		<div className="bg-white rounded-lg p-4 border">
			<h4 className="font-medium text-gray-900 mb-4 flex items-center gap-2">
				<Trophy size={16} className="text-blue-500" />
				Performance Analytics
			</h4>

			<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
				{/* Left Column - Attempt Stats */}
				<div className="space-y-3">
					<AnalyticCard
						icon={Users}
						label="Total Attempts"
						value={exam.attempt_count.toLocaleString()}
						color="blue"
					/>
					<AnalyticCard
						icon={Users}
						label="Completed Attempts"
						value={exam.attempt_count.toLocaleString()}
						color="green"
					/>
				</div>

				{/* Right Column - Performance Stats */}
				<div className="space-y-3">
					<AnalyticCard
						icon={Star}
						label="Average Score"
						value={`${exam.average_score?.toFixed(1)}%`}
						color="amber"
					/>
					<AnalyticCard
						icon={Trophy}
						label="Pass Rate"
						value={`${exam.pass_rate?.toFixed(1)}%`}
						color={
							(exam.pass_rate || 0) >= 80
								? "emerald"
								: (exam.pass_rate || 0) >= 60
								? "amber"
								: "red"
						}
					/>
				</div>
			</div>

			{/* Progress Bar for Pass Rate */}
			<div className="mt-4 pt-3 border-t border-gray-100">
				<div className="flex items-center justify-between text-xs text-gray-600 mb-1">
					<span>Success Rate</span>
					<span>{exam.pass_rate?.toFixed(1)}%</span>
				</div>
				<div className="w-full bg-gray-200 rounded-full h-2">
					<div
						className={cn(
							"h-2 rounded-full transition-all duration-300",
							(exam.pass_rate || 0) >= 80
								? "bg-emerald-500"
								: (exam.pass_rate || 0) >= 60
								? "bg-amber-500"
								: "bg-red-500"
						)}
						style={{
							width: `${exam.pass_rate || 0}%`,
						}}
					/>
				</div>
			</div>
		</div>
	);
}

function NoParticipationMessage() {
	return (
		<div className="bg-white rounded-lg p-4 border">
			<h4 className="font-medium text-gray-900 mb-4 flex items-center gap-2">
				<Trophy size={16} className="text-blue-500" />
				Performance Analytics
			</h4>
			<div className="text-center py-8">
				<div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
					<Users size={24} className="text-gray-400" />
				</div>
				<h5 className="text-lg font-medium text-gray-900 mb-2">
					No Attempts Yet
				</h5>
				<p className="text-sm text-gray-600 mb-4">
					This exam is available but no one has attempted it yet.
				</p>
				<div className="inline-flex items-center gap-2 text-sm text-blue-600 bg-blue-50 px-3 py-1 rounded-full">
					<Trophy size={14} />
					<span>Be the first to attempt!</span>
				</div>
			</div>
		</div>
	);
}

interface AnalyticCardProps {
	icon: React.ComponentType<{ size: number; className?: string }>;
	label: string;
	value: string;
	color: "blue" | "green" | "amber" | "emerald" | "red";
}

function AnalyticCard({ icon: Icon, label, value, color }: AnalyticCardProps) {
	const colorClasses = {
		blue: "bg-blue-50 text-blue-600 text-blue-900 text-blue-700",
		green: "bg-green-50 text-green-600 text-green-900 text-green-700",
		amber: "bg-amber-50 text-amber-600 text-amber-900 text-amber-700",
		emerald:
			"bg-emerald-50 text-emerald-600 text-emerald-900 text-emerald-700",
		red: "bg-red-50 text-red-600 text-red-900 text-red-700",
	};

	return (
		<div
			className={`p-3 ${
				colorClasses[color].split(" ")[0]
			} rounded-lg flex items-center justify-between`}
		>
			<div className="flex items-center gap-2">
				<Icon size={16} className={colorClasses[color].split(" ")[1]} />
				<span
					className={`text-sm font-medium ${
						colorClasses[color].split(" ")[2]
					}`}
				>
					{label}
				</span>
			</div>
			<span
				className={`text-lg font-bold ${
					colorClasses[color].split(" ")[3]
				}`}
			>
				{value}
			</span>
		</div>
	);
}
