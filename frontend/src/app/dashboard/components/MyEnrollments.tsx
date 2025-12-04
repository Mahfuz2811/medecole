import { designTokens } from "@/styles/design-tokens";
import Link from "next/link";

interface Enrollment {
	id: number;
	packageId: number;
	packageName: string;
	packageSlug: string;
	packageType: string;
	progress: number;
	totalExams: number;
	completedExams: number;
	expiryDate?: string;
	status: string;
}

interface MyEnrollmentsProps {
	enrollments: Enrollment[];
	loading: boolean;
	error: string | null;
}

export function MyEnrollments({
	enrollments,
	loading,
	error,
}: MyEnrollmentsProps) {
	if (loading) {
		return (
			<div className="bg-white rounded-lg shadow-sm p-6 mb-6">
				<div className="flex justify-between items-center mb-4">
					<h2 className="text-lg font-semibold text-gray-900">
						My Packages
					</h2>
				</div>
				<div className="flex items-center justify-center py-8">
					<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
					<span className="ml-2 text-gray-600">
						Loading packages...
					</span>
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="bg-white rounded-lg shadow-sm p-6 mb-6">
				<div className="flex justify-between items-center mb-4">
					<h2 className="text-lg font-semibold text-gray-900">
						My Packages
					</h2>
				</div>
				<div className="text-center py-8">
					<div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg
							className="w-8 h-8 text-red-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
					</div>
					<h3 className="text-lg font-medium text-gray-900 mb-2">
						Failed to load packages
					</h3>
					<p className="text-gray-600 mb-4">{error}</p>
					<button
						onClick={() => window.location.reload()}
						className="text-sm text-blue-600 hover:text-blue-700 font-medium"
					>
						Try again
					</button>
				</div>
			</div>
		);
	}

	if (enrollments.length === 0) {
		return (
			<div className="bg-white rounded-lg shadow-sm p-6 mb-6">
				<div className="flex justify-between items-center mb-4">
					<h2 className="text-lg font-semibold text-gray-900">
						My Packages
					</h2>
				</div>
				<div className="text-center py-8">
					<div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg
							className="w-8 h-8 text-gray-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2M4 13h2m8-8v2m0 6v2"
							/>
						</svg>
					</div>
					<h3 className="text-lg font-medium text-gray-900 mb-2">
						No packages yet
					</h3>
					<p className="text-gray-600 mb-4">
						Start learning by enrolling in a package
					</p>
					<Link
						href="/"
						className={`${designTokens.components.button.primary} inline-flex`}
					>
						Browse Packages
					</Link>
				</div>
			</div>
		);
	}

	return (
		<div className="bg-white rounded-lg shadow-sm p-6 mb-6">
			<div className="flex justify-between items-center mb-4">
				<h2 className="text-lg font-semibold text-gray-900">
					My Packages
				</h2>
			</div>
			<div className="space-y-4">
				{enrollments.map((enrollment) => (
					<div
						key={enrollment.id}
						className="border border-gray-100 rounded-lg p-4 hover:border-blue-200 transition-colors"
					>
						<div className="flex justify-between items-start mb-3">
							<div>
								<h3 className="font-medium text-gray-900 mb-1">
									{enrollment.packageName}
								</h3>
								<div className="flex items-center gap-2">
									<span
										className={`px-2 py-1 rounded-full text-xs font-medium ${
											enrollment.packageType === "FREE"
												? "bg-green-100 text-green-800"
												: "bg-blue-100 text-blue-800"
										}`}
									>
										{enrollment.packageType}
									</span>
									{enrollment.expiryDate && (
										<span className="text-xs text-gray-500">
											Expires:{" "}
											{new Date(
												enrollment.expiryDate
											).toLocaleDateString()}
										</span>
									)}
								</div>
							</div>
							<div className="text-right">
								<p className="text-sm font-medium text-gray-900">
									{enrollment.completedExams}/
									{enrollment.totalExams} Exams
								</p>
								<p className="text-xs text-gray-500">
									{enrollment.progress}% Complete
								</p>
							</div>
						</div>
						{/* Progress Bar */}
						<div className="w-full bg-gray-200 rounded-full h-2 mb-2">
							<div
								className={`h-2 rounded-full transition-all duration-300 ${
									enrollment.progress > 70
										? "bg-green-500"
										: enrollment.progress > 30
										? "bg-yellow-500"
										: "bg-blue-500"
								}`}
								style={{
									width: `${enrollment.progress}%`,
								}}
							></div>
						</div>
						<div className="flex justify-between items-center">
							<span
								className={`text-xs px-2 py-1 rounded-full ${
									enrollment.status === "active"
										? "bg-green-50 text-green-700"
										: "bg-gray-50 text-gray-700"
								}`}
							>
								{enrollment.status === "active"
									? "Active"
									: "Enrolled"}
							</span>
							<Link
								href={`/exams/${enrollment.packageSlug}`}
								className="text-sm text-blue-600 hover:text-blue-700 font-medium"
							>
								Continue →
							</Link>
						</div>
					</div>
				))}
			</div>
		</div>
	);
}
