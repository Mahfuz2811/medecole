import { QuestionType } from "../types";

type FilterType = "all" | QuestionType;

interface QuestionFilterProps {
	currentFilter: FilterType;
	onFilterChange: (filter: FilterType) => void;
	questionCounts: {
		total: number;
		sba: number;
		trueFalse: number;
	};
}

export function QuestionFilter({
	currentFilter,
	onFilterChange,
	questionCounts,
}: QuestionFilterProps) {
	const handleFilterClick = (filterValue: QuestionType) => {
		// If the same filter is clicked, toggle it off (show all)
		if (currentFilter === filterValue) {
			onFilterChange("all");
		} else {
			onFilterChange(filterValue);
		}
	};

	const filterOptions = [
		{
			value: "SBA" as QuestionType,
			label: "SBA",
			count: questionCounts.sba,
			color: "bg-white text-emerald-700 border-emerald-200",
			activeColor: "bg-emerald-500 text-white border-emerald-500",
		},
		{
			value: "TRUE_FALSE" as QuestionType,
			label: "True/False",
			count: questionCounts.trueFalse,
			color: "bg-white text-purple-700 border-purple-200",
			activeColor: "bg-purple-500 text-white border-purple-500",
		},
	];

	return (
		<div className="bg-gradient-to-r from-blue-50 to-indigo-50 p-4">
			<div className="max-w-6xl mx-auto">
				<div className="flex items-center justify-center">
					<div className="flex gap-3">
						{filterOptions.map((option) => (
							<button
								key={option.value}
								onClick={() => handleFilterClick(option.value)}
								className={`px-4 py-2.5 text-sm font-semibold rounded-xl border-2 transition-all duration-200 shadow-sm hover:shadow-md transform hover:-translate-y-0.5 ${
									currentFilter === option.value
										? option.activeColor + " shadow-lg"
										: option.color +
										  " hover:bg-opacity-80 bg-white"
								}`}
							>
								{option.label}
								<span
									className={`ml-2 px-2 py-1 text-xs font-bold rounded-full ${
										currentFilter === option.value
											? "bg-white bg-opacity-25 text-white"
											: "bg-gray-100 text-gray-700"
									}`}
								>
									{option.count}
								</span>
							</button>
						))}
					</div>
				</div>
			</div>
		</div>
	);
}
