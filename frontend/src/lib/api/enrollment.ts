import axios from "axios";
import type {
	CouponValidationRequest,
	CouponValidationResponse,
	DashboardEnrollmentsResponse,
	EnrollmentRequest,
	EnrollmentResponse,
	EnrollmentStatusResponse,
	ErrorResponse,
} from "./enrollment-types";
import { auth } from "./utils";

// Create a dedicated authenticated API client for enrollment operations
const enrollmentApiClient = axios.create({
	baseURL:
		process.env.NEXT_PUBLIC_PACKAGES_API_URL || "http://localhost:8080/api",
	timeout: 10000,
	headers: {
		"Content-Type": "application/json",
	},
});

// Add authorization interceptor ONLY for enrollment API
enrollmentApiClient.interceptors.request.use(
	(config) => {
		const token = localStorage.getItem("authToken");
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

// Handle auth errors ONLY for enrollment API
enrollmentApiClient.interceptors.response.use(
	(response) => response,
	(error) => {
		// Only force logout for critical enrollment operations, not status checks
		if (error.response?.status === 401) {
			const isStatusCheck = error.config?.url?.includes(
				"/enrollments/status"
			);
			const isCouponValidation =
				error.config?.url?.includes("/validate-coupon");

			// Don't force logout for optional operations like status checks and coupon validation
			if (!isStatusCheck && !isCouponValidation) {
				// Clear all auth data and cache on unauthorized error
				auth.clearAuthData();
				window.location.href = "/auth";
			}
		}
		return Promise.reject(error);
	}
);

export class EnrollmentAPI {
	/**
	 * Enroll user in a package
	 */
	static async enrollInPackage(
		request: EnrollmentRequest
	): Promise<EnrollmentResponse> {
		try {
			const response = await enrollmentApiClient.post<EnrollmentResponse>(
				"/enrollments",
				request
			);
			return response.data;
		} catch (error: unknown) {
			throw this.handleError(error);
		}
	}

	/**
	 * Get user's enrollments optimized for dashboard display
	 */
	static async getDashboardEnrollments(): Promise<DashboardEnrollmentsResponse> {
		try {
			const response =
				await enrollmentApiClient.get<DashboardEnrollmentsResponse>(
					"/dashboard/enrollments"
				);
			return response.data;
		} catch (error: unknown) {
			throw this.handleError(error);
		}
	}

	/**
	 * Check enrollment status for a package
	 */
	static async checkEnrollmentStatus(
		packageId: number
	): Promise<EnrollmentStatusResponse> {
		try {
			const response =
				await enrollmentApiClient.get<EnrollmentStatusResponse>(
					`/enrollments/status`,
					{
						params: { package_id: packageId },
					}
				);
			return response.data;
		} catch (error: unknown) {
			throw this.handleError(error);
		}
	}

	/**
	 * Validate coupon for a package
	 */
	static async validateCoupon(
		request: CouponValidationRequest
	): Promise<CouponValidationResponse> {
		try {
			const response =
				await enrollmentApiClient.post<CouponValidationResponse>(
					"/enrollments/validate-coupon",
					request
				);
			return response.data;
		} catch (error: unknown) {
			throw this.handleError(error);
		}
	}

	/**
	 * Handle API errors
	 */
	private static handleError(error: unknown): Error {
		if (
			typeof error === "object" &&
			error !== null &&
			"response" in error
		) {
			const axiosError = error as { response?: { data?: ErrorResponse } };
			if (axiosError.response?.data) {
				const errorData = axiosError.response.data;
				return new Error(
					errorData.message || errorData.error || "An error occurred"
				);
			}
		}

		if (error instanceof Error) {
			return new Error(error.message);
		}

		return new Error("Network error or server unavailable");
	}
}

export default EnrollmentAPI;
