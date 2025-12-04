import { ExamNotAvailableError } from "@/lib/api/types";
import { Calendar, Clock } from "lucide-react";
import Link from "next/link";

interface ExamNotAvailableProps {
	error: ExamNotAvailableError;
	packageSlug: string;
}

export function ExamNotAvailableScreen({
	error,
	packageSlug,
}: ExamNotAvailableProps) {
	const formatDate = (dateString?: string) => {
		if (!dateString) return null;
		const date = new Date(dateString);
		return {
			date: date.toLocaleDateString("en-US", { timeZone: "UTC" }),
			time: date.toLocaleTimeString("en-US", {
				hour: "2-digit",
				minute: "2-digit",
				timeZone: "UTC",
			}),
		};
	};

	const availableFrom = formatDate(error.available_from);
	const availableUntil = formatDate(error.available_until);

	return (
		<div className="min-h-screen bg-blue-50 flex items-center justify-center">
			<div className="max-w-lg text-center p-6">
				<div className="flex items-center justify-center mb-6">
					<div className="w-16 h-16 bg-yellow-100 rounded-full flex items-center justify-center">
						<Clock className="w-8 h-8 text-yellow-600" />
					</div>
				</div>

				<h2 className="text-2xl font-bold text-gray-900 mb-3">
					Exam Not Available
				</h2>

				<p className="text-gray-600 mb-6">{error.message}</p>

				{(availableFrom || availableUntil) && (
					<div className="bg-white rounded-lg p-4 mb-6 border border-gray-200">
						<div className="space-y-3">
							{availableFrom && (
								<div className="flex items-center justify-between">
									<div className="flex items-center space-x-2">
										<Calendar className="w-4 h-4 text-gray-500" />
										<span className="text-sm text-gray-500">
											Available From
										</span>
									</div>
									<div className="text-right">
										<p className="text-sm font-medium text-gray-900">
											{availableFrom.date}
										</p>
										<p className="text-xs text-gray-500">
											{availableFrom.time}
										</p>
									</div>
								</div>
							)}

							{availableUntil && (
								<div className="flex items-center justify-between">
									<div className="flex items-center space-x-2">
										<Calendar className="w-4 h-4 text-gray-500" />
										<span className="text-sm text-gray-500">
											Available Until
										</span>
									</div>
									<div className="text-right">
										<p className="text-sm font-medium text-gray-900">
											{availableUntil.date}
										</p>
										<p className="text-xs text-gray-500">
											{availableUntil.time}
										</p>
									</div>
								</div>
							)}
						</div>
					</div>
				)}

				<div className="space-y-3">
					<button
						onClick={() => window.location.reload()}
						className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-blue-700 transition-colors"
					>
						Check Again
					</button>
					<Link
						href={`/exams/${packageSlug}`}
						className="block w-full bg-gray-100 text-gray-700 py-3 px-4 rounded-lg font-medium hover:bg-gray-200 transition-colors"
					>
						Back to Exam List
					</Link>
					<Link
						href="/dashboard"
						className="block w-full text-gray-500 py-2 px-4 rounded-lg hover:text-gray-700 transition-colors"
					>
						Dashboard
					</Link>
				</div>
			</div>
		</div>
	);
}
