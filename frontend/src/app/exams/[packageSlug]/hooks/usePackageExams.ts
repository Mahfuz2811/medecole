import { packagesAPI } from "@/lib/api/packages";
import { ExamResponse, PackageResponse } from "@/lib/api/types";
import { useEffect, useMemo, useState } from "react";

// Legacy interfaces for user attempt data (will be replaced when user attempts API is ready)
interface ExamAttempt {
	id: number;
	examId: number;
	status: "STARTED" | "COMPLETED" | "AUTO_SUBMITTED" | "ABANDONED";
	started_at: string;
	completed_at?: string;
	score?: number;
	correct_answers?: number;
	is_passed?: boolean;
	actualTimeSpent?: number;
	total_questions: number;
}

// Extended exam type that includes user attempt data
export interface ExtendedExam extends ExamResponse {
	userAttempt?: ExamAttempt;
}

// Transform PackageResponse to include computed fields for UI
interface PackageData extends PackageResponse {
	completedExams: number;
	progress: number;
}

interface UsePackageExamsResult {
	packageData: PackageData | null;
	allExams: ExtendedExam[];
	filteredExams: ExtendedExam[];
	paginatedExams: ExtendedExam[];
	loading: boolean;
	error: string | null;
	totalPages: number;
	searchQuery: string;
	setSearchQuery: (query: string) => void;
	selectedType: string;
	setSelectedType: (type: string) => void;
	selectedStatus: string;
	setSelectedStatus: (status: string) => void;
	currentPage: number;
	setCurrentPage: (page: number) => void;
	itemsPerPage: number;
}

export function usePackageExams(packageSlug: string): UsePackageExamsResult {
	// State management
	const [packageData, setPackageData] = useState<PackageData | null>(null);
	const [allExams, setAllExams] = useState<ExtendedExam[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	// Filter and pagination state
	const [searchQuery, setSearchQuery] = useState("");
	const [selectedType, setSelectedType] = useState<string>("ALL");
	const [selectedStatus, setSelectedStatus] = useState<string>("ALL");
	const [currentPage, setCurrentPage] = useState(1);
	const itemsPerPage = 12;

	// Fetch package and exams data from API
	useEffect(() => {
		const fetchPackageData = async () => {
			try {
				setLoading(true);
				setError(null);

				// Fetch combined package and exams data from the enhanced API
				const examData = await packagesAPI.getPackageExams(packageSlug);

				// Transform API response to extended exam format
				const extendedExams: ExtendedExam[] = examData.exams.map(
					(exam) => ({
						...exam,
						// Use real user attempt data from the API response
						userAttempt: exam.user_attempt
							? {
									id: exam.user_attempt.id,
									examId: exam.id,
									status: exam.user_attempt.status,
									started_at: exam.user_attempt.started_at,
									completed_at:
										exam.user_attempt.completed_at,
									score: exam.user_attempt.score,
									correct_answers:
										exam.user_attempt.correct_answers,
									is_passed: exam.user_attempt.is_passed,
									actualTimeSpent:
										exam.user_attempt.time_spent,
									total_questions: exam.total_questions,
							  }
							: undefined,
					})
				);

				// Calculate completion stats
				const completedExams = extendedExams.filter(
					(exam) =>
						exam.userAttempt?.status === "COMPLETED" ||
						exam.userAttempt?.status === "AUTO_SUBMITTED" ||
						exam.userAttempt?.status === "ABANDONED"
				).length;

				const progress =
					examData.package.total_exams > 0
						? Math.round(
								(completedExams /
									examData.package.total_exams) *
									100
						  )
						: 0;

				// Create package data with computed fields from the combined response
				const enhancedPackageData: PackageData = {
					...examData.package,
					// Convert from PackageInfoResponse to PackageResponse format
					id: examData.package.id,
					name: examData.package.name,
					slug: examData.package.slug,
					description: examData.package.description,
					package_type: examData.package.package_type,
					price: examData.package.price,
					validity_type: examData.package.validity_type,
					validity_days: examData.package.validity_days,
					validity_date: examData.package.validity_date,
					total_exams: examData.package.total_exams,
					enrollment_count: examData.package.enrollment_count,
					active_enrollment_count:
						examData.package.active_enrollment_count,
					// Add default values for fields not in PackageInfoResponse
					images: {
						original: "",
						mobile: "",
						tablet: "",
						desktop: "",
						thumbnail: "",
						alt_text: "",
					},
					is_active: true,
					sort_order: 0,
					created_at: "",
					updated_at: "",
					exams: [], // Not needed for this use case
					// Computed fields
					completedExams,
					progress,
				};

				setPackageData(enhancedPackageData);
				setAllExams(extendedExams);
			} catch (err) {
				console.error("Error fetching package data:", err);
				setError(
					err instanceof Error
						? err.message
						: "Failed to load package data"
				);
			} finally {
				setLoading(false);
			}
		};

		if (packageSlug) {
			fetchPackageData();
		}
	}, [packageSlug]);

	// Filter and search logic
	const filteredExams = useMemo(() => {
		return allExams.filter((exam) => {
			// Search filter
			const matchesSearch =
				searchQuery === "" ||
				exam.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
				exam.description
					.toLowerCase()
					.includes(searchQuery.toLowerCase());

			// Type filter
			const matchesType =
				selectedType === "ALL" || exam.exam_type === selectedType;

			// Status filter
			let matchesStatus = true;
			if (selectedStatus !== "ALL") {
				if (selectedStatus === "COMPLETED") {
					matchesStatus = Boolean(
						exam.userAttempt?.status === "COMPLETED"
					);
				} else if (selectedStatus === "AVAILABLE") {
					matchesStatus =
						!exam.userAttempt ||
						exam.userAttempt.status !== "COMPLETED";
				}
			}

			return matchesSearch && matchesType && matchesStatus;
		});
	}, [allExams, searchQuery, selectedType, selectedStatus]);

	// Pagination logic
	const totalPages = Math.ceil(filteredExams.length / itemsPerPage);
	const paginatedExams = useMemo(() => {
		const startIndex = (currentPage - 1) * itemsPerPage;
		return filteredExams.slice(startIndex, startIndex + itemsPerPage);
	}, [filteredExams, currentPage, itemsPerPage]);

	// Reset to first page when filters change
	useEffect(() => {
		setCurrentPage(1);
	}, [searchQuery, selectedType, selectedStatus]);

	return {
		packageData,
		allExams,
		filteredExams,
		paginatedExams,
		loading,
		error,
		totalPages,
		searchQuery,
		setSearchQuery,
		selectedType,
		setSelectedType,
		selectedStatus,
		setSelectedStatus,
		currentPage,
		setCurrentPage,
		itemsPerPage,
	};
}
