import { useAuth } from "@/contexts/AuthContext";
import { EnrollmentAPI } from "@/lib/api/enrollment";
import type { DashboardEnrollmentsResponse } from "@/lib/api/enrollment-types";
import { useEffect, useState } from "react";

export function useProfileEnrollments() {
	const { user } = useAuth();
	const [enrollments, setEnrollments] =
		useState<DashboardEnrollmentsResponse | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const fetchEnrollments = async () => {
			try {
				setLoading(true);
				setError(null);
				const data = await EnrollmentAPI.getDashboardEnrollments();
				setEnrollments(data);
			} catch (error) {
				console.error("Failed to fetch enrollments:", error);
				setError("Failed to load subscription data");
			} finally {
				setLoading(false);
			}
		};

		// Only fetch if user is authenticated
		if (user) {
			fetchEnrollments();
		}
	}, [user]);

	return {
		enrollments: enrollments?.enrollments || [],
		totalEnrollments: enrollments?.total || 0,
		activeEnrollments: enrollments?.active || 0,
		loading,
		error,
	};
}
