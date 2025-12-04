"use client";

import { User, auth, authAPI } from "@/lib/api";
import React, {
	ReactNode,
	createContext,
	useContext,
	useEffect,
	useState,
} from "react";

// Profile cache to prevent unnecessary API calls on page reload
const PROFILE_CACHE_TTL = 10 * 60 * 1000; // 10 minutes

// Helper functions for profile cache persistence
const getProfileCacheTimestamp = (): number => {
	if (typeof window === "undefined") return 0;
	const timestamp = localStorage.getItem("profile_cache_timestamp");
	return timestamp ? parseInt(timestamp, 10) : 0;
};

const setProfileCacheTimestamp = (timestamp: number): void => {
	if (typeof window !== "undefined") {
		localStorage.setItem("profile_cache_timestamp", timestamp.toString());
	}
};

interface AuthContextType {
	user: User | null;
	isAuthenticated: boolean;
	isLoading: boolean;
	login: (msisdn: string, password: string) => Promise<void>;
	register: (name: string, msisdn: string, password: string) => Promise<void>;
	logout: () => Promise<void>;
	refreshProfile: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
	children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
	const [user, setUser] = useState<User | null>(null);
	const [isLoading, setIsLoading] = useState(true);

	const isAuthenticated = !!user;

	// Initialize auth state on mount
	useEffect(() => {
		const initializeAuth = async () => {
			try {
				const { token, user: storedUser } = auth.getAuthData();

				if (token && storedUser) {
					// Set stored user immediately for better UX
					setUser(storedUser);

					// Only refresh profile from server if cache is expired
					const now = Date.now();
					const lastProfileFetch = getProfileCacheTimestamp();
					const profileCacheExpired =
						now - lastProfileFetch > PROFILE_CACHE_TTL;

					if (process.env.NODE_ENV === "development") {
						console.log("Profile Cache Check:", {
							lastFetch: new Date(lastProfileFetch).toISOString(),
							cacheAge: Math.round(
								(now - lastProfileFetch) / 1000
							),
							expired: profileCacheExpired,
						});
					}

					if (profileCacheExpired) {
						if (process.env.NODE_ENV === "development") {
							console.log("🔄 Fetching fresh profile from API");
						}
						try {
							const freshUser = await authAPI.getProfile();
							setUser(freshUser);
							// Update stored user data and cache timestamp
							localStorage.setItem(
								"user",
								JSON.stringify(freshUser)
							);
							setProfileCacheTimestamp(now);
						} catch (error) {
							// If profile fetch fails, keep using stored user
							console.warn("Failed to refresh profile:", error);
							// Keep the stored user, no need to update timestamp
						}
					} else {
						if (process.env.NODE_ENV === "development") {
							console.log("✅ Using cached profile data");
						}
					}
				}
			} catch (error) {
				console.error("Auth initialization error:", error);
				// Clear invalid auth data
				auth.clearAuthData();
			} finally {
				setIsLoading(false);
			}
		};

		initializeAuth();
	}, []);

	const login = async (msisdn: string, password: string): Promise<void> => {
		setIsLoading(true);
		try {
			const formattedMSISDN = auth.formatMSISDN(msisdn);
			const authResponse = await authAPI.login({
				msisdn: formattedMSISDN,
				password,
			});

			// Store auth data
			auth.setAuthData(authResponse);
			setUser(authResponse.user);
		} catch (error) {
			throw error; // Re-throw to be handled by the component
		} finally {
			setIsLoading(false);
		}
	};

	const register = async (
		name: string,
		msisdn: string,
		password: string
	): Promise<void> => {
		setIsLoading(true);
		try {
			const formattedMSISDN = auth.formatMSISDN(msisdn);
			const authResponse = await authAPI.register({
				name,
				msisdn: formattedMSISDN,
				password,
			});

			// Store auth data
			auth.setAuthData(authResponse);
			setUser(authResponse.user);
		} catch (error) {
			throw error; // Re-throw to be handled by the component
		} finally {
			setIsLoading(false);
		}
	};

	const logout = async (): Promise<void> => {
		setIsLoading(true);
		try {
			await authAPI.logout();
		} catch (error) {
			console.error("Logout error:", error);
		} finally {
			// Always clear local state regardless of API success
			auth.clearAuthData();
			setUser(null);
			setIsLoading(false);
		}
	};

	const refreshProfile = async (): Promise<void> => {
		if (!auth.isAuthenticated()) return;

		try {
			const freshUser = await authAPI.getProfile();
			setUser(freshUser);
			localStorage.setItem("user", JSON.stringify(freshUser));
			setProfileCacheTimestamp(Date.now()); // Update cache timestamp
		} catch (error) {
			console.error("Failed to refresh profile:", error);
			// If refresh fails and it's an auth error, the interceptor will handle logout
		}
	};

	const value: AuthContextType = {
		user,
		isAuthenticated,
		isLoading,
		login,
		register,
		logout,
		refreshProfile,
	};

	return (
		<AuthContext.Provider value={value}>{children}</AuthContext.Provider>
	);
}

export function useAuth(): AuthContextType {
	const context = useContext(AuthContext);
	if (context === undefined) {
		throw new Error("useAuth must be used within an AuthProvider");
	}
	return context;
}

// Higher-order component for protecting routes
export function withAuth<P extends object>(
	WrappedComponent: React.ComponentType<P>
): React.ComponentType<P> {
	return function AuthenticatedComponent(props: P) {
		const { isAuthenticated, isLoading } = useAuth();

		if (isLoading) {
			return (
				<div className="min-h-screen flex items-center justify-center">
					<div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
				</div>
			);
		}

		if (!isAuthenticated) {
			// Redirect to auth page if not authenticated
			if (typeof window !== "undefined") {
				window.location.href = "/auth";
			}
			return null;
		}

		return <WrappedComponent {...props} />;
	};
}
