import { PackageResponse } from "@/lib/api/types";

// Transform PackageResponse to include computed fields for UI
interface PackageData extends PackageResponse {
	completedExams: number;
	progress: number;
}

interface PackageHeaderProps {
	packageData: PackageData | null;
	loading: boolean;
}

export function PackageHeader({ packageData, loading }: PackageHeaderProps) {
	if (loading) {
		return (
			<div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg p-6">
				<div className="animate-pulse">
					<div className="flex items-center gap-4">
						<div className="w-16 h-16 bg-white/20 rounded-lg"></div>
						<div className="flex-1">
							<div className="h-8 bg-white/20 rounded w-2/3 mb-2"></div>
							<div className="h-4 bg-white/20 rounded w-full mb-2"></div>
							<div className="h-4 bg-white/20 rounded w-3/4"></div>
						</div>
					</div>
				</div>
			</div>
		);
	}

	if (!packageData) {
		return (
			<div className="bg-red-50 border border-red-200 p-6">
				<div className="text-red-800 text-center">
					<h2 className="text-xl font-semibold mb-2">
						Package Not Found
					</h2>
					<p>The requested package could not be found.</p>
				</div>
			</div>
		);
	}

	return (
		<div className="bg-white shadow-sm p-6 mb-6">
			<div className="flex justify-between items-start mb-4">
				<div className="flex-1">
					<div className="flex items-center gap-3 mb-2">
						<h1 className="text-2xl font-bold text-gray-900">
							{packageData.name}
						</h1>
						<span
							className={`px-3 py-1 rounded-full text-sm font-medium ${
								packageData.package_type === "FREE"
									? "bg-green-100 text-green-800"
									: "bg-blue-100 text-blue-800"
							}`}
						>
							{packageData.package_type}
						</span>
					</div>
				</div>
			</div>

			{/* Progress Overview */}
			<div className="border-t pt-4">
				<div className="flex justify-between items-center mb-2">
					<span className="text-sm font-medium text-gray-700">
						Overall Progress
					</span>
					<span className="text-sm text-gray-600">
						{packageData.completedExams}/{packageData.total_exams}{" "}
						Exams Complete
					</span>
				</div>
				<div className="w-full bg-gray-200 rounded-full h-3">
					<div
						className={`h-3 rounded-full transition-all duration-300 ${
							packageData.progress > 70
								? "bg-green-500"
								: packageData.progress > 30
								? "bg-yellow-500"
								: "bg-blue-500"
						}`}
						style={{ width: `${packageData.progress}%` }}
					></div>
				</div>
				<div className="text-right mt-1">
					<span className="text-sm font-medium text-gray-700">
						{packageData.progress.toFixed(1)}% Complete
					</span>
				</div>
			</div>
		</div>
	);
}
