import { EnrollmentAPI } from "@/lib/api/enrollment";
import type {
	CouponValidationRequest,
	CouponValidationResponse,
	EnrollmentRequest,
	EnrollmentResponse,
	EnrollmentStatusResponse,
} from "@/lib/api/enrollment-types";
import { useCallback, useState } from "react";

// Simple cache for enrollment status to prevent duplicate requests
// Cache expires after 5 minutes to ensure reasonably fresh data
const ENROLLMENT_STATUS_CACHE_TTL = 5 * 60 * 1000; // 5 minutes in milliseconds

interface CachedEnrollmentStatus {
	data: EnrollmentStatusResponse;
	timestamp: number;
}

// Helper functions for persistent cache across page reloads
// Helper function to get cached enrollment status
const getCachedEnrollmentStatus = (
	packageId: number
): CachedEnrollmentStatus | null => {
	const cacheKey = `enrollment_status_${packageId}`;
	const cached = localStorage.getItem(cacheKey);
	if (!cached) return null;

	try {
		const data: CachedEnrollmentStatus = JSON.parse(cached);
		const now = Date.now();
		const isExpired = now - data.timestamp > ENROLLMENT_STATUS_CACHE_TTL;

		if (process.env.NODE_ENV === "development") {
			console.log(`Enrollment Cache Check [${packageId}]:`, {
				cached: data.data.has_active_enrollment,
				cacheAge: Math.round((now - data.timestamp) / 1000),
				expired: isExpired,
			});
		}

		if (isExpired) {
			localStorage.removeItem(cacheKey);
			return null;
		}

		return data;
	} catch {
		localStorage.removeItem(cacheKey);
		return null;
	}
};

const setCachedEnrollmentStatus = (
	packageId: number,
	data: CachedEnrollmentStatus
): void => {
	if (typeof window !== "undefined") {
		try {
			localStorage.setItem(
				`enrollment_status_${packageId}`,
				JSON.stringify(data)
			);
		} catch {
			// Ignore localStorage errors (quota exceeded, etc.)
		}
	}
};

const clearCachedEnrollmentStatus = (packageId?: number): void => {
	if (typeof window === "undefined") return;

	if (packageId) {
		localStorage.removeItem(`enrollment_status_${packageId}`);
	} else {
		// Clear all enrollment status cache
		const keys = Object.keys(localStorage);
		keys.forEach((key) => {
			if (key.startsWith("enrollment_status_")) {
				localStorage.removeItem(key);
			}
		});
	}
};

// Global cache to persist across component re-mounts (but not page reloads)
const enrollmentStatusCache = new Map<number, CachedEnrollmentStatus>();

// Global flag to prevent duplicate concurrent requests
const ongoingRequests = new Map<number, Promise<EnrollmentStatusResponse>>();

interface UseEnrollmentState {
	loading: boolean;
	error: string | null;
	enrollmentStatus: EnrollmentStatusResponse | null;
}

interface UseEnrollmentReturn extends UseEnrollmentState {
	enrollInPackage: (
		request: EnrollmentRequest
	) => Promise<EnrollmentResponse | null>;
	checkEnrollmentStatus: (packageId: number) => Promise<void>;
	validateCoupon: (
		request: CouponValidationRequest
	) => Promise<CouponValidationResponse | null>;
	clearError: () => void;
	clearEnrollmentStatus: () => void;
	clearEnrollmentStatusCache: (packageId?: number) => void;
	initializeWithCache: (packageId: number) => boolean;
	getCacheStats?: () => {
		totalEntries: number;
		validEntries: number;
		expiredEntries: number;
		ongoingRequests: number;
	} | null; // Only available in development
}

