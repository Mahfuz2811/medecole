import axios from "axios";
import { authApiClient } from "./client";
import { AuthResponse, LoginRequest, RegisterRequest, User } from "./types";
import { auth } from "./utils";

export const authAPI = {
	// Register a new user
	register: async (data: RegisterRequest): Promise<AuthResponse> => {
		try {
			const response = await authApiClient.post("/auth/register", data);
			return response.data;
		} catch (error) {
			if (axios.isAxiosError(error)) {
				throw new Error(
					error.response?.data?.message || "Registration failed"
				);
			}
			throw new Error("Registration failed");
		}
	},

	// Login user
	login: async (data: LoginRequest): Promise<AuthResponse> => {
		try {
			const response = await authApiClient.post("/auth/login", data);
			return response.data;
		} catch (error) {
			if (axios.isAxiosError(error)) {
				throw new Error(
					error.response?.data?.message || "Login failed"
				);
			}
			throw new Error("Login failed");
		}
	},

	// Logout user
	logout: async (): Promise<void> => {
		try {
			const token = localStorage.getItem("authToken");
			if (token) {
				await authApiClient.post("/auth/logout");
			}
		} catch (error) {
			// Even if logout fails on server, we still clear local storage
			console.error("Logout error:", error);
		} finally {
			// Use the comprehensive clearAuthData function to clear all auth data and cache
			auth.clearAuthData();
		}
	},

	// Get current user profile
	getProfile: async (): Promise<User> => {
		try {
			const response = await authApiClient.get("/auth/profile");
			return response.data;
		} catch (error) {
			if (axios.isAxiosError(error)) {
				throw new Error(
					error.response?.data?.message || "Failed to get profile"
				);
			}
			throw new Error("Failed to get profile");
		}
	},

	// Google OAuth authentication with credential
	googleAuth: async (credential: string): Promise<AuthResponse> => {
		try {
			const response = await authApiClient.post(
				"/auth/google/credential",
				{ credential }
			);
			return response.data;
		} catch (error) {
			if (axios.isAxiosError(error)) {
				throw new Error(
					error.response?.data?.message ||
						"Google authentication failed"
				);
			}
			throw new Error("Google authentication failed");
		}
	},

	// Facebook OAuth authentication with access token
	facebookAuth: async (accessToken: string): Promise<AuthResponse> => {
		try {
			const response = await authApiClient.post("/auth/facebook/token", {
				access_token: accessToken,
			});
			return response.data;
		} catch (error) {
			if (axios.isAxiosError(error)) {
				throw new Error(
					error.response?.data?.message ||
						"Facebook authentication failed"
				);
			}
			throw new Error("Facebook authentication failed");
		}
	},
};
