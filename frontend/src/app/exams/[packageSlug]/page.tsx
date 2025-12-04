"use client";

import { BottomNav, StickyHeader } from "@/components/layout";
import { withAuth } from "@/contexts/AuthContext";
import { layouts } from "@/styles/design-tokens";
import { useParams } from "next/navigation";
import {
	ExamCard,
	ExamFilters,
	ExamPagination,
	PackageHeader,
} from "./components";
import { usePackageExams, type ExtendedExam } from "./hooks";

function PackageExamsPage() {
	const params = useParams();
	const packageSlug = params.packageSlug as string;

	// Use our custom hook for all exam logic
	const {
		packageData,
		allExams,
		filteredExams,
		paginatedExams,
		loading,
		error,
		totalPages,
		searchQuery,
		setSearchQuery,
		selectedType,
		setSelectedType,
		selectedStatus,
		setSelectedStatus,
		currentPage,
		setCurrentPage,
		itemsPerPage,
	} = usePackageExams(packageSlug);

	if (loading) {
		return (
			<main className="bg-blue-50 min-h-screen relative">
				<div className={layouts.stickyHeader}>
					<div className={`${layouts.container} bg-white shadow`}>
						<StickyHeader />
					</div>
				</div>
				<div className={`${layouts.container} ${layouts.pageContent}`}>
					<div className="flex items-center justify-center py-16">
						<div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
						<span className="ml-4 text-gray-600">
							Loading exams...
						</span>
					</div>
				</div>
			</main>
		);
	}

	if (error) {
		return (
			<main className="bg-blue-50 min-h-screen relative">
				<div className={layouts.stickyHeader}>
					<div className={`${layouts.container} bg-white shadow`}>
						<StickyHeader />
					</div>
				</div>
				<div className={`${layouts.container} ${layouts.pageContent}`}>
					<div className="flex flex-col items-center justify-center py-16">
						<div className="text-red-600 text-center">
							<svg
								className="w-16 h-16 mx-auto mb-4"
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
							<h2 className="text-xl font-semibold mb-2">
								Error Loading Exams
							</h2>
							<p className="text-gray-600">{error}</p>
						</div>
					</div>
				</div>
			</main>
		);
	}

	return (
		<main className="bg-blue-50 min-h-screen relative">
			<div className={layouts.stickyHeader}>
				<div className={`${layouts.container} bg-white shadow`}>
					<StickyHeader />
				</div>
			</div>

			<div className={`${layouts.container} ${layouts.pageContent}`}>
				<div className="space-y-6">
					{/* Package Header */}
					<PackageHeader
						packageData={packageData}
						loading={loading}
					/>

					{/* Search and Filters */}
					<ExamFilters
						searchQuery={searchQuery}
						setSearchQuery={setSearchQuery}
						selectedType={selectedType}
						setSelectedType={setSelectedType}
						selectedStatus={selectedStatus}
						setSelectedStatus={setSelectedStatus}
						filteredCount={filteredExams.length}
						totalCount={allExams.length}
						currentPage={currentPage}
						totalPages={totalPages}
						itemsPerPage={itemsPerPage}
					/>

					{/* Exams List */}
					<div className="space-y-4">
						<div className="flex items-center justify-between px-2 sm:px-0">
							<div className="flex items-center gap-3">
								<div className="flex items-center gap-2">
									<div className="w-8 h-8 bg-blue-500 rounded-lg flex items-center justify-center">
										<svg
											className="w-4 h-4 text-white"
											fill="none"
											stroke="currentColor"
											viewBox="0 0 24 24"
										>
											<path
												strokeLinecap="round"
												strokeLinejoin="round"
												strokeWidth={2}
												d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
											/>
										</svg>
									</div>
									<h2 className="text-xl font-bold text-gray-900">
										Exams
									</h2>
								</div>
								<div className="bg-blue-100 text-blue-700 px-3 py-1 rounded-full text-sm font-semibold">
									{filteredExams.length}{" "}
									{filteredExams.length === 1
										? "exam"
										: "exams"}
								</div>
								{(searchQuery ||
									selectedType !== "ALL" ||
									selectedStatus !== "ALL") && (
									<div className="hidden sm:flex items-center text-xs text-gray-500">
										<svg
											className="w-3 h-3 mr-1"
											fill="none"
											stroke="currentColor"
											viewBox="0 0 24 24"
										>
											<path
												strokeLinecap="round"
												strokeLinejoin="round"
												strokeWidth={2}
												d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.414A1 1 0 013 6.707V4z"
											/>
										</svg>
										filtered
									</div>
								)}
							</div>
						</div>

						{paginatedExams.length > 0 ? (
							<div className="space-y-4">
								{paginatedExams.map(
									(exam: ExtendedExam, index: number) => (
										<ExamCard
											key={exam.id || index}
											exam={exam}
										/>
									)
								)}
							</div>
						) : (
							<div className="bg-white shadow-sm p-12 text-center">
								<svg
									className="w-16 h-16 text-gray-400 mx-auto mb-4"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										strokeLinecap="round"
										strokeLinejoin="round"
										strokeWidth={2}
										d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
									/>
								</svg>
								<h3 className="text-lg font-medium text-gray-900 mb-2">
									No Exams Found
								</h3>
								<p className="text-gray-600">
									{searchQuery ||
									selectedType !== "ALL" ||
									selectedStatus !== "ALL"
										? "Try adjusting your search filters."
										: "No exams are available for this package yet."}
								</p>
							</div>
						)}
					</div>

					{/* Pagination */}
					{filteredExams.length > itemsPerPage && (
						<ExamPagination
							currentPage={currentPage}
							totalPages={totalPages}
							onPageChange={setCurrentPage}
							filteredCount={filteredExams.length}
							itemsPerPage={itemsPerPage}
						/>
					)}
				</div>
			</div>

			{/* Blue separator above BottomNav */}
			<div className={`${layouts.container} h-3 bg-blue-100`} />
			<BottomNav />
		</main>
	);
}

export default withAuth(PackageExamsPage);
