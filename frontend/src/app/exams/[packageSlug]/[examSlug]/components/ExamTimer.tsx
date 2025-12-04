import { Clock } from "lucide-react";

interface ExamTimerProps {
	timeRemaining: number;
	isReviewMode: boolean;
}

export function ExamTimer({ timeRemaining, isReviewMode }: ExamTimerProps) {
	const formatTime = (seconds: number) => {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		const remainingSeconds = seconds % 60;

		if (hours > 0) {
			return `${hours}:${minutes
				.toString()
				.padStart(2, "0")}:${remainingSeconds
				.toString()
				.padStart(2, "0")}`;
		}
		return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
	};

	const getTimeColor = () => {
		if (isReviewMode) return "text-gray-600";
		if (timeRemaining <= 300) return "text-red-600"; // Last 5 minutes
		if (timeRemaining <= 600) return "text-orange-600"; // Last 10 minutes
		return "text-gray-900";
	};

	const getBackgroundColor = () => {
		if (isReviewMode) return "bg-gray-50";
		if (timeRemaining <= 300) return "bg-red-50"; // Last 5 minutes
		if (timeRemaining <= 600) return "bg-orange-50"; // Last 10 minutes
		return "bg-white";
	};

	if (isReviewMode) {
		return (
			<div className="flex items-center justify-center py-2">
				<div className="flex items-center gap-2 px-3 py-1 bg-gray-50 rounded-lg">
					<Clock className="w-4 h-4 text-gray-500" />
					<span className="text-sm font-medium text-gray-600">
						Review Mode
					</span>
				</div>
			</div>
		);
	}

	return (
		<div className="flex items-center justify-center py-2">
			<div
				className={`flex items-center gap-2 px-4 py-2 rounded-lg border ${getBackgroundColor()}`}
			>
				<Clock className={`w-4 h-4 ${getTimeColor()}`} />
				<span
					className={`text-sm font-mono font-semibold ${getTimeColor()}`}
				>
					{formatTime(timeRemaining)}
				</span>
				{timeRemaining <= 300 && (
					<span className="text-xs text-red-600 font-medium">
						Hurry!
					</span>
				)}
			</div>
		</div>
	);
}
