import { AlertCircle } from "lucide-react";
import Link from "next/link";

interface ErrorScreenProps {
	error: string;
}

export function ErrorScreen({ error }: ErrorScreenProps) {
	return (
		<div className="min-h-screen bg-blue-50 flex items-center justify-center">
			<div className="max-w-md text-center">
				<div className="flex items-center justify-center mb-4">
					<AlertCircle className="w-12 h-12 text-red-500" />
				</div>
				<h2 className="text-xl font-semibold text-gray-900 mb-2">
					Unable to Load Exam
				</h2>
				<p className="text-gray-600 mb-6">{error}</p>
				<div className="space-y-3">
					<button
						onClick={() => window.location.reload()}
						className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-blue-700 transition-colors"
					>
						Try Again
					</button>
					<Link
						href="/dashboard"
						className="block w-full bg-gray-100 text-gray-700 py-3 px-4 rounded-lg font-medium hover:bg-gray-200 transition-colors"
					>
						Back to Dashboard
					</Link>
				</div>
			</div>
		</div>
	);
}
