// Enrollment API Types
export interface EnrollmentRequest {
	package_id: number;
	coupon_code?: string;
}

export interface EnrollmentResponse {
	id: number;
	user_id: number;
	package_id: number;
	enrollment_type: string;
	enrolled_at: string;
	expires_at?: string;
	enrolled_package_type: string;
	enrolled_price: number;
	payment_status: string;
	payment_amount?: number;
	coupon_code?: string;
	original_price?: number;
	discount_percentage?: number;
	discount_amount?: number;
	final_price?: number;
	effective_status: string;
	can_access_content: boolean;
	is_active: boolean;
	package?: {
		id: number;
		name: string;
		slug: string;
		package_type: string;
		price: number;
	};
	price_calculation?: PriceCalculationResult;
	created_at: string;
	updated_at: string;
}

export interface UserEnrollmentsResponse {
	enrollments: EnrollmentResponse[];
	total: number;
	active: number;
	expired: number;
}

// Dashboard optimized enrollment types
export interface DashboardEnrollmentDTO {
	id: number;
	package_id: number;
	package_name: string;
	package_slug: string; // Added slug for navigation
	package_type: string;
	expiry_date?: string; // ISO string or null
	status: string; // "active" or "enrolled"
	progress: number; // Percentage 0-100
	total_exams: number;
	completed_exams: number;
}

export interface DashboardEnrollmentsResponse {
	enrollments: DashboardEnrollmentDTO[];
	total: number;
	active: number;
}

export interface CouponValidationRequest {
	coupon_code: string;
	package_id: number;
}

export interface CouponValidationResponse {
	valid: boolean;
	coupon_code: string;
	discount_percentage: number;
	message: string;
	price_calculation?: PriceCalculationResult;
}

export interface PriceCalculationResult {
	original_price: number;
	discount_percentage: number;
	discount_amount: number;
	final_price: number;
	coupon_code?: string;
}

export interface EnrollmentStatusResponse {
	has_active_enrollment: boolean;
	enrollment?: EnrollmentResponse;
}

// Error Response
export interface ErrorResponse {
	error: string;
	message: string;
}
