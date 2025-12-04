import { useState } from "react";
import { Question, UserAnswer } from "../types";

interface ExamQuestionProps {
	question: Question;
	currentAnswer?: UserAnswer;
	onUpdateAnswer: (questionId: number, selectedOptions: string[]) => void;
	questionNumber: number;
	totalQuestions: number;
}

export function ExamQuestion({
	question,
	currentAnswer,
	onUpdateAnswer,
	questionNumber,
	totalQuestions,
}: ExamQuestionProps) {
	// For SBA questions, track selected option
	const [selectedOptions, setSelectedOptions] = useState<string[]>(() => {
		// Only use selectedOptions for SBA questions
		if (
			question.question_type === "SBA" &&
			currentAnswer?.selectedOptions
		) {
			return currentAnswer.selectedOptions;
		}
		return [];
	});

	// For TRUE_FALSE questions, we need to track true/false selections for each option
	const [trueFalseSelections, setTrueFalseSelections] = useState<
		Record<string, "true" | "false" | null>
	>(() => {
		if (
			question.question_type === "TRUE_FALSE" &&
			currentAnswer?.selectedOptions
		) {
			const selections: Record<string, "true" | "false" | null> = {};
			// Parse the selections from the stored format
			// Assuming stored format is ["a:true", "b:false", "c:true"] etc.
			currentAnswer.selectedOptions.forEach((selection) => {
				const [optionKey, value] = selection.split(":");
				if (value === "true" || value === "false") {
					selections[optionKey] = value as "true" | "false";
				}
			});
			return selections;
		}
		return {};
	});

	const handleOptionChange = (optionKey: string) => {
		// This function is only for SBA questions
		if (question.question_type === "SBA") {
			// Single Best Answer - only one option can be selected
			const newSelectedOptions = [optionKey];
			setSelectedOptions(newSelectedOptions);
			onUpdateAnswer(question.id, newSelectedOptions);
		}
	};

	const handleTrueFalseChange = (
		optionKey: string,
		value: "true" | "false"
	) => {
		const newSelections = {
			...trueFalseSelections,
			[optionKey]: value,
		};
		setTrueFalseSelections(newSelections);

		// Convert to the format expected by the backend
		const formattedSelections = Object.entries(newSelections)
			.filter(([, val]) => val !== null)
			.map(([key, val]) => `${key}:${val}`);

		onUpdateAnswer(question.id, formattedSelections);
	};

	const isOptionSelected = (optionKey: string) => {
		return selectedOptions.includes(optionKey);
	};

	const getTrueFalseSelection = (
		optionKey: string
	): "true" | "false" | null => {
		return trueFalseSelections[optionKey] || null;
	};

	return (
		<div className="bg-white rounded-lg shadow-sm border border-gray-200">
			{/* Question Header */}
			<div className="py-6 px-2 sm:px-6 border-b border-gray-200">
				<div className="flex items-start justify-between mb-4">
					<div className="flex items-center gap-3">
						<span className="inline-flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-600 rounded-full text-sm font-semibold">
							{questionNumber}
						</span>
						<div>
							<span className="text-sm text-gray-500">
								Question {questionNumber} of {totalQuestions}
							</span>
							<div className="flex items-center gap-2 mt-1">
								<span
									className={`text-xs px-2 py-1 rounded-full font-medium ${
										question.question_type === "SBA"
											? "bg-green-100 text-green-700"
											: "bg-purple-100 text-purple-700"
									}`}
								>
									{question.question_type === "SBA"
										? "SBA"
										: "True/False"}
								</span>
								<span className="text-xs text-gray-500">
									{question.points}{" "}
									{question.points === 1 ? "point" : "points"}
								</span>
							</div>
						</div>
					</div>
				</div>

				{/* Question Text */}
				<div className="prose prose-sm max-w-none">
					<p className="text-gray-900 font-medium leading-relaxed">
						{question.question_text}
					</p>
				</div>

				{/* Instructions */}
				<div className="mt-4 p-3 bg-blue-50 rounded-lg">
					<p className="text-sm text-blue-800">
						{question.question_type === "SBA"
							? "Select the ONE best answer."
							: "Mark each option as TRUE or FALSE."}
					</p>
				</div>
			</div>

			{/* Options */}
			<div className="py-6 px-2 sm:px-6">
				<div className="space-y-3">
					{Object.entries(question.options).map(
						([optionKey, option]) => {
							if (question.question_type === "SBA") {
								// SBA Question Layout - circular option letters
								const isSelected = isOptionSelected(optionKey);
								return (
									<label
										key={optionKey}
										className={`flex items-center gap-4 p-4 rounded-lg border-2 cursor-pointer transition-all duration-200 ${
											isSelected
												? "border-blue-500 bg-blue-50"
												: "border-gray-200 hover:border-gray-300 hover:bg-gray-50"
										}`}
									>
										{/* Hidden radio input for form behavior */}
										<input
											type="radio"
											name={`question-${question.id}`}
											value={optionKey}
											checked={isSelected}
											onChange={() =>
												handleOptionChange(optionKey)
											}
											className="sr-only"
										/>
										<div className="flex items-center gap-3 flex-1 min-w-0">
											{/* Circular option letter */}
											<span
												className={`inline-flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold border-2 transition-colors ${
													isSelected
														? "bg-blue-500 text-white border-blue-500"
														: "bg-gray-100 text-gray-700 border-gray-300"
												}`}
											>
												{optionKey.toUpperCase()}
											</span>
											<div className="flex-1 min-w-0">
												<p
													className={`text-sm leading-relaxed ${
														isSelected
															? "text-blue-900 font-medium"
															: "text-gray-700"
													}`}
												>
													{option.text}
												</p>
											</div>
										</div>
									</label>
								);
							} else {
								// TRUE_FALSE Question Layout
								const selection =
									getTrueFalseSelection(optionKey);
								return (
									<div
										key={optionKey}
										className="flex items-center justify-between p-4 rounded-lg border-2 border-gray-200 hover:border-gray-300 transition-all duration-200"
									>
										{/* Option Text (Left Side) */}
										<div className="flex-1 min-w-0 pr-4">
											<p className="text-sm leading-relaxed text-gray-700">
												<span className="font-semibold text-gray-700">
													{optionKey.toUpperCase()}.
												</span>{" "}
												{option.text}
											</p>
										</div>{" "}
										{/* TRUE/FALSE Radio Buttons (Right Side) */}
										<div className="flex items-center gap-4">
											{/* TRUE Button */}
											<label className="cursor-pointer">
												<div className="relative">
													<input
														type="radio"
														name={`question-${question.id}-option-${optionKey}`}
														value="true"
														checked={
															selection === "true"
														}
														onChange={() =>
															handleTrueFalseChange(
																optionKey,
																"true"
															)
														}
														className="sr-only"
													/>
													<div
														className={`w-8 h-8 rounded-full border-2 flex items-center justify-center text-xs font-bold transition-all duration-200 ${
															selection === "true"
																? "border-green-500 bg-green-500 text-white"
																: "border-gray-300 bg-white text-gray-400 hover:border-green-300"
														}`}
													>
														T
													</div>
												</div>
											</label>

											{/* FALSE Button */}
											<label className="cursor-pointer">
												<div className="relative">
													<input
														type="radio"
														name={`question-${question.id}-option-${optionKey}`}
														value="false"
														checked={
															selection ===
															"false"
														}
														onChange={() =>
															handleTrueFalseChange(
																optionKey,
																"false"
															)
														}
														className="sr-only"
													/>
													<div
														className={`w-8 h-8 rounded-full border-2 flex items-center justify-center text-xs font-bold transition-all duration-200 ${
															selection ===
															"false"
																? "border-red-500 bg-red-500 text-white"
																: "border-gray-300 bg-white text-gray-400 hover:border-red-300"
														}`}
													>
														F
													</div>
												</div>
											</label>
										</div>
									</div>
								);
							}
						}
					)}
				</div>

				{/* Answer Status */}
				<div className="mt-6 p-4 bg-gray-50 rounded-lg">
					<div className="flex items-center justify-between">
						<div className="flex items-center gap-2">
							{question.question_type === "SBA" ? (
								// SBA Status
								selectedOptions.length > 0 ? (
									<>
										<div className="w-2 h-2 bg-green-500 rounded-full"></div>
										<span className="text-sm text-gray-700">
											Answer saved:{" "}
											{selectedOptions
												.join(", ")
												.toUpperCase()}
										</span>
									</>
								) : (
									<>
										<div className="w-2 h-2 bg-gray-400 rounded-full"></div>
										<span className="text-sm text-gray-500">
											Not answered yet
										</span>
									</>
								)
							) : (
								// TRUE_FALSE Status
								(() => {
									const totalOptions = Object.keys(
										question.options
									).length;
									const answeredOptions = Object.values(
										trueFalseSelections
									).filter((val) => val !== null).length;
									const isComplete =
										answeredOptions === totalOptions;

									return isComplete ? (
										<>
											<div className="w-2 h-2 bg-green-500 rounded-full"></div>
											<span className="text-sm text-gray-700">
												All options answered (
												{answeredOptions}/{totalOptions}
												)
											</span>
										</>
									) : (
										<>
											<div className="w-2 h-2 bg-yellow-500 rounded-full"></div>
											<span className="text-sm text-gray-600">
												{answeredOptions > 0
													? `Partial: ${answeredOptions}/${totalOptions} answered`
													: "Not answered yet"}
											</span>
										</>
									);
								})()
							)}
						</div>
						{((question.question_type === "SBA" &&
							selectedOptions.length > 0) ||
							(question.question_type === "TRUE_FALSE" &&
								Object.values(trueFalseSelections).some(
									(val) => val !== null
								))) && (
							<button
								onClick={() => {
									if (question.question_type === "SBA") {
										setSelectedOptions([]);
										onUpdateAnswer(question.id, []);
									} else {
										setTrueFalseSelections({});
										onUpdateAnswer(question.id, []);
									}
								}}
								className="text-sm text-blue-600 hover:text-blue-700 font-medium"
							>
								Clear Answer
							</button>
						)}
					</div>
				</div>
			</div>
		</div>
	);
}
