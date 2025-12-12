"use client";

import { useAuth } from "@/contexts/AuthContext";
import { usePackages } from "@/hooks/usePackages";
import type { PackageResponse } from "@/lib/api/types";
import { layouts } from "@/styles/design-tokens";
import Link from "next/link";

// New organized imports
import { PackageCard } from "@/components/features/packages";
import { BottomNav, StickyHeader } from "@/components/layout";

export default function CardsPage() {
	const { isAuthenticated, isLoading } = useAuth();
	const {
		data: packagesData,
		loading: packagesLoading,
		error: packagesError,
	} = usePackages({
		sort: "sort_order",
		order: "asc",
		per_page: 10,
	});

	if (isLoading || packagesLoading) {
		return (
			<div className="min-h-screen flex items-center justify-center bg-blue-50">
				<div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
			</div>
		);
	}

	return (
		<main className="min-h-screen relative">
			{/* Sticky Header: Full width on mobile, centered on desktop */}
			<div className={layouts.stickyHeader}>
				<div className={`${layouts.container} bg-white shadow`}>
					<StickyHeader />
				</div>
			</div>

			<div className={`${layouts.container} ${layouts.pageContent}`}>
				{/* Welcome Message */}
				<div className="bg-gradient-to-r from-blue-500 to-purple-600 p-6 text-white mb-6">
					<h2 className="text-xl font-bold mb-2">
						{isAuthenticated
							? "Welcome back to Medecole! 🎓"
							: "Welcome to Medecole! 🎓"}
					</h2>
					<p className="text-blue-100 mb-4">
						{isAuthenticated
							? "Continue your learning journey and practice with our comprehensive MCQ collection."
							: "Sign in to track your progress, save your scores, and unlock personalized features."}
					</p>
					{!isAuthenticated && (
						<Link
							href="/auth"
							className="inline-flex items-center px-4 py-2 bg-white text-blue-600 rounded-lg font-medium hover:bg-blue-50 transition-colors"
						>
							Sign In / Register
							<svg
								className="w-4 h-4 ml-2"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M9 5l7 7-7 7"
								/>
							</svg>
						</Link>
					)}
					{isAuthenticated && (
						<Link
							href="/dashboard"
							className="inline-flex items-center px-4 py-2 bg-white text-blue-600 rounded-lg font-medium hover:bg-blue-50 transition-colors"
						>
							Go to Dashboard
							<svg
								className="w-4 h-4 ml-2"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M9 5l7 7-7 7"
								/>
							</svg>
						</Link>
					)}
				</div>

				<div className="grid sm:grid-cols-2 gap-x-4 gap-y-4 p-4">
					{packagesError && (
						<div className="col-span-2 text-center text-red-600 bg-red-50 p-4">
							{packagesError}
						</div>
					)}
					{packagesData?.packages?.map((pkg: PackageResponse) => (
						<PackageCard
							key={`package-${pkg.slug}`}
							package={pkg}
						/>
					))}
					{!packagesData?.packages?.length && !packagesError && (
						<div className="col-span-2 text-center text-gray-500 bg-gray-50 p-4 rounded-lg">
							No packages available at the moment.
						</div>
					)}
				</div>
			</div>

			{/* Blue separator above BottomNav */}
			<div className={`${layouts.container} h-3 bg-blue-100`} />

			<BottomNav />
		</main>
	);
}
