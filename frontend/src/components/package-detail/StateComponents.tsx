import { ArrowLeft } from "lucide-react";
import Link from "next/link";

interface PackageNotFoundProps {
	error?: string | null;
}

export function PackageNotFound({ error }: PackageNotFoundProps) {
	return (
		<div className="min-h-screen bg-gray-50 flex items-center justify-center">
			<div className="text-center">
				<h1 className="text-2xl font-bold text-gray-900 mb-4">
					Package Not Found
				</h1>
				<p className="text-gray-600 mb-6">
					{error ||
						"The package you're looking for doesn't exist or has been removed."}
				</p>
				<Link
					href="/packages"
					className="inline-flex items-center gap-2 bg-blue-600 text-white px-4 py-2 hover:bg-blue-700 transition-colors"
				>
					<ArrowLeft size={16} />
					Back to Packages
				</Link>
			</div>
		</div>
	);
}

export function PackageDetailSkeleton() {
	return (
		<div className="min-h-screen bg-gray-50">
			<div className="bg-white border-b border-gray-200">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
					<div className="flex items-center gap-4 py-4">
						<div className="w-32 h-6 bg-gray-200 rounded animate-pulse" />
					</div>
				</div>
			</div>

			<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				<div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
					<div className="lg:col-span-2 space-y-6">
						{/* Hero Skeleton */}
						<div className="bg-white shadow-sm overflow-hidden">
							<div className="h-96 bg-gray-200 animate-pulse" />
							<div className="p-6 space-y-4">
								<div className="w-3/4 h-8 bg-gray-200 rounded animate-pulse" />
								<div className="flex gap-4">
									<div className="w-20 h-4 bg-gray-200 rounded animate-pulse" />
									<div className="w-24 h-4 bg-gray-200 rounded animate-pulse" />
									<div className="w-16 h-4 bg-gray-200 rounded animate-pulse" />
								</div>
							</div>
						</div>

						{/* Description Skeleton */}
						<div className="bg-white shadow-sm p-6">
							<div className="w-48 h-6 bg-gray-200 rounded animate-pulse mb-4" />
							<div className="space-y-3">
								<div className="w-full h-4 bg-gray-200 rounded animate-pulse" />
								<div className="w-full h-4 bg-gray-200 rounded animate-pulse" />
								<div className="w-3/4 h-4 bg-gray-200 rounded animate-pulse" />
							</div>
						</div>
					</div>

					{/* Sidebar Skeleton */}
					<div className="space-y-6">
						<div className="bg-white shadow-sm p-6">
							<div className="text-center mb-6">
								<div className="w-24 h-8 bg-gray-200 rounded animate-pulse mx-auto" />
							</div>
							<div className="w-full h-12 bg-gray-200 rounded animate-pulse" />
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
