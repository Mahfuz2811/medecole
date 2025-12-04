interface ExamPaginationProps {
	currentPage: number;
	totalPages: number;
	onPageChange: (page: number) => void;
	filteredCount: number;
	itemsPerPage: number;
}

export function ExamPagination({
	currentPage,
	totalPages,
	onPageChange,
	filteredCount,
	itemsPerPage,
}: ExamPaginationProps) {
	if (totalPages <= 1) return null;

	const getVisiblePages = () => {
		const delta = 2;
		const range = [];
		const rangeWithDots = [];

		for (
			let i = Math.max(2, currentPage - delta);
			i <= Math.min(totalPages - 1, currentPage + delta);
			i++
		) {
			range.push(i);
		}

		if (currentPage - delta > 2) {
			rangeWithDots.push(1, "...");
		} else {
			rangeWithDots.push(1);
		}

		rangeWithDots.push(...range);

		if (currentPage + delta < totalPages - 1) {
			rangeWithDots.push("...", totalPages);
		} else {
			if (totalPages > 1) {
				rangeWithDots.push(totalPages);
			}
		}

		return rangeWithDots;
	};

	const visiblePages = getVisiblePages();
	const startItem = (currentPage - 1) * itemsPerPage + 1;
	const endItem = Math.min(currentPage * itemsPerPage, filteredCount);

	return (
		<div className="bg-white rounded-lg shadow-sm p-6">
			<div className="flex flex-col sm:flex-row justify-between items-center gap-4">
				{/* Results info */}
				<div className="text-sm text-gray-700">
					Showing <span className="font-medium">{startItem}</span> to{" "}
					<span className="font-medium">{endItem}</span> of{" "}
					<span className="font-medium">{filteredCount}</span> results
				</div>

				{/* Pagination buttons */}
				<div className="flex items-center gap-2">
					{/* Previous button */}
					<button
						onClick={() => onPageChange(currentPage - 1)}
						disabled={currentPage === 1}
						className="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Previous
					</button>

					{/* Page numbers */}
					<div className="flex items-center gap-1">
						{visiblePages.map((page, index) => (
							<span key={index}>
								{page === "..." ? (
									<span className="px-3 py-2 text-sm text-gray-500">
										...
									</span>
								) : (
									<button
										onClick={() =>
											onPageChange(page as number)
										}
										className={`px-3 py-2 text-sm font-medium rounded-md ${
											currentPage === page
												? "bg-blue-600 text-white"
												: "text-gray-500 bg-white border border-gray-300 hover:bg-gray-50"
										}`}
									>
										{page}
									</button>
								)}
							</span>
						))}
					</div>

					{/* Next button */}
					<button
						onClick={() => onPageChange(currentPage + 1)}
						disabled={currentPage === totalPages}
						className="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Next
					</button>
				</div>
			</div>
		</div>
	);
}
