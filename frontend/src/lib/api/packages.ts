import axios from "axios";
import { packagesApiClient } from "./client";
import {
	ExamListResponse,
	PackageListRequest,
	PackageListResponse,
	PackageResponse,
} from "./types";

export const packagesAPI = {
	// Get all packages with optional filtering
	getPackages: async (
		params?: PackageListRequest
	): Promise<PackageListResponse> => {
		try {
			const queryParams = new URLSearchParams();

			if (params?.page)
				queryParams.append("page", params.page.toString());
			if (params?.per_page)
				queryParams.append("per_page", params.per_page.toString());
			if (params?.type) queryParams.append("type", params.type);
			if (params?.sort) queryParams.append("sort", params.sort);
			if (params?.order) queryParams.append("order", params.order);
			if (params?.search) queryParams.append("search", params.search);

			const url = `/packages${
				queryParams.toString() ? `?${queryParams.toString()}` : ""
			}`;
			const response = await packagesApiClient.get(url);
			return response.data;
		} catch (error) {
			if (axios.isAxiosError(error)) {
				throw new Error(
					error.response?.data?.message || "Failed to fetch packages"
				);
			}
			throw new Error("Failed to fetch packages");
		}
	},

	// Get package by slug or ID
	getPackage: async (slugOrId: string | number): Promise<PackageResponse> => {
		try {
			const response = await packagesApiClient.get(
				`/packages/${slugOrId}`
			);
			return response.data;
		} catch (error) {
			if (axios.isAxiosError(error)) {
				throw new Error(
					error.response?.data?.message || "Failed to fetch package"
				);
			}
			throw new Error("Failed to fetch package");
		}
	},

	// Get exams for a specific package (requires authentication)
	getPackageExams: async (packageSlug: string): Promise<ExamListResponse> => {
		try {
			// Check if user is authenticated
			const token = localStorage.getItem("authToken");
			if (!token) {
				throw new Error("Authentication required");
			}

			const response = await packagesApiClient.get(
				`/exams/${packageSlug}`,
				{
					headers: {
						Authorization: `Bearer ${token}`,
					},
				}
			);
			return response.data;
		} catch (error) {
			if (axios.isAxiosError(error)) {
				if (error.response?.status === 401) {
					throw new Error("Authentication required");
				}
				if (error.response?.status === 404) {
					throw new Error("Package not found");
				}
				throw new Error(
					error.response?.data?.message ||
						"Failed to fetch package exams"
				);
			}
			throw new Error("Failed to fetch package exams");
		}
	},
};
