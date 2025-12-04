import axios, { AxiosInstance } from "axios";
import { auth } from "./utils";

// Base API configuration
const createApiClient = (baseURL: string): AxiosInstance => {
	return axios.create({
		baseURL,
		timeout: 10000,
		headers: {
			"Content-Type": "application/json",
		},
	});
};

// Auth API client
export const authApiClient = createApiClient(
	process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"
);

// Packages API client
export const packagesApiClient = createApiClient(
	process.env.NEXT_PUBLIC_PACKAGES_API_URL || "http://localhost:8080/api"
);

// Request interceptor to add auth token
authApiClient.interceptors.request.use(
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

// Response interceptor to handle auth errors
authApiClient.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response?.status === 401) {
			// Clear all auth data and cache on unauthorized error
			auth.clearAuthData();
			window.location.href = "/auth";
		}
		return Promise.reject(error);
	}
);

// Add request interceptor to packages API if needed
packagesApiClient.interceptors.request.use(
	(config) => {
		// Add auth token for protected routes
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

// Add response interceptor to packages API for auth errors
packagesApiClient.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response?.status === 401) {
			// Clear all auth data and cache on unauthorized error
			auth.clearAuthData();
			window.location.href = "/auth";
		}
		return Promise.reject(error);
	}
);

// Create named object before exporting
const apiClients = { authApiClient, packagesApiClient };

export default apiClients;
