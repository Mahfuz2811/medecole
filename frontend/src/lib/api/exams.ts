import { AxiosError } from "axios";
import { packagesApiClient } from "./client";
import {
	ExamApiError,
	ExamListResponse,
	ExamMeta,
	ExamResponse,
	ExamSessionResponse,
	GetExamResponse,
	RawAnswersData,
	StartExamRequest,
	StartExamResponse,
	SubmitExamRequest,
	SubmitExamResponse,
	SyncAnswerRequest,
	SyncAnswerResponse,
} from "./types";

export const examAPI = {
	/**
	 * Get exam metadata directly (preferred method to avoid duplicate API calls)
	 */
	async getExamMeta(examSlug: string): Promise<ExamMeta> {
		try {
			const response = await packagesApiClient.get<ExamMeta>(
				`/exams/meta/${examSlug}`
			);
			return response.data;
		} catch (error) {
			const axiosError = error as AxiosError;
			if (axiosError.response?.status === 404) {
				throw new Error(
					JSON.stringify({
						type: "EXAM_NOT_FOUND",
						message: "Exam not found",
					})
				);
			}
			throw new Error(
				JSON.stringify({
					type: "NETWORK_ERROR",
					message:
						axiosError.message || "Failed to fetch exam metadata",
				})
			);
		}
	},

	/**
	 * Get exam metadata from package exams list (fallback method)
	 */
	async getExamFromPackage(
		packageSlug: string,
		examSlug: string
	): Promise<GetExamResponse> {
		try {
			const response = await packagesApiClient.get<ExamListResponse>(
				`/exams/${packageSlug}`
			);

			// Find the specific exam from the list
			const exam = response.data.exams.find(
				(e: ExamResponse) => e.slug === examSlug
			);
			if (!exam) {
				throw new Error(
					JSON.stringify({
						type: "EXAM_NOT_FOUND",
						message: "Exam not found in package",
					})
				);
			}

			// Transform to expected format
			return {
				exam_meta: {
					id: exam.id,
					title: exam.title,
					slug: exam.slug,
					duration_minutes: exam.duration_minutes,
					total_questions: exam.total_questions,
					passing_score: exam.passing_score,
					total_marks: exam.total_marks,
					max_attempts: exam.max_attempts,
					instructions: exam.instructions || "",
				},
			};
		} catch (error) {
			const axiosError = error as AxiosError;
			if (axiosError.response?.status === 404) {
				throw new Error(
					JSON.stringify({
						type: "EXAM_NOT_FOUND",
						message: "Exam not found",
					})
				);
			}
			throw error;
		}
	},

	/**
	 * Start a new exam session
	 */
	async startExam(
		examSlug: string,
		request: StartExamRequest
	): Promise<StartExamResponse> {
		try {
			const response = await packagesApiClient.post<StartExamResponse>(
				`/exams/${examSlug}/start`,
				request
			);
			return response.data;
		} catch (error) {
			const axiosError = error as AxiosError;
			// Handle specific error codes
			if (axiosError.response?.status === 409) {
				const errorData = axiosError.response.data as ExamApiError;
				throw new Error(
					JSON.stringify({
						type: "EXAM_CONFLICT",
						data: errorData,
					})
				);
			}

			if (axiosError.response?.status === 403) {
				const errorData = axiosError.response.data as ExamApiError;
				throw new Error(
					JSON.stringify({
						type: "EXAM_FORBIDDEN",
						data: errorData,
					})
				);
			}

			throw error;
		}
	},

	/**
	 * Get current exam session
	 */
	async getSession(sessionId: string): Promise<ExamSessionResponse> {
		try {
			const response = await packagesApiClient.get<ExamSessionResponse>(
				`/exams/session/${sessionId}`
			);
			return response.data;
		} catch (error) {
			const axiosError = error as AxiosError;
			if (axiosError.response?.status === 404) {
				throw new Error(
					JSON.stringify({
						type: "SESSION_NOT_FOUND",
						message: "Exam session not found or expired",
					})
				);
			}

			if (axiosError.response?.status === 410) {
				throw new Error(
					JSON.stringify({
						type: "SESSION_EXPIRED",
						message: "Exam session has expired",
					})
				);
			}

			throw error;
		}
	},

	/**
	 * Sync user answers for a session
	 */
	async syncSession(
		sessionId: string,
		request: SyncAnswerRequest
	): Promise<SyncAnswerResponse> {
		try {
			const response = await packagesApiClient.put<SyncAnswerResponse>(
				`/exams/session/${sessionId}/sync`,
				request
			);
			return response.data;
		} catch (error) {
			const axiosError = error as AxiosError;
			if (axiosError.response?.status === 404) {
				throw new Error(
					JSON.stringify({
						type: "SESSION_NOT_FOUND",
						message: "Exam session not found",
					})
				);
			}

			if (axiosError.response?.status === 409) {
				throw new Error(
					JSON.stringify({
						type: "EXAM_ALREADY_SUBMITTED",
						message: "Cannot modify answers after exam submission",
					})
				);
			}

			throw error;
		}
	},

	/**
	 * Submit exam and finalize session
	 */
	async submitExam(request: SubmitExamRequest): Promise<SubmitExamResponse> {
		try {
			const response = await packagesApiClient.post<SubmitExamResponse>(
				"/exams/submit",
				request
			);
			return response.data;
		} catch (error) {
			const axiosError = error as AxiosError;
			if (axiosError.response?.status === 404) {
				throw new Error(
					JSON.stringify({
						type: "SESSION_NOT_FOUND",
						message: "Exam session not found",
					})
				);
			}

			if (axiosError.response?.status === 409) {
				throw new Error(
					JSON.stringify({
						type: "EXAM_ALREADY_SUBMITTED",
						message: "Exam has already been submitted",
					})
				);
			}

			throw error;
		}
	},
};

// Helper function to parse API errors
export const parseExamApiError = (
	error: Error
): { type: string; data?: ExamApiError; message: string } => {
	try {
		const parsed = JSON.parse(error.message);
		return parsed;
	} catch {
		return {
			type: "UNKNOWN_ERROR",
			message: error.message || "An unexpected error occurred",
		};
	}
};

// Session-based exam results API - returns raw answers_data for frontend processing
export const getExamResultsBySession = async (
	sessionId: string
): Promise<RawAnswersData> => {
	try {
		const response = await packagesApiClient.get<{
			success: boolean;
			message: string;
			data: RawAnswersData; // Raw answers_data from backend
		}>(`/exams/results/${sessionId}`);

		return response.data.data;
	} catch (error) {
		const axiosError = error as AxiosError;
		if (axiosError.response?.status === 404) {
			throw new Error(
				JSON.stringify({
					type: "RESULTS_NOT_FOUND",
					message: "Exam results not found",
				})
			);
		}

		if (axiosError.response?.status === 401) {
			throw new Error(
				JSON.stringify({
					type: "UNAUTHORIZED",
					message: "You are not authorized to view these results",
				})
			);
		}

		throw error;
	}
};
