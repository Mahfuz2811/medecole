import { Loader2 } from "lucide-react";

export function LoadingScreen() {
	return (
		<div className="min-h-screen bg-blue-50 flex items-center justify-center">
			<div className="text-center">
				<div className="flex items-center justify-center mb-4">
					<Loader2 className="w-8 h-8 text-blue-600 animate-spin" />
				</div>
				<h2 className="text-xl font-semibold text-gray-900 mb-2">
					Loading Exam
				</h2>
				<p className="text-gray-600">
					Please wait while we prepare your exam...
				</p>
			</div>
		</div>
	);
}
