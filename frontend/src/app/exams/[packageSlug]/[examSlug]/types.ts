// Exam Types for the exam interface
export type QuestionType = "SBA" | "TRUE_FALSE";

export interface QuestionOption {
	text: string;
	is_correct?: boolean;
}

export interface Question {
	id: number;
	question_text: string;
	question_type: QuestionType;
	options: Record<string, QuestionOption>;
	explanation?: string;
	reference?: string;
	points: number;
}

export interface ExamData {
	id: number;
	title: string;
	slug: string;
	description?: string;
	total_questions: number;
	duration_minutes: number;
	passing_score: number;
	total_marks?: number;
	questions: Question[];
	instructions?: string | null;
}

export interface UserAnswer {
	questionId: number;
	selectedOptions: string[];
	timeSpent: number;
	answeredAt: Date;
	isSkipped: boolean;
}

export interface ExamSession {
	examId: number;
	userId: number;
	startedAt: Date;
	timeLimit: number; // in seconds
	answers: Record<number, UserAnswer>;
	currentQuestionIndex: number;
	isSubmitted: boolean;
}

export interface ExamState {
	exam: ExamData | null;
	session: ExamSession | null;
	loading: boolean;
	error: string | null;
	timeRemaining: number;
	currentQuestion: Question | null;
	currentQuestionIndex: number;
	userAnswers: Record<number, UserAnswer>;
	isReviewing: boolean;
	showResults: boolean;
}

export interface ExamResult {
	score: number;
	correctAnswers: number;
	totalQuestions: number;
	timeSpent: number;
	isPassed: boolean;
	answers: UserAnswer[];
}

// UI State Types
export type ExamPhase =
	| "loading"
	| "instructions"
	| "exam"
	| "review"
	| "results"
	| "error";

export interface NavigationState {
	canGoNext: boolean;
	canGoPrevious: boolean;
	canSubmit: boolean;
	hasAnswered: boolean;
}
