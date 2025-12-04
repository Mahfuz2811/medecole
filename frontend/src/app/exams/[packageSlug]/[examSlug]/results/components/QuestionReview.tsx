"use client";

import { CheckCircle, ChevronDown, ChevronUp, XCircle } from "lucide-react";
import { useRef, useState } from "react";

// Local interface definition to match enhanced API response
interface ExamResultOption {
	id: number;
	option_text: string;
	is_correct: boolean;
}

interface ExamResultQuestion {
	id: number;
	question: string;
	question_type: "SBA" | "TRUE_FALSE";
	options: ExamResultOption[]; // Enhanced to include option details
	correct_answer: number | boolean[]; // number for SBA (option index), boolean[] for TRUE_FALSE
	user_answer: number | boolean[] | (boolean | undefined)[] | null; // user's answer, can be null if not answered, undefined for unanswered TRUE_FALSE options
	is_correct: boolean;
	points: number; // Earned points
	max_points?: number; // Optional max points for the question
	time_spent?: number; // time spent on this question in seconds
	explanation?: string;
}

interface QuestionReviewProps {
	questions: ExamResultQuestion[];
}

export function QuestionReview({ questions }: QuestionReviewProps) {
	const [expandedQuestions, setExpandedQuestions] = useState<Set<number>>(
		new Set(questions.map((question) => question.id))
	);
	const [explanationVisible, setExplanationVisible] = useState<Set<number>>(
		new Set()
	);
	const questionRefs = useRef<{ [key: number]: HTMLDivElement | null }>({});

	const toggleQuestion = (questionId: number) => {
		const questionElement = questionRefs.current[questionId];

		setExpandedQuestions((prev) => {
			const newSet = new Set(prev);
			if (newSet.has(questionId)) {
				newSet.delete(questionId);
				// Also hide explanation when collapsing question
				setExplanationVisible((prevExp) => {
					const newExpSet = new Set(prevExp);
					newExpSet.delete(questionId);
					return newExpSet;
				});
			} else {
				newSet.add(questionId);
				// Gentle scroll to ensure question stays in comfortable view when expanding
				if (questionElement) {
					setTimeout(() => {
						const rect = questionElement.getBoundingClientRect();
						const isVisible =
							rect.top >= 0 &&
							rect.top <= window.innerHeight * 0.3;

						// Only scroll if question header is not comfortably visible
						if (!isVisible) {
							questionElement.scrollIntoView({
								behavior: "smooth",
								block: "nearest",
								inline: "nearest",
							});
						}
					}, 100); // Slightly longer delay to allow content to fully expand
				}
			}
			return newSet;
		});
	};

	const toggleExplanation = (questionId: number) => {
		setExplanationVisible((prev) => {
			const newSet = new Set(prev);
			if (newSet.has(questionId)) {
				newSet.delete(questionId);
			} else {
				newSet.add(questionId);
			}
			return newSet;
		});
	};

	return (
		<div className="space-y-3 sm:space-y-4">
			{questions.map((question, index) => (
				<div
					key={question.id}
					ref={(el) => {
						questionRefs.current[question.id] = el;
					}}
					className="bg-gray-50 rounded-lg border border-gray-200"
				>
					<div
						className="flex items-start justify-between cursor-pointer p-3 sm:p-4"
						onClick={() => toggleQuestion(question.id)}
					>
						<div className="flex items-start gap-2 sm:gap-3 flex-1 min-w-0">
							<div
								className={`w-6 h-6 sm:w-8 sm:h-8 rounded-full flex items-center justify-center text-xs sm:text-sm font-medium flex-shrink-0 ${
									question.is_correct
										? "bg-green-100 text-green-700"
										: "bg-red-100 text-red-700"
								}`}
							>
								{index + 1}
							</div>

							<div className="flex-1 min-w-0">
								<p className="text-sm sm:text-base font-medium text-gray-900 leading-snug mb-2">
									{question.question}
								</p>
								<div className="flex flex-wrap items-center gap-1.5 sm:gap-3">
									<div className="flex items-center gap-1">
										{question.is_correct ? (
											<CheckCircle
												size={14}
												className="text-green-600"
											/>
										) : (
											<XCircle
												size={14}
												className="text-red-600"
											/>
										)}
										<span
											className={`text-xs font-medium ${
												question.is_correct
													? "text-green-600"
													: "text-red-600"
											}`}
										>
											{question.is_correct
												? "Correct"
												: "Incorrect"}
										</span>
									</div>
									<span className="text-xs bg-gray-100 text-gray-700 px-2 py-0.5 rounded-full font-medium">
										{question.question_type}
									</span>
									<span className="text-xs text-gray-500 font-medium">
										{question.points}
										{question.max_points
											? `/${question.max_points}`
											: ""}{" "}
										pts
									</span>
								</div>
							</div>
						</div>

						<button className="ml-2 p-1 hover:bg-gray-100 rounded flex-shrink-0">
							{expandedQuestions.has(question.id) ? (
								<ChevronUp
									size={16}
									className="text-gray-400"
								/>
							) : (
								<ChevronDown
									size={16}
									className="text-gray-400"
								/>
							)}
						</button>
					</div>

					{expandedQuestions.has(question.id) && (
						<div className="px-3 sm:px-4 pb-3 sm:pb-4 border-t border-gray-200">
							<div className="mt-3 sm:mt-4">
								{question.question_type === "SBA" ? (
									<div className="space-y-2">
										{question.options.map(
											(
												option: ExamResultOption,
												optionIndex: number
											) => {
												const isCorrectAnswer =
													optionIndex ===
													question.correct_answer;
												const isUserAnswer =
													question.user_answer ===
													optionIndex;
												const isUserCorrect =
													isCorrectAnswer &&
													isUserAnswer;

												return (
													<div
														key={optionIndex}
														className={`p-2 sm:p-3 rounded-lg border text-sm transition-colors ${
															isCorrectAnswer
																? "bg-green-50 border-green-200 text-green-800"
																: isUserAnswer
																? "bg-red-50 border-red-200 text-red-800"
																: "bg-white border-gray-200 text-gray-700"
														}`}
													>
														<div className="flex items-start gap-3">
															<span className="text-sm font-medium text-gray-600 flex-shrink-0 mt-0.5">
																{String.fromCharCode(
																	65 +
																		optionIndex
																)}
																.
															</span>
															<div className="flex-1 min-w-0">
																<span className="block text-sm leading-relaxed">
																	{
																		option.option_text
																	}
																</span>
																{(isCorrectAnswer ||
																	isUserAnswer) && (
																	<div className="flex flex-wrap items-center gap-2 mt-2">
																		{isUserCorrect ? (
																			<span className="inline-flex items-center gap-1 text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full font-medium">
																				<svg
																					className="w-3 h-3"
																					fill="currentColor"
																					viewBox="0 0 20 20"
																				>
																					<path
																						fillRule="evenodd"
																						d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
																						clipRule="evenodd"
																					/>
																				</svg>
																				Correct
																				Answer
																			</span>
																		) : (
																			<>
																				{isCorrectAnswer && (
																					<span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full font-medium">
																						Correct
																						Answer
																					</span>
																				)}
																				{isUserAnswer &&
																					!isCorrectAnswer && (
																						<span className="inline-flex items-center gap-1 text-xs bg-red-100 text-red-800 px-2 py-1 rounded-full font-medium">
																							<svg
																								className="w-3 h-3"
																								fill="currentColor"
																								viewBox="0 0 20 20"
																							>
																								<path
																									fillRule="evenodd"
																									d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
																									clipRule="evenodd"
																								/>
																							</svg>
																							Your
																							Answer
																						</span>
																					)}
																			</>
																		)}
																	</div>
																)}
															</div>
														</div>
													</div>
												);
											}
										)}
									</div>
								) : (
									<div className="space-y-2">
										{question.options.map(
											(
												option: ExamResultOption,
												optionIndex: number
											) => {
												const correctAnswer =
													Array.isArray(
														question.correct_answer
													)
														? question
																.correct_answer[
																optionIndex
														  ]
														: false;
												const userAnswerArray =
													Array.isArray(
														question.user_answer
													)
														? (question.user_answer as (
																| boolean
																| undefined
														  )[])
														: null;
												const userAnswer =
													userAnswerArray
														? userAnswerArray[
																optionIndex
														  ]
														: undefined;
												const wasAnswered =
													userAnswer !== undefined;
												const isCorrect = wasAnswered
													? correctAnswer ===
													  userAnswer
													: false;

												return (
													<div
														key={optionIndex}
														className={`p-2 sm:p-3 rounded-lg border text-sm transition-colors ${
															isCorrect &&
															wasAnswered
																? "bg-green-50 border-green-200"
																: wasAnswered
																? "bg-red-50 border-red-200"
																: "bg-gray-50 border-gray-200"
														}`}
													>
														<div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
															<span className="flex-1 text-gray-800 min-w-0">
																{
																	option.option_text
																}
															</span>
															<div className="flex flex-wrap items-center gap-1 sm:gap-2 flex-shrink-0">
																<span
																	className={`text-xs px-2 py-0.5 rounded-full font-medium whitespace-nowrap ${
																		correctAnswer
																			? "bg-green-100 text-green-800"
																			: "bg-gray-100 text-gray-800"
																	}`}
																>
																	Correct:{" "}
																	{correctAnswer
																		? "TRUE"
																		: "FALSE"}
																</span>
																{wasAnswered && (
																	<span
																		className={`text-xs px-2 py-0.5 rounded-full font-medium whitespace-nowrap ${
																			isCorrect
																				? "bg-blue-100 text-blue-800"
																				: "bg-red-100 text-red-800"
																		}`}
																	>
																		Your
																		Answer:{" "}
																		{userAnswer
																			? "TRUE"
																			: "FALSE"}
																	</span>
																)}
																{wasAnswered &&
																	(isCorrect ? (
																		<CheckCircle
																			size={
																				14
																			}
																			className="text-green-600 flex-shrink-0"
																		/>
																	) : (
																		<XCircle
																			size={
																				14
																			}
																			className="text-red-600 flex-shrink-0"
																		/>
																	))}
															</div>
														</div>
													</div>
												);
											}
										)}
									</div>
								)}

								{/* Explanation Section */}
								{question.explanation && (
									<div className="mt-4 pt-4 border-t border-gray-200">
										<button
											onClick={() =>
												toggleExplanation(question.id)
											}
											className="flex items-center gap-2 text-sm font-medium text-blue-600 hover:text-blue-700 transition-colors mb-3 touch-target"
										>
											<svg
												className="w-4 h-4 flex-shrink-0"
												fill="none"
												stroke="currentColor"
												viewBox="0 0 24 24"
											>
												<path
													strokeLinecap="round"
													strokeLinejoin="round"
													strokeWidth={2}
													d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
												/>
											</svg>
											{explanationVisible.has(question.id)
												? "Hide Explanation"
												: "Show Explanation"}
										</button>

										{explanationVisible.has(
											question.id
										) && (
											<div className="bg-blue-50 border border-blue-200 rounded-lg p-3 sm:p-4">
												<h5 className="font-medium text-blue-900 mb-2 text-sm sm:text-base">
													Explanation
												</h5>
												<p className="text-sm text-blue-800 leading-relaxed">
													{question.explanation}
												</p>
											</div>
										)}
									</div>
								)}
							</div>
						</div>
					)}
				</div>
			))}
		</div>
	);
}
