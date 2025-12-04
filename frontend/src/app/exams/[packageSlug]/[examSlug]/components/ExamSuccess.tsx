import { CheckCircle, ExternalLink } from "lucide-react";
import Link from "next/link";

interface ExamSuccessProps {
	examTitle: string;
}

export function ExamSuccess({ examTitle }: ExamSuccessProps) {
	return (
		<div className="min-h-screen bg-gradient-to-br from-green-50 to-blue-50 flex items-center justify-center p-4">
			<div className="max-w-md w-full">
				{/* Success Card */}
				<div className="bg-white rounded-2xl shadow-xl p-8 text-center border border-gray-100">
					{/* Success Icon */}
					<div className="inline-flex items-center justify-center w-20 h-20 bg-green-100 rounded-full mb-6">
						<CheckCircle className="w-10 h-10 text-green-600" />
					</div>

					{/* Success Message */}
					<h1 className="text-2xl font-bold text-gray-900 mb-3">
						Exam Submitted Successfully!
					</h1>

					<p className="text-gray-600 mb-6 leading-relaxed">
						Your exam{" "}
						<span className="font-medium text-gray-900">
							&ldquo;{examTitle}&rdquo;
						</span>{" "}
						has been submitted successfully.
					</p>

					{/* Result Message */}
					<div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
						<p className="text-blue-800 text-sm font-medium mb-1">
							📊 Results Coming Soon
						</p>
						<p className="text-blue-700 text-xs">
							You will receive your detailed results shortly in
							your dashboard.
						</p>
					</div>

					{/* Action Buttons */}
					<div className="space-y-3">
						{/* Dashboard Link */}
						<Link
							href="/dashboard"
							className="w-full inline-flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-colors"
						>
							Go to Dashboard
							<ExternalLink className="w-4 h-4" />
						</Link>
					</div>
				</div>

				{/* Additional Info */}
				<div className="mt-6 text-center">
					<p className="text-xs text-gray-500">
						Results are typically available within 24 hours
					</p>
				</div>
			</div>
		</div>
	);
}
