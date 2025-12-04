interface ExamFiltersProps {
	searchQuery: string;
	setSearchQuery: (query: string) => void;
	selectedType: string;
	setSelectedType: (type: string) => void;
	selectedStatus: string;
	setSelectedStatus: (status: string) => void;
	filteredCount: number;
	totalCount: number;
	currentPage: number;
	totalPages: number;
	itemsPerPage: number;
}

export function ExamFilters({
	searchQuery,
	setSearchQuery,
	selectedType,
	setSelectedType,
	selectedStatus,
	setSelectedStatus,
	filteredCount,
	totalCount,
	currentPage,
	totalPages,
	itemsPerPage,
}: ExamFiltersProps) {
	return (
		<div className="bg-white shadow-sm p-6">
			<div className="flex flex-col lg:flex-row gap-4">
				{/* Search */}
				<div className="flex-1">
					<div className="relative">
						<input
							type="text"
							placeholder="Search exams..."
							value={searchQuery}
							onChange={(e) => setSearchQuery(e.target.value)}
							className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
						<svg
							className="absolute left-3 top-2.5 h-5 w-5 text-gray-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
							/>
						</svg>
					</div>
				</div>

				{/* Filters */}
				<div className="flex gap-3">
					<select
						value={selectedType}
						onChange={(e) => setSelectedType(e.target.value)}
						className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					>
						<option value="ALL">All Types</option>
						<option value="DAILY">Daily</option>
						<option value="MOCK">Mock</option>
						<option value="REVIEW">Review</option>
						<option value="FINAL">Final</option>
					</select>

					<select
						value={selectedStatus}
						onChange={(e) => setSelectedStatus(e.target.value)}
						className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					>
						<option value="ALL">All Status</option>
						<option value="AVAILABLE">Available</option>
						<option value="COMPLETED">Completed</option>
					</select>
				</div>
			</div>

			{/* Results Info */}
			<div className="mt-4 flex justify-between items-center text-sm text-gray-600">
				<span>
					{filteredCount} of {totalCount} exams
					{searchQuery && ` matching "${searchQuery}"`}
				</span>
				{filteredCount > itemsPerPage && (
					<span>
						Page {currentPage} of {totalPages}
					</span>
				)}
			</div>
		</div>
	);
}
