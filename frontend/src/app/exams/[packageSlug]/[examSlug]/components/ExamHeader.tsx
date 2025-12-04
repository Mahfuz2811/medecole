import { Clock } from "lucide-react";
import { ExamData } from "../types";

interface ExamHeaderProps {
	exam: ExamData;
	progress: number;
	answeredCount: number;
	totalQuestions: number;
	isReviewMode: boolean;
	timeRemaining: number;
}

export function ExamHeader({
	exam,
	progress,
	answeredCount,
	totalQuestions,
	isReviewMode,
	timeRemaining,
}: ExamHeaderProps) {
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

	const getTimerStyles = () => {
		if (isReviewMode) {
			return {
				bg: "bg-gray-100",
				text: "text-gray-600",
				icon: "text-gray-500",
			};
		}
		if (timeRemaining <= 300) {
			// Last 5 minutes
			return {
				bg: "bg-red-100",
				text: "text-red-700",
				icon: "text-red-600",
			};
		}
		if (timeRemaining <= 600) {
			// Last 10 minutes
			return {
				bg: "bg-orange-100",
				text: "text-orange-700",
				icon: "text-orange-600",
			};
		}
		return {
			bg: "bg-blue-100",
			text: "text-blue-700",
			icon: "text-blue-600",
		};
	};

	const timerStyles = getTimerStyles();

	return (
		<div className="px-4 sm:px-6 py-4 bg-gradient-to-r from-gray-50 to-white">
			<div className="flex items-start justify-between gap-4">
				{/* Left: Title and Progress */}
				<div className="flex-1 min-w-0">
					{/* Title */}
					<h1 className="text-base sm:text-lg font-bold text-gray-900 leading-tight truncate mb-2">
						{exam.title}
					</h1>

					{/* Progress Section - Below Title */}
					<div className="flex items-center gap-3">
						<div className="flex items-center gap-2 text-xs sm:text-sm">
							<span className="font-semibold text-gray-900">
								{answeredCount}/{totalQuestions}
							</span>
							<span className="text-gray-500 hidden sm:inline">
								({Math.round(progress)}% Complete)
							</span>
						</div>
						<div className="w-16 sm:w-24 bg-gray-200 rounded-full h-1.5 sm:h-2 overflow-hidden">
							<div
								className="bg-gradient-to-r from-blue-500 to-blue-600 h-full rounded-full transition-all duration-500 ease-out"
								style={{ width: `${progress}%` }}
							/>
						</div>
					</div>
				</div>

				{/* Right: Timer */}
				<div className="flex-shrink-0">
					<div
						className={`flex items-center gap-2 px-3 sm:px-4 py-2 rounded-xl ${timerStyles.bg} border border-white shadow-sm`}
					>
						<Clock
							className={`w-3 h-3 sm:w-4 sm:h-4 ${timerStyles.icon}`}
						/>
						<span
							className={`text-xs sm:text-sm font-mono font-bold ${timerStyles.text}`}
						>
							{formatTime(timeRemaining)}
						</span>
						{!isReviewMode && timeRemaining <= 300 && (
							<span className="text-xs font-bold text-red-600 animate-pulse hidden sm:inline">
								HURRY!
							</span>
						)}
					</div>
				</div>
			</div>
		</div>
	);
}
