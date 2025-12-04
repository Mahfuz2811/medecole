"use client";

import { StickyHeader } from "@/components/layout";
import { useAuth } from "@/contexts/AuthContext";
import { usePackage } from "@/hooks/usePackages";
import { cn, layouts } from "@/styles/design-tokens";
import { CheckCircle, Trophy } from "lucide-react";
import { useState } from "react";

// Import refactored components
import {
	ExamScheduleSection,
	LoginModal,
	MobileTabNavigation,
	PackageDescription,
	PackageDetailSkeleton,
	PackageFeatures,
	PackageHero,
	PackageNotFound,
	PackagePurchaseCard,
	PackageStats,
	PackageValidityInfo,
	type MobileTab,
} from "@/components/package-detail";

interface PackageDetailClientProps {
	slug: string;
}

export function PackageDetailClient({ slug }: PackageDetailClientProps) {
	const { data: packageData, loading, error } = usePackage(slug);
	const { isAuthenticated } = useAuth();
	const [showLoginModal, setShowLoginModal] = useState(false);
	const [activeMobileTab, setActiveMobileTab] =
		useState<MobileTab>("overview");

	if (loading) {
		return <PackageDetailSkeleton />;
	}

	if (error || !packageData) {
		return <PackageNotFound error={error} />;
	}

	return (
		<div className="min-h-screen bg-gray-50">
			{/* Sticky Header: Full width on mobile, centered on desktop */}
			<div className={layouts.stickyHeader}>
				<div className={`${layouts.container} bg-white shadow`}>
					<StickyHeader />
				</div>
			</div>

			{/* Main Content */}
			<div className={`${layouts.container} ${layouts.pageContent}`}>
				{/* Mobile Tabs - Show only on mobile */}
				<div className="lg:hidden mb-6">
					<MobileTabNavigation
						activeTab={activeMobileTab}
						onTabChange={setActiveMobileTab}
						package={packageData}
					/>
				</div>

				{/* Desktop Layout - Hidden on mobile when tabs are active */}
				<div className="hidden lg:grid lg:grid-cols-3 gap-8">
					{/* Left Column - Package Details */}
					<div className="lg:col-span-2 space-y-6">
						<PackageHero package={packageData} />
						<PackageDescription package={packageData} />
						<PackageFeatures package={packageData} />
						<ExamScheduleSection package={packageData} />
					</div>

					{/* Right Column - Sidebar */}
					<div className="space-y-6">
						<PackagePurchaseCard
							package={packageData}
							isAuthenticated={isAuthenticated}
							onShowLogin={() => setShowLoginModal(true)}
						/>
						<PackageStats package={packageData} />
						<PackageValidityInfo package={packageData} />
					</div>
				</div>

				{/* Mobile Content Based on Active Tab */}
				<div className="lg:hidden">
					{activeMobileTab === "overview" && (
						<div className="space-y-6">
							<PackageHero package={packageData} />
							<PackageDescription package={packageData} />
							<PackageFeatures package={packageData} />
							<div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
								<PackageStats package={packageData} />
								<PackageValidityInfo package={packageData} />
							</div>
						</div>
					)}

					{activeMobileTab === "exams" && (
						<div className="space-y-6">
							<ExamScheduleSection package={packageData} />
						</div>
					)}

					{activeMobileTab === "enroll" && (
						<div className="space-y-6">
							<PackagePurchaseCard
								package={packageData}
								isAuthenticated={isAuthenticated}
								onShowLogin={() => setShowLoginModal(true)}
							/>
							<PackageValidityInfo package={packageData} />
						</div>
					)}

					{/* Floating Action Button - Show on overview and exams tabs */}
					{(activeMobileTab === "overview" ||
						activeMobileTab === "exams") && (
						<div className="fixed bottom-6 right-6 z-40">
							<button
								onClick={() => setActiveMobileTab("enroll")}
								className={cn(
									"w-14 h-14 rounded-full shadow-lg flex items-center justify-center transition-all duration-200 hover:scale-105",
									packageData.package_type === "FREE"
										? "bg-green-600 hover:bg-green-700 text-white"
										: "bg-blue-600 hover:bg-blue-700 text-white"
								)}
								aria-label="Quick Enroll"
							>
								{packageData.package_type === "FREE" ? (
									<CheckCircle size={24} />
								) : (
									<Trophy size={24} />
								)}
							</button>
						</div>
					)}
				</div>
			</div>

			{/* Login Modal */}
			{showLoginModal && (
				<LoginModal
					package={packageData}
					onClose={() => setShowLoginModal(false)}
				/>
			)}
		</div>
	);
}
