// API Response Types
export interface ApiError {
	message: string;
	details?: string;
}

// Package Types
export interface PackageImageResponse {
	original: string;
	mobile: string;
	tablet: string;
	desktop: string;
	thumbnail: string;
	alt_text: string;
	metadata?: {
		width: number;
		height: number;
		file_size: number;
		format: string;
		original_url?: string;
	};
}

export interface PackageResponse {
	id: number;
	name: string;
	slug: string;
	description?: string;
	package_type: "FREE" | "PREMIUM";
	price: number;
	images: PackageImageResponse;
	coupon_code?: string;
	validity_type: "FIXED" | "RELATIVE";
	validity_days?: number;
	validity_date?: string;
	total_exams: number;
	is_active: boolean;
	sort_order: number;
	created_at: string;
	updated_at: string;
	// Analytics fields (optional for backward compatibility)
	enrollment_count?: number;
	active_enrollment_count?: number;
	last_enrollment_at?: string;
	// Exam data - now required for package details
	exams: PackageExamResponse[];
}

export interface PackageListResponse {
	packages: PackageResponse[];
	total: number;
	page: number;
	per_page: number;
	total_pages: number;
}

export interface PackageListRequest {
	page?: number;
	per_page?: number;
	type?: "FREE" | "PREMIUM";
	sort?: "name" | "price" | "created_at" | "sort_order";
	order?: "asc" | "desc";
	search?: string;
}

// Exam Types
export type ExamType = "DAILY" | "MOCK" | "REVIEW" | "FINAL";
export type ExamStatus = "DRAFT" | "SCHEDULED" | "ACTIVE" | "COMPLETED";
export type ComputedStatus = "UPCOMING" | "LIVE" | "COMPLETED" | "AVAILABLE";
export type AttemptStatus =
	| "STARTED"
	| "COMPLETED"
	| "AUTO_SUBMITTED"
	| "ABANDONED";

export interface ExamResponse {
	id: number;
	title: string;
	slug: string;
	description: string;
	exam_type: ExamType;
	total_questions: number;
	duration_minutes: number;
	passing_score: number;
	total_marks: number;
	max_attempts: number;
	instructions?: string;
	scheduled_start_date?: string;
	scheduled_end_date?: string;
	attempt_count: number;
	average_score: number;
	pass_rate: number;
	computed_status: ComputedStatus;
	sort_order: number;
	has_attempted: boolean;
	user_attempt?: UserAttempt;
}

export interface UserAttempt {
	id: number;
	status: AttemptStatus;
	started_at: string;
	completed_at?: string;
	score?: number;
	correct_answers?: number;
	is_passed?: boolean;
	time_spent?: number;
	session_id?: string;
}

export interface ExamListResponse {
	package: PackageInfoResponse;
	exams: ExamResponse[];
}

export interface PackageInfoResponse {
	id: number;
	name: string;
	slug: string;
	description?: string;
	package_type: "FREE" | "PREMIUM";
	price: number;
	validity_type: "FIXED" | "RELATIVE";
	validity_days?: number;
	validity_date?: string;
	total_exams: number;
	enrollment_count: number;
	active_enrollment_count: number;
}

export interface PackageExamResponse {
	exam: ExamResponse;
}

// Auth Types
export interface User {
	id: number;
	name: string;
	msisdn?: string; // Optional for social auth users
	email?: string; // For social auth users
	auth_provider: string; // "local", "google", "facebook"
	profile_picture?: string; // Avatar URL from social provider
	email_verified: boolean; // Email verification status
	isActive: boolean;
	createdAt: string;
	updatedAt?: string;
}

export interface AuthResponse {
	user: User;
	token: string;
}

export interface LoginRequest {
	msisdn: string;
	password: string;
}

export interface RegisterRequest {
	name: string;
	msisdn: string;
	password: string;
}

// Exam Session Types
export interface DeviceInfo {
	browser: string;
	user_agent: string;
	ip_address: string;
}

export interface StartExamRequest {
	package_slug: string;
	device_info: DeviceInfo;
}

export interface ExamMeta {
	id: number;
	title: string;
	slug: string;
	duration_minutes: number;
	total_questions: number;
	passing_score: number;
	total_marks: number;
	max_attempts: number;
	instructions: string | null;
}

export interface StartExamResponse {
	session_id: string;
	attempt_id: number;
	exam_meta: ExamMeta;
}

export interface GetExamResponse {
	exam_meta: ExamMeta;
}

