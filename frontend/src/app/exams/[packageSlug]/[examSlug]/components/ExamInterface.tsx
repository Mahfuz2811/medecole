"use client";

import type {
	ExamNotAvailableError,
	MaxAttemptsExceededError,
} from "@/lib/api/types";
import { useMemo, useState } from "react";
import {
	ErrorScreen,
	ExamAlreadySubmittedScreen,
	ExamHeader,
	ExamInstructions,
	ExamNavigationList,
	ExamNotAvailableScreen,
	ExamQuestionListItem,
	ExamReview,
	LoadingScreen,
	MaxAttemptsExceededScreen,
	QuestionFilter,
} from ".";
import { useExamState } from "../hooks/useExamState";
import { QuestionType } from "../types";
import { ExamSuccess } from "./ExamSuccess";

interface ExamInterfaceProps {
	packageSlug: string;
	examSlug: string;
}

type FilterType = "all" | QuestionType;

export function ExamInterface({ packageSlug, examSlug }: ExamInterfaceProps) {
	const examHook = useExamState(packageSlug, examSlug);
	const {
		examState,
		phase,
		progress,
		answeredCount,
		totalQuestions,
		canSubmit,
		examError, // New error state
		startExam,
		updateAnswer,
		reviewAnswers,
		backToExam,
		submitExam,
	} = examHook;

	const [questionFilter, setQuestionFilter] = useState<FilterType>("all");

	// Calculate question counts for filter
	const questionCounts = useMemo(() => {
		if (!examState.exam) {
			return { total: 0, sba: 0, trueFalse: 0 };
		}

		const sbaCount = examState.exam.questions.filter(
			(q) => q.question_type === "SBA"
		).length;
		const trueFalseCount = examState.exam.questions.filter(
			(q) => q.question_type === "TRUE_FALSE"
		).length;

		return {
			total: examState.exam.questions.length,
			sba: sbaCount,
			trueFalse: trueFalseCount,
		};
	}, [examState.exam]);

	// Filter questions based on selected filter
	const filteredQuestions = useMemo(() => {
		if (!examState.exam) return [];

		if (questionFilter === "all") {
			return examState.exam.questions;
		}

		return examState.exam.questions.filter(
			(question) => question.question_type === questionFilter
		);
	}, [examState.exam, questionFilter]);

	// Auto-submit when time runs out
	if (
		examState.timeRemaining === 0 &&
		phase === "exam" &&
		!examState.showResults
	) {
		submitExam();
	}

	if (phase === "loading") {
		return <LoadingScreen />;
	}

	if (phase === "error") {
		// Handle specific API errors
		if (examError) {
			switch (examError.type) {
				case "ALREADY_SUBMITTED":
					return (
						<ExamAlreadySubmittedScreen
							error={examError.data as { error: string; message: string; data: { previous_attempt: { attempt_id: number; submitted_at: string; score: number; status: string }; can_retry: boolean } }}
							examTitle={examState.exam?.title}
							packageSlug={packageSlug}
						/>
					);
				case "MAX_ATTEMPTS":
					return (
						<MaxAttemptsExceededScreen
							error={examError.data as MaxAttemptsExceededError}
							packageSlug={packageSlug}
						/>
					);
				case "NOT_AVAILABLE":
					return (
						<ExamNotAvailableScreen
							error={examError.data as ExamNotAvailableError}
							packageSlug={packageSlug}
						/>
					);
				default:
					return <ErrorScreen error={examError.message} />;
			}
		}
		return <ErrorScreen error={examState.error || "Unknown error"} />;
	}

	if (phase === "instructions") {
		return (
			<ExamInstructions exam={examState.exam!} onStartExam={startExam} />
		);
	}

	if (phase === "results") {
		return <ExamSuccess examTitle={examState.exam!.title} />;
	}

	// Main exam interface (exam and review phases)
	return (
		<div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
			{/* Fixed Header */}
			<div className="fixed top-0 left-0 right-0 z-50 bg-white shadow-lg border-b border-gray-200">
				<div className="max-w-6xl mx-auto">
					<ExamHeader
						exam={examState.exam!}
						progress={progress}
						answeredCount={answeredCount}
						totalQuestions={totalQuestions}
						isReviewMode={phase === "review"}
						timeRemaining={examState.timeRemaining}
					/>
				</div>

				{/* Question Filter - Only show during exam phase */}
				{phase === "exam" && (
					<div className="border-t border-gray-100">
						<QuestionFilter
							currentFilter={questionFilter}
							onFilterChange={setQuestionFilter}
							questionCounts={questionCounts}
						/>
					</div>
				)}
			</div>

			{/* Main Content */}
			<div className={phase === "review" ? "pt-24 pb-24" : "pt-40 pb-24"}>
				<div className="max-w-6xl mx-auto sm:px-6">
					{phase === "review" ? (
						<ExamReview
							exam={examState.exam!}
							userAnswers={examState.userAnswers}
							onBackToExam={backToExam}
						/>
					) : (
						<div className="space-y-6">
							{/* Filter Info */}
							{questionFilter !== "all" && (
								<div className="bg-white rounded-lg border border-gray-200 p-4">
									<p className="text-sm text-gray-600">
										Showing{" "}
										<span className="font-semibold text-gray-900">
											{filteredQuestions.length}
										</span>{" "}
										{questionFilter === "SBA"
											? "SBA"
											: "True/False"}{" "}
										questions out of{" "}
										<span className="font-semibold text-gray-900">
											{totalQuestions}
										</span>{" "}
										total questions.
									</p>
								</div>
							)}

							{/* Questions List */}
							{filteredQuestions.map((question) => {
								// Calculate the actual question number (not filtered index)
								const actualQuestionNumber =
									examState.exam!.questions.findIndex(
										(q) => q.id === question.id
									) + 1;

								return (
									<ExamQuestionListItem
										key={question.id}
										question={question}
										questionNumber={actualQuestionNumber}
										currentAnswer={
											examState.userAnswers[question.id]
										}
										onUpdateAnswer={updateAnswer}
									/>
								);
							})}

							{/* Empty State */}
							{filteredQuestions.length === 0 && (
								<div className="bg-white rounded-lg border border-gray-200 p-8 text-center">
									<p className="text-gray-500">
										No questions found for the selected
										filter.
									</p>
								</div>
							)}
						</div>
					)}
				</div>
			</div>

			{/* Fixed Navigation */}
			<div className="fixed bottom-0 left-0 right-0 z-50 bg-white shadow-lg border-t">
				<ExamNavigationList
					canSubmit={canSubmit}
					isReviewMode={phase === "review"}
					onReview={reviewAnswers}
					onSubmit={submitExam}
					answeredCount={answeredCount}
				/>
			</div>
		</div>
	);
}