export function useEnrollment(packageId?: number): UseEnrollmentReturn {
	// Initialize state with cached data if packageId is provided
	const getInitialState = (): UseEnrollmentState => {
		const baseState = {
			loading: false,
			error: null,
			enrollmentStatus: null,
		};

		// If packageId is provided, try to initialize with cached data
		if (packageId && typeof window !== "undefined") {
			const cachedEntry = getCachedEnrollmentStatus(packageId);
			if (cachedEntry) {
				const isExpired =
					Date.now() - cachedEntry.timestamp >
					ENROLLMENT_STATUS_CACHE_TTL;
				if (!isExpired) {
					// Initialize with cached data to prevent blink
					enrollmentStatusCache.set(packageId, cachedEntry);
					return {
						...baseState,
						enrollmentStatus: cachedEntry.data,
					};
				}
			}
		}

		return baseState;
	};

	const [state, setState] = useState<UseEnrollmentState>(getInitialState());

	// Initialize with cached data if available (prevents UI blinks)
	const initializeWithCache = useCallback((packageId: number) => {
		// Check memory cache first
		let cachedEntry = enrollmentStatusCache.get(packageId);

		// If not in memory, check localStorage
		if (!cachedEntry) {
			const persistedEntry = getCachedEnrollmentStatus(packageId);
			if (persistedEntry) {
				cachedEntry = persistedEntry;
				enrollmentStatusCache.set(packageId, persistedEntry);
			}
		}

		// If we have valid cached data, initialize with it immediately
		if (cachedEntry) {
			const isExpired =
				Date.now() - cachedEntry.timestamp >
				ENROLLMENT_STATUS_CACHE_TTL;
			if (!isExpired) {
				setState((prev) => ({
					...prev,
					enrollmentStatus: cachedEntry.data,
				}));
				return true; // Indicates we initialized with cache
			} else {
				// Clean up expired cache
				enrollmentStatusCache.delete(packageId);
				clearCachedEnrollmentStatus(packageId);
			}
		}
		return false; // No cache available
	}, []);

	const clearError = useCallback(() => {
		setState((prev) => ({ ...prev, error: null }));
	}, []);

	const clearEnrollmentStatus = useCallback(() => {
		setState((prev) => ({ ...prev, enrollmentStatus: null }));
	}, []);

	const clearEnrollmentStatusCache = useCallback((packageId?: number) => {
		if (packageId) {
			// Clear cache for specific package from both memory and localStorage
			enrollmentStatusCache.delete(packageId);
			clearCachedEnrollmentStatus(packageId);
		} else {
			// Clear entire cache from both memory and localStorage
			enrollmentStatusCache.clear();
			clearCachedEnrollmentStatus();
		}
	}, []);

	// Development helper: Get cache statistics
	const getCacheStats = useCallback(() => {
		if (process.env.NODE_ENV !== "development") return null;

		const now = Date.now();
		const cacheEntries = Array.from(enrollmentStatusCache.entries());

		return {
			totalEntries: cacheEntries.length,
			validEntries: cacheEntries.filter(
				([, entry]) =>
					now - entry.timestamp <= ENROLLMENT_STATUS_CACHE_TTL
			).length,
			expiredEntries: cacheEntries.filter(
				([, entry]) =>
					now - entry.timestamp > ENROLLMENT_STATUS_CACHE_TTL
			).length,
			ongoingRequests: ongoingRequests.size,
		};
	}, []);

	const enrollInPackage = useCallback(
		async (
			request: EnrollmentRequest
		): Promise<EnrollmentResponse | null> => {
			setState((prev) => ({ ...prev, loading: true, error: null }));

			try {
				const response = await EnrollmentAPI.enrollInPackage(request);

				// Clear enrollment status cache for this package after successful enrollment
				enrollmentStatusCache.delete(request.package_id);
				clearCachedEnrollmentStatus(request.package_id);

				setState((prev) => ({ ...prev, loading: false }));
				return response;
			} catch (error) {
				const errorMessage =
					error instanceof Error
						? error.message
						: "Failed to enroll in package";
				setState((prev) => ({
					...prev,
					loading: false,
					error: errorMessage,
				}));
				return null;
			}
		},
		[]
	);

	const checkEnrollmentStatus = useCallback(
		async (packageId: number): Promise<void> => {
			// Try to initialize with cached data first (prevents blink)
			const hasCache = initializeWithCache(packageId);

			// If we have cached data, we can skip the API call entirely
			if (hasCache) {
				if (process.env.NODE_ENV === "development") {
					console.log(
						"✅ Using cached enrollment status (no API call)"
					);
				}
				return;
			}

			// No cache available, proceed with API call
			setState((prev) => ({ ...prev, loading: true, error: null }));

			// Check if there's already an ongoing request for this package
			const ongoingRequest = ongoingRequests.get(packageId);
			if (ongoingRequest) {
				try {
					const status = await ongoingRequest;
					setState((prev) => ({
						...prev,
						enrollmentStatus: status,
						loading: false,
					}));
					return;
				} catch {
					// If ongoing request fails, we'll make a new one below
				}
			}

			if (process.env.NODE_ENV === "development") {
				console.log("🔄 Fetching fresh enrollment status from API");
			}

			// Create and store the API request promise to prevent duplicates
			const requestPromise =
				EnrollmentAPI.checkEnrollmentStatus(packageId);
			ongoingRequests.set(packageId, requestPromise);

			try {
				const status = await requestPromise;

				// Cache the response with timestamp in both memory and localStorage
				const cacheEntry = {
					data: status,
					timestamp: Date.now(),
				};

				enrollmentStatusCache.set(packageId, cacheEntry);
				setCachedEnrollmentStatus(packageId, cacheEntry);

				setState((prev) => ({
					...prev,
					loading: false,
					enrollmentStatus: status,
				}));
			} catch (error) {
				const errorMessage =
					error instanceof Error
						? error.message
						: "Failed to check enrollment status";
				setState((prev) => ({
					...prev,
					loading: false,
					error: errorMessage,
				}));
			} finally {
				// Remove the request from ongoing requests
				ongoingRequests.delete(packageId);
			}
		},
		[initializeWithCache]
	);

	const validateCoupon = useCallback(
		async (
			request: CouponValidationRequest
		): Promise<CouponValidationResponse | null> => {
			setState((prev) => ({ ...prev, loading: true, error: null }));

			try {
				const response = await EnrollmentAPI.validateCoupon(request);
				setState((prev) => ({ ...prev, loading: false }));
				return response;
			} catch (error) {
				const errorMessage =
					error instanceof Error
						? error.message
						: "Failed to validate coupon";
				setState((prev) => ({
					...prev,
					loading: false,
					error: errorMessage,
				}));
				return null;
			}
		},
		[]
	);

	return {
		...state,
		enrollInPackage,
		checkEnrollmentStatus,
		validateCoupon,
		clearError,
		clearEnrollmentStatus,
		clearEnrollmentStatusCache,
		initializeWithCache,
		...(process.env.NODE_ENV === "development" && { getCacheStats }),
	};
}
