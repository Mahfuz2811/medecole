"use client";

import { BottomNav, StickyHeader } from "@/components/layout";
import { getExamResultsBySession, parseExamApiError } from "@/lib/api/exams";
import type { ExamResultQuestion } from "@/lib/api/types";
import { transformRawAnswersData } from "@/lib/utils/examResults";
import { layouts } from "@/styles/design-tokens";
import { Calendar, Clock, Target, Trophy } from "lucide-react";
import { useEffect, useState } from "react";
import { QuestionReview } from "./QuestionReview";

interface ExamResultInterfaceProps {
	packageSlug: string;
	examSlug: string;
	sessionId?: string | null;
}

interface TransformedResultData {
	exam: {
		id: number;
		title: string;
		description: string;
		total_questions: number;
		duration_minutes: number;
		passing_score: number;
		exam_type: string;
	};
	attempt: {
		id: number;
		status: string;
		started_at: string;
		completed_at: string;
		score: number;
		max_score: number;
		score_percentage: number;
		correct_answers: number;
		wrong_answers: number;
		is_passed: boolean;
		time_spent: number;
	};
	package: {
		name: string;
		slug: string;
	};
	questions: ExamResultQuestion[];
}

export function ExamResultInterface({
	packageSlug, // eslint-disable-line @typescript-eslint/no-unused-vars
	examSlug, // eslint-disable-line @typescript-eslint/no-unused-vars
	sessionId,
}: ExamResultInterfaceProps) {
	const [resultData, setResultData] = useState<TransformedResultData | null>(
		null
	);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	// Use provided sessionId or fall back to hardcoded session ID
	const SESSION_ID = sessionId || "0417c4c6e265a439bda9e201e7ee207f";

	useEffect(() => {
		const fetchResults = async () => {
			try {
				setLoading(true);
				setError(null);
				const rawData = await getExamResultsBySession(SESSION_ID);
				const transformedData = transformRawAnswersData(rawData);
				setResultData(transformedData);
			} catch (err) {
				const parsedError = parseExamApiError(err as Error);
				setError(parsedError.message);
				console.error("Failed to fetch exam results:", err);
			} finally {
				setLoading(false);
			}
		};

		fetchResults();
	}, [SESSION_ID]); // Remove dependency on packageSlug and examSlug since we're using session-based API

	const formatDate = (dateString: string) => {
		const date = new Date(dateString);
		return date.toLocaleDateString("en-US", {
			year: "numeric",
			month: "long",
			day: "numeric",
			hour: "2-digit",
			minute: "2-digit",
			timeZone: "UTC",
		});
	};

	if (loading) {
		return (
			<main className="bg-blue-50 min-h-screen relative pb-20">
				<div className={layouts.stickyHeader}>
					<div className={`${layouts.container} bg-white shadow`}>
						<StickyHeader />
					</div>
				</div>
				<div className={`${layouts.container} ${layouts.pageContent}`}>
					<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-8 text-center">
						<div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
						<p className="text-gray-600">Loading exam results...</p>
					</div>
				</div>
				<BottomNav />
			</main>
		);
	}

	if (error || !resultData) {
		return (
			<main className="bg-blue-50 min-h-screen relative pb-20">
				<div className={layouts.stickyHeader}>
					<div className={`${layouts.container} bg-white shadow`}>
						<StickyHeader />
					</div>
				</div>
				<div className={`${layouts.container} ${layouts.pageContent}`}>
					<div className="bg-white rounded-lg shadow-sm border border-gray-200 p-8 text-center">
						<div className="text-red-600 mb-4">
							<svg
								className="w-12 h-12 mx-auto mb-2"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z"
								/>
							</svg>
						</div>
						<h2 className="text-xl font-semibold text-gray-900 mb-2">
							Unable to Load Results
						</h2>
						<p className="text-gray-600">
							{error ||
								"Failed to load exam results. Please try again later."}
						</p>
					</div>
				</div>
				<BottomNav />
			</main>
		);
	}

	const { exam, attempt, questions } = resultData;

	return (
		<main className="bg-blue-50 min-h-screen relative pb-20">
			{/* Sticky Header */}
			<div className={layouts.stickyHeader}>
				<div className={`${layouts.container} bg-white shadow`}>
					<StickyHeader />
				</div>
			</div>

			{/* Content */}
			<div className={`${layouts.container} ${layouts.pageContent}`}>
				{/* Exam Header */}
				<div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden mb-6">
					{/* Result Header */}
					<div className="p-6 border-b border-gray-200">
						<div className="flex items-center justify-between mb-4">
							<div>
								<h1 className="text-2xl font-bold text-gray-900 mb-2">
									{exam.title}
								</h1>
								<p className="text-gray-600">
									{exam.description}
								</p>
							</div>
							<div
								className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
									attempt.is_passed
										? "bg-green-100 text-green-700"
										: "bg-red-100 text-red-700"
								}`}
							>
								{attempt.is_passed ? "PASSED" : "FAILED"}
							</div>
						</div>

						<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
							<div className="flex items-center gap-2 text-sm text-gray-600">
								<Calendar size={16} />
								<span>{formatDate(attempt.completed_at)}</span>
							</div>
							<div className="flex items-center gap-2 text-sm text-gray-600">
								<Clock size={16} />
								<span>
									Duration: {exam.duration_minutes} minutes
								</span>
							</div>
							<div className="flex items-center gap-2 text-sm text-gray-600">
								<Target size={16} />
								<span>
									Passing Score: {exam.passing_score}%
								</span>
							</div>
						</div>
					</div>

					{/* Score Overview */}
					<div className="p-6 border-b border-gray-200">
						<h3 className="text-lg font-semibold text-gray-900 mb-4">
							Your Performance
						</h3>
						<div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
							{/* Score Card */}
							<div
								className={`flex items-center gap-3 p-4 rounded-lg ${
									attempt.is_passed
										? "bg-green-50"
										: "bg-red-50"
								}`}
							>
								<Trophy
									className={`w-5 h-5 ${
										attempt.is_passed
											? "text-green-600"
											: "text-red-600"
									}`}
								/>
								<div>
									<p className="text-sm font-medium text-gray-900">
										Final Score
									</p>
									<p
										className={`text-lg font-bold ${
											attempt.is_passed
												? "text-green-600"
												: "text-red-600"
										}`}
									>
										{attempt.score}
										{attempt.max_score
											? `/${attempt.max_score}`
											: ""}
									</p>
									{attempt.score_percentage && (
										<p className="text-sm text-gray-600">
											{attempt.score_percentage.toFixed(
												1
											)}
											%
										</p>
									)}
								</div>
							</div>

							{/* Correct Answers */}
							<div className="flex items-center gap-3 p-4 bg-blue-50 rounded-lg">
								<Target className="w-5 h-5 text-blue-600" />
								<div>
									<p className="text-sm font-medium text-gray-900">
										Correct
									</p>
									<p className="text-lg font-bold text-blue-600">
										{attempt.correct_answers}/
										{exam.total_questions}
									</p>
								</div>
							</div>

							{/* Accuracy */}
							<div className="flex items-center gap-3 p-4 bg-purple-50 rounded-lg">
								<Calendar className="w-5 h-5 text-purple-600" />
								<div>
									<p className="text-sm font-medium text-gray-900">
										Accuracy
									</p>
									<p className="text-lg font-bold text-purple-600">
										{Math.round(
											(attempt.correct_answers /
												exam.total_questions) *
												100
										)}
										%
									</p>
								</div>
							</div>
						</div>
					</div>

					{/* Question Review Section */}
					<div className="p-6">
						<h3 className="text-lg font-semibold text-gray-900 mb-4">
							Question Review
						</h3>
						<QuestionReview questions={questions} />
					</div>
				</div>
			</div>

			{/* Bottom Navigation */}
			<BottomNav />
		</main>
	);
}
