import { MaxAttemptsExceededError } from "@/lib/api/types";
import { RefreshCw, XCircle } from "lucide-react";
import Link from "next/link";

interface MaxAttemptsExceededProps {
	error: MaxAttemptsExceededError;
	packageSlug: string;
}

export function MaxAttemptsExceededScreen({
	error,
	packageSlug,
}: MaxAttemptsExceededProps) {
	return (
		<div className="min-h-screen bg-blue-50 flex items-center justify-center">
			<div className="max-w-lg text-center p-6">
				<div className="flex items-center justify-center mb-6">
					<div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center">
						<XCircle className="w-8 h-8 text-red-600" />
					</div>
				</div>

				<h2 className="text-2xl font-bold text-gray-900 mb-3">
					Maximum Attempts Reached
				</h2>

				<p className="text-gray-600 mb-6">
					You have reached the maximum number of attempts for this
					exam.
				</p>

				<div className="bg-white rounded-lg p-4 mb-6 border border-gray-200">
					<div className="grid grid-cols-2 gap-4 text-sm">
						<div>
							<p className="text-gray-500">Attempts Used</p>
							<p className="font-medium text-gray-900">
								{error.attempted_count}
							</p>
						</div>
						<div>
							<p className="text-gray-500">Max Allowed</p>
							<p className="font-medium text-gray-900">
								{error.max_attempts}
							</p>
						</div>
					</div>
				</div>

				<div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
					<div className="flex items-start space-x-2">
						<RefreshCw className="w-5 h-5 text-yellow-600 mt-0.5" />
						<div className="text-left">
							<p className="text-sm font-medium text-yellow-800">
								Need More Attempts?
							</p>
							<p className="text-sm text-yellow-700 mt-1">
								Contact your instructor or exam administrator to
								reset your attempts if additional practice is
								needed.
							</p>
						</div>
					</div>
				</div>

				<div className="space-y-3">
					<Link
						href={`/exams/${packageSlug}`}
						className="block w-full bg-blue-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-blue-700 transition-colors"
					>
						Back to Exam List
					</Link>
					<Link
						href="/dashboard"
						className="block w-full bg-gray-100 text-gray-700 py-3 px-4 rounded-lg font-medium hover:bg-gray-200 transition-colors"
					>
						Dashboard
					</Link>
				</div>
			</div>
		</div>
	);
}
