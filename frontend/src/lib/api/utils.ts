import { AuthResponse, User } from "./types";

export const auth = {
	// Store auth data in localStorage
	setAuthData: (authResponse: AuthResponse) => {
		localStorage.setItem("authToken", authResponse.token);
		localStorage.setItem("user", JSON.stringify(authResponse.user));
	},

	// Get stored auth data
	getAuthData: (): { token: string | null; user: User | null } => {
		const token = localStorage.getItem("authToken");
		const userStr = localStorage.getItem("user");
		const user = userStr ? JSON.parse(userStr) : null;
		return { token, user };
	},

	// Check if user is authenticated
	isAuthenticated: (): boolean => {
		const token = localStorage.getItem("authToken");
		return !!token;
	},

	// Clear all auth data and cache
	clearAuthData: () => {
		// Clear specific auth data
		const authKeys = ["authToken", "user", "profile_cache_timestamp"];
		authKeys.forEach((key) => localStorage.removeItem(key));

		// Clear all cache data with known patterns
		const enrollmentCacheKeys: string[] = [];
		const packageCacheKeys: string[] = [];
		const packageDetailCacheKeys: string[] = [];
		const otherCacheKeys: string[] = [];

		// Iterate through all localStorage keys to find cache entries
		for (let i = 0; i < localStorage.length; i++) {
			const key = localStorage.key(i);
			if (key) {
				// Remove enrollment status cache (pattern: enrollment_status_${packageId})
				if (key.startsWith("enrollment_status_")) {
					enrollmentCacheKeys.push(key);
				}
				// Remove package list cache (pattern: packages_cache_${hash})
				else if (key.startsWith("packages_cache_")) {
					packageCacheKeys.push(key);
				}
				// Remove individual package details cache (pattern: package_${slug})
				else if (key.startsWith("package_")) {
					packageDetailCacheKeys.push(key);
				}
				// Remove any other cache entries (patterns containing _cache)
				else if (key.includes("_cache")) {
					otherCacheKeys.push(key);
				}
			}
		}

		// Remove all identified cache keys
		const allCacheKeys = [
			...enrollmentCacheKeys,
			...packageCacheKeys,
			...packageDetailCacheKeys,
			...otherCacheKeys,
		];
		allCacheKeys.forEach((key) => {
			localStorage.removeItem(key);
		});

		// For development: Log what was cleared with detailed breakdown
		if (process.env.NODE_ENV === "development") {
			console.log("🧹 Cleared localStorage:", {
				authData: authKeys,
				enrollmentCache:
					enrollmentCacheKeys.length > 0
						? enrollmentCacheKeys
						: "none",
				packageListCache:
					packageCacheKeys.length > 0 ? packageCacheKeys : "none",
				packageDetailCache:
					packageDetailCacheKeys.length > 0
						? packageDetailCacheKeys
						: "none",
				otherCache: otherCacheKeys.length > 0 ? otherCacheKeys : "none",
				totalCleared: authKeys.length + allCacheKeys.length,
			});
		}
	},

	// Clear ALL localStorage data (more aggressive approach if needed)
	clearAllAppData: () => {
		// Get all keys first (since localStorage.length changes as we remove items)
		const allKeys: string[] = [];
		for (let i = 0; i < localStorage.length; i++) {
			const key = localStorage.key(i);
			if (key) {
				allKeys.push(key);
			}
		}

		// Remove all keys
		allKeys.forEach((key) => {
			localStorage.removeItem(key);
		});

		// For development: Log what was cleared
		if (process.env.NODE_ENV === "development") {
			console.log("🧹 Cleared ALL localStorage data:", {
				totalKeys: allKeys.length,
				keys: allKeys,
			});
		}
	},

	// Format MSISDN for API (local Bangladesh format)
	formatMSISDN: (msisdn: string): string => {
		// Remove all non-digits
		const digits = msisdn.replace(/\D/g, "");

		// Ensure it's 11 digits starting with 0
		if (digits.length === 0) return "";

		// If user types without leading 0, add it
		if (digits.length > 0 && !digits.startsWith("0")) {
			return "0" + digits.slice(0, 10);
		}

		// Limit to 11 digits max
		return digits.slice(0, 11);
	},
};
