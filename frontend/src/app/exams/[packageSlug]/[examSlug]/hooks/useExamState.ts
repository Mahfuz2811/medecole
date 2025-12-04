"use client";

import { examAPI, parseExamApiError } from "@/lib/api";
import type { ExamAlreadySubmittedError, ExamApiError } from "@/lib/api/types";
import { useCallback, useEffect, useRef, useState } from "react";
import { ExamPhase, ExamResult, ExamState, UserAnswer } from "../types";

// Global interface for storing exam result
declare global {
	interface Window {
		examResult?: ExamResult;
	}
}

interface ExamErrorState {
	type:
		| "ALREADY_SUBMITTED"
		| "MAX_ATTEMPTS"
		| "NOT_AVAILABLE"
		| "SESSION_EXPIRED"
		| "UNKNOWN";
	message: string;
	data?: ExamApiError;
}

export function useExamState(packageSlug: string, examSlug: string) {
	const [examState, setExamState] = useState<ExamState>({
		exam: null,
		session: null,
		loading: true,
		error: null,
		timeRemaining: 0,
		currentQuestion: null,
		currentQuestionIndex: 0,
		userAnswers: {},
		isReviewing: false,
		showResults: false,
	});

	const [phase, setPhase] = useState<ExamPhase>("loading");
	const [examError, setExamError] = useState<ExamErrorState | null>(null);
	const [sessionId, setSessionId] = useState<string | null>(null);
	const timerRef = useRef<NodeJS.Timeout | null>(null);
	const startTimeRef = useRef<Date | null>(null);

	// Helper function to get device info
	const getDeviceInfo = useCallback(() => {
		return {
			browser: navigator.userAgent.includes("Chrome")
				? "Chrome"
				: navigator.userAgent.includes("Firefox")
				? "Firefox"
				: navigator.userAgent.includes("Safari")
				? "Safari"
				: "Unknown",
			user_agent: navigator.userAgent,
			ip_address: "0.0.0.0", // Will be determined by backend
		};
	}, []);

	// Initialize exam metadata only (not session)
	useEffect(() => {
		const loadExamMetadata = async () => {
			try {
				setExamState((prev) => ({
					...prev,
					loading: true,
					error: null,
				}));
				setExamError(null);

				// Get exam metadata directly using the new endpoint
				const examMeta = await examAPI.getExamMeta(examSlug);

				// Create exam data with metadata only
				const examData = {
					id: examMeta.id,
					title: examMeta.title,
					slug: examMeta.slug,
					total_questions: examMeta.total_questions,
					duration_minutes: examMeta.duration_minutes,
					passing_score: examMeta.passing_score,
					total_marks: examMeta.total_marks,
					instructions: examMeta.instructions,
					questions: [], // No questions yet
				};
				setExamState((prev) => ({
					...prev,
					exam: examData,
					timeRemaining: examData.duration_minutes * 60,
					loading: false,
				}));

				setPhase("instructions");
			} catch (error) {
				console.error("Failed to load exam metadata:", error);

				const parsedError = parseExamApiError(error as Error);

				setExamError({
					type: "UNKNOWN",
					message: parsedError.message,
					data: parsedError.data,
				});

				setExamState((prev) => ({
					...prev,
					error: parsedError.message,
					loading: false,
				}));
				setPhase("error");
			}
		};

		loadExamMetadata();
	}, [packageSlug, examSlug]);

	// Timer management
	const startTimer = useCallback(() => {
		if (timerRef.current) {
			clearInterval(timerRef.current);
		}

		startTimeRef.current = new Date();

		timerRef.current = setInterval(() => {
			setExamState((prev: ExamState) => {
				const newTimeRemaining = Math.max(0, prev.timeRemaining - 1);

				// Auto-submit when time runs out
				if (newTimeRemaining === 0 && !prev.showResults) {
					// Use setTimeout to avoid dependency issues
					setTimeout(async () => {
						if (sessionId) {
							try {
								await examAPI.submitExam({
									session_id: sessionId,
								});
							} catch (error) {
								console.error("Auto-submit failed:", error);
								// Don't retry auto-submit on failure - just log the error
								// This prevents infinite loops when submission fails
							}
						}
					}, 100);
				}

				return {
					...prev,
					timeRemaining: newTimeRemaining,
				};
			});
		}, 1000);
	}, [sessionId]);

	const stopTimer = useCallback(() => {
		if (timerRef.current) {
			clearInterval(timerRef.current);
			timerRef.current = null;
		}
	}, []);

	// Navigation functions
	const goToQuestion = useCallback((index: number) => {
		setExamState((prev: ExamState) => {
			if (
				!prev.exam ||
				index < 0 ||
				index >= prev.exam.questions.length
			) {
				return prev;
			}

			return {
				...prev,
				currentQuestionIndex: index,
				currentQuestion: prev.exam.questions[index],
			};
		});
	}, []);

	const nextQuestion = useCallback(() => {
		setExamState((prev: ExamState) => {
			if (
				!prev.exam ||
				prev.currentQuestionIndex >= prev.exam.questions.length - 1
			) {
				return prev;
			}

			const newIndex = prev.currentQuestionIndex + 1;
			return {
				...prev,
				currentQuestionIndex: newIndex,
				currentQuestion: prev.exam.questions[newIndex],
			};
		});
	}, []);

	const previousQuestion = useCallback(() => {
		setExamState((prev: ExamState) => {
			if (prev.currentQuestionIndex <= 0) {
				return prev;
			}

			const newIndex = prev.currentQuestionIndex - 1;
			return {
				...prev,
				currentQuestionIndex: newIndex,
				currentQuestion: prev.exam?.questions[newIndex] || null,
			};
		});
	}, []);

	// Answer management with API sync
	const updateAnswer = useCallback(
		async (questionId: number, selectedOptions: string[]) => {
			const answer: UserAnswer = {
				questionId,
				selectedOptions,
				timeSpent: 0, // Calculate based on actual time
				answeredAt: new Date(),
				isSkipped: selectedOptions.length === 0,
			};

			// First, calculate what the updated answers will be
			const currentUserAnswers = examState.userAnswers;
			const updatedUserAnswers = {
				...currentUserAnswers,
				[questionId]: answer,
			};

			// Update state
			setExamState((prev: ExamState) => ({
				...prev,
				userAnswers: updatedUserAnswers,
			}));

			// Sync ALL answers to backend immediately (including previous ones)
			if (sessionId) {
				try {
					// Collect all answered questions (including the one just answered)
					const allAnswers = Object.values(updatedUserAnswers)
						.filter(
							(ans) =>
								!ans.isSkipped && ans.selectedOptions.length > 0
						)
						.map((ans) => ({
							question_id: ans.questionId,
							selected_option:
								ans.selectedOptions.length === 1
									? ans.selectedOptions[0] // Single option for SBA
									: JSON.stringify(ans.selectedOptions), // Multiple options for TRUE_FALSE as JSON array
						}));

					if (allAnswers.length > 0) {
						await examAPI.syncSession(sessionId, {
							answers: allAnswers,
						});
					}
				} catch (error) {
					console.warn("Failed to sync answers:", error);
				}
			}
		},
		[sessionId, examState.userAnswers]
	);

	// Exam flow functions
	const startExam = useCallback(async () => {
		try {
			setExamState((prev) => ({ ...prev, loading: true }));

			// Start the exam session
			const startResponse = await examAPI.startExam(examSlug, {
				package_slug: packageSlug,
				device_info: getDeviceInfo(),
			});

			setSessionId(startResponse.session_id);

			// Fetch session data including questions
			const sessionResponse = await examAPI.getSession(
				startResponse.session_id
			);

			// Transform questions to frontend format
			const transformedQuestions = sessionResponse.exam.questions.map(
				(q) => ({
					id: q.id,
					question_text: q.question_text,
					question_type: (q.question_type === "SINGLE_CHOICE"
						? "SBA"
						: q.question_type === "TRUE_FALSE"
						? "TRUE_FALSE"
						: "SBA") as "SBA" | "TRUE_FALSE",
					options: Object.entries(q.options || {}).reduce(
						(acc, [key, option]) => {
							acc[key] = {
								text: option.text || "",
								is_correct: false, // Not provided in secure response
								explanation: undefined,
							};
							return acc;
						},
						{} as Record<
							string,
							{
								text: string;
								is_correct: boolean;
								explanation?: string;
							}
						>
					),
					explanation: undefined, // Not provided in secure response
					points: q.points,
				})
			);

			// Restore saved answers if available (from Redis cache)
			const restoredUserAnswers: Record<number, UserAnswer> = {};
			if (
				sessionResponse.session.saved_answers &&
				sessionResponse.session.saved_answers.length > 0
			) {
				sessionResponse.session.saved_answers.forEach((savedAnswer) => {
					// Handle both single options and JSON arrays (for TRUE_FALSE questions)
					let selectedOptions: string[];
					try {
						// Try to parse as JSON array first (for TRUE_FALSE multiple selections)
						const parsed = JSON.parse(savedAnswer.selected_option);
						selectedOptions = Array.isArray(parsed)
							? parsed
							: [savedAnswer.selected_option];
					} catch {
						// If parsing fails, treat as single option (for SBA questions)
						selectedOptions = [savedAnswer.selected_option];
					}

					restoredUserAnswers[savedAnswer.question_id] = {
						questionId: savedAnswer.question_id,
						selectedOptions,
						timeSpent: 0, // We don't track individual time spent in cache yet
						answeredAt: new Date(), // Use current time as placeholder
						isSkipped: false,
					};
				});
				console.log(
					`Restored ${
						Object.keys(restoredUserAnswers).length
					} saved answers from session`
				);
			}

			// Update exam state with questions and restored answers
			setExamState((prev) => ({
				...prev,
				exam: prev.exam
					? {
							...prev.exam,
							questions: transformedQuestions,
					  }
					: null,
				currentQuestion: transformedQuestions[0],
				userAnswers: restoredUserAnswers, // Restore saved answers
				timeRemaining: sessionResponse.session.time_remaining, // Use actual remaining time from backend
				loading: false,
			}));

			setPhase("exam");
			startTimer();
		} catch (error) {
			console.error("Failed to start exam:", error);

			const parsedError = parseExamApiError(error as Error);

			if (parsedError.type === "EXAM_CONFLICT") {
				const examData = parsedError.data as ExamAlreadySubmittedError;
				setExamError({
					type: "ALREADY_SUBMITTED",
					message: examData.message,
					data: examData,
				});
				setPhase("error");
			} else if (parsedError.type === "EXAM_FORBIDDEN") {
				setExamError({
					type: "MAX_ATTEMPTS",
					message: parsedError.message,
					data: parsedError.data,
				});
				setPhase("error");
			} else {
				setExamState((prev) => ({
					...prev,
					error: parsedError.message,
					loading: false,
				}));
			}
		}
	}, [packageSlug, examSlug, getDeviceInfo, startTimer]);

	const reviewAnswers = useCallback(() => {
		setExamState((prev: ExamState) => ({ ...prev, isReviewing: true }));
		setPhase("review");
		// Timer continues running during review
	}, []);

	const backToExam = useCallback(() => {
		setExamState((prev: ExamState) => ({ ...prev, isReviewing: false }));
		setPhase("exam");
		// Timer continues running, no need to restart
	}, []);

	const submitExam = useCallback(async () => {
		if (!sessionId) return;

		stopTimer();

		try {
			// Submit exam via API
			const result = await examAPI.submitExam({
				session_id: sessionId,
			});

			// Transform API result to match frontend types
			const examResult: ExamResult = {
				score: result.score,
				correctAnswers: result.correct_answers,
				totalQuestions: result.total_questions,
				timeSpent: result.time_taken_seconds,
				isPassed: result.passed,
				answers: Object.values(examState.userAnswers),
			};

			setExamState((prev: ExamState) => ({
				...prev,
				showResults: true,
			}));

			setPhase("results");

			// Store result for later use
			window.examResult = examResult;
		} catch (error) {
			console.error("Failed to submit exam:", error);

			// Parse error to check if it's a known API error
			const { type } = parseExamApiError(error as Error);

			let errorMessage = "Failed to submit exam. Please try again.";
			if (type === "SESSION_NOT_FOUND") {
				errorMessage = "Exam session expired. Please restart the exam.";
			} else if (type === "EXAM_ALREADY_SUBMITTED") {
				errorMessage = "Exam has already been submitted.";
				// Don't show error for already submitted - this might be expected
				return;
			}

			setExamState((prev) => ({
				...prev,
				error: errorMessage,
			}));

			// Re-enable navigation if submission failed
			setPhase("exam");
		}
	}, [sessionId, examState.userAnswers, stopTimer]);

	const restartExam = useCallback(() => {
		// For restart, we need to start a new session
		setSessionId(null);
		setExamState((prev: ExamState) => ({
			...prev,
			currentQuestionIndex: 0,
			currentQuestion: prev.exam?.questions[0] || null,
			userAnswers: {},
			timeRemaining: prev.exam?.duration_minutes
				? prev.exam.duration_minutes * 60
				: 0,
			isReviewing: false,
			showResults: false,
			loading: true,
		}));
		setPhase("loading");

		// Trigger re-initialization
		window.location.reload(); // Simple approach for now
	}, []);

	// Cleanup timers on unmount
	useEffect(() => {
		return () => {
			if (timerRef.current) {
				clearInterval(timerRef.current);
			}
		};
	}, []);

	// Derived state
	const answeredCount = Object.keys(examState.userAnswers).filter(
		(questionId) => {
			const answer = examState.userAnswers[parseInt(questionId)];
			return (
				answer && !answer.isSkipped && answer.selectedOptions.length > 0
			);
		}
	).length;
	const totalQuestions = examState.exam?.questions.length || 0;

	const progress =
		totalQuestions > 0 ? (answeredCount / totalQuestions) * 100 : 0;

	const currentAnswer = examState.currentQuestion
		? examState.userAnswers[examState.currentQuestion.id]
		: undefined;

	const canGoNext = examState.currentQuestionIndex < totalQuestions - 1;
	const canGoPrevious = examState.currentQuestionIndex > 0;
	const canSubmit = answeredCount > 0; // Allow submission with partial answers

	return {
		// State
		examState,
		phase,
		progress,
		answeredCount,
		totalQuestions,
		currentAnswer,
		examError, // New error state for specific exam errors
		sessionId, // Expose session ID for debugging

		// Navigation state
		canGoNext,
		canGoPrevious,
		canSubmit,

		// Actions
		startExam,
		goToQuestion,
		nextQuestion,
		previousQuestion,
		updateAnswer,
		reviewAnswers,
		backToExam,
		submitExam,
		restartExam,
	};
}