export interface ExamQuestion {
	id: number;
	question_text: string;
	question_type: "SINGLE_CHOICE" | "MULTIPLE_CHOICE" | "TRUE_FALSE";
	options: ExamOption[];
	points: number;
	time_limit_seconds?: number;
	explanation?: string;
}

export interface ExamOption {
	id: number;
	option_text: string;
	is_correct: boolean;
	explanation?: string;
}

export interface ExamSessionResponse {
	exam: {
		id: number;
		title: string;
		slug: string;
		description?: string;
		exam_type: string;
		duration_minutes: number;
		passing_score: number;
		max_attempts: number;
		instructions?: string;
		questions: SecureExamQuestion[];
	};
	session: {
		session_id: string;
		attempt_id: number;
		status: string;
		time_remaining: number;
		time_limit_seconds: number;
		started_at: string;
		can_submit: boolean;
		can_pause: boolean;
		last_activity: string;
		saved_answers?: Array<{
			question_id: number;
			selected_option: string; // Single option for SBA, JSON array string for TRUE_FALSE
		}>; // Optional field for saved answers from Redis
	};
}

export interface SecureExamQuestion {
	id: number;
	question_text: string;
	question_type: string;
	options: Record<string, { text: string }>;
	points: number;
}

export interface UserAnswer {
	question_id: number;
	selected_option_ids: number[];
	is_flagged: boolean;
	time_spent_seconds: number;
	answered_at?: string;
}

export interface SyncAnswerRequest {
	answers: Array<{
		question_id: number;
		selected_option: string; // Single option for SBA questions, JSON array string for TRUE_FALSE questions
	}>;
}

export interface SyncAnswerResponse {
	success: boolean;
	synced_count: number;
	last_sync_at: string;
	time_remaining: number;
}

export interface SubmitExamRequest {
	session_id: string;
}

export interface SubmitExamResponse {
	session_id: string;
	score: number;
	passed: boolean;
	total_questions: number;
	correct_answers: number;
	time_taken_seconds: number;
	submitted_at: string;
}

// Error Types for specific API scenarios
export interface ExamAlreadySubmittedError extends ApiError {
	error_code: "EXAM_ALREADY_SUBMITTED";
	exam_title: string;
	submission_date: string;
	score?: number;
	passed?: boolean;
}

export interface MaxAttemptsExceededError extends ApiError {
	error_code: "MAX_ATTEMPTS_EXCEEDED";
	max_attempts: number;
	attempted_count: number;
}

export interface ExamNotAvailableError extends ApiError {
	error_code: "EXAM_NOT_AVAILABLE";
	available_from?: string;
	available_until?: string;
}

export type ExamApiError =
	| ExamAlreadySubmittedError
	| MaxAttemptsExceededError
	| ExamNotAvailableError
	| ApiError;

// Exam Results API Types
export interface ExamResultOption {
	id: number;
	option_text: string;
	is_correct: boolean;
}

export interface ExamResultQuestion {
	id: number;
	question: string;
	question_type: "SBA" | "TRUE_FALSE";
	options: ExamResultOption[]; // Enhanced to include option details
	correct_answer: number | boolean[]; // number for SBA (option index), boolean[] for TRUE_FALSE
	user_answer: number | boolean[] | (boolean | undefined)[] | null; // user's answer, can be null if not answered, undefined for unanswered TRUE_FALSE options
	is_correct: boolean;
	points: number; // Earned points (not max points)
	max_points?: number; // Optional max points for the question
	time_spent?: number; // time spent on this question in seconds
	explanation?: string;
}

// Raw question format from backend response
export interface RawQuestionOption {
	key: string;
	text: string;
	is_correct: boolean;
}

export interface RawQuestionData {
	question_id: number;
	question_text: string;
	question_type: "SBA" | "TRUE_FALSE";
	options: RawQuestionOption[];
	correct_answer: Record<string, boolean> | string; // Object for TRUE_FALSE, string for SBA
	user_answer: string[]; // e.g., ["a:true", "b:true", "c:true", "d:false", "e:false"] or ["c"]
	is_correct: boolean;
	points_earned: number;
	max_points: number;
	explanation?: string;
}

// Raw answers data type returned directly from backend (session-based API)
export interface RawAnswersData {
	answers: RawQuestionData[];
	exam_snapshot: {
		duration_minutes: number;
		passing_score: number;
		total_questions: number;
		actual_time_spent: number; // in seconds
		score?: number | null; // can be null if not scored yet
		correct_answers?: number | null; // can be null if not scored yet
		is_passed?: boolean | null; // can be null if not scored yet
	};
	submission_timestamp: string;
}
