interface QuickStat {
	label: string;
	value: string;
	color: string;
	trend: string;
}

interface QuickStatsProps {
	stats: QuickStat[];
}

export function QuickStats({ stats }: QuickStatsProps) {
	return (
		<div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
			{stats.map((stat, index) => (
				<div
					key={index}
					className="bg-white rounded-lg shadow-sm p-4 relative overflow-hidden"
				>
					<div
						className={`w-2 h-2 ${stat.color} rounded-full mb-2`}
					></div>
					<p className="text-2xl font-bold text-gray-900 mb-1">
						{stat.value}
					</p>
					<p className="text-sm text-gray-600 mb-2">{stat.label}</p>
					<p className="text-xs text-green-600 font-medium">
						{stat.trend}
					</p>
				</div>
			))}
		</div>
	);
}
