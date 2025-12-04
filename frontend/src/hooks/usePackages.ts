"use client";

import { packagesAPI } from "@/lib/api/packages";
import type {
	PackageListRequest,
	PackageListResponse,
	PackageResponse,
} from "@/lib/api/types";
import { useEffect, useRef, useState } from "react";

// Simple cache for package data to prevent duplicate requests
const packageCache = new Map<string | number, PackageResponse>();

// Package cache configuration
const PACKAGE_CACHE_TTL = 10 * 60 * 1000; // 10 minutes

// Helper functions for persistent package cache
const getCachedPackage = (
	slugOrId: string | number
): PackageResponse | null => {
	if (typeof window === "undefined") return null;

	const cacheKey = `package_${slugOrId}`;
	const cached = localStorage.getItem(cacheKey);
	if (!cached) return null;

	try {
		const { data, timestamp } = JSON.parse(cached);
		const now = Date.now();
		const isExpired = now - timestamp > PACKAGE_CACHE_TTL;

		if (process.env.NODE_ENV === "development") {
			console.log(`Package Cache Check [${slugOrId}]:`, {
				cacheAge: Math.round((now - timestamp) / 1000),
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

const setCachedPackage = (
	slugOrId: string | number,
	data: PackageResponse
): void => {
	if (typeof window === "undefined") return;

	const cacheKey = `package_${slugOrId}`;
	const cacheData = {
		data,
		timestamp: Date.now(),
	};

	try {
		localStorage.setItem(cacheKey, JSON.stringify(cacheData));
	} catch {
		// Ignore localStorage errors
	}
};

// Hook for fetching packages list
export function usePackages(params?: PackageListRequest) {
	const [data, setData] = useState<PackageListResponse | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	// Extract individual param values to create stable dependencies
	const page = params?.page;
	const perPage = params?.per_page;
	const type = params?.type;
	const sort = params?.sort;
	const order = params?.order;
	const search = params?.search;

	useEffect(() => {
		const fetchPackages = async () => {
			try {
				setLoading(true);
				setError(null);
				const requestParams = {
					...(page && { page }),
					...(perPage && { per_page: perPage }),
					...(type && { type }),
					...(sort && { sort }),
					...(order && { order }),
					...(search && { search }),
				};
				const response = await packagesAPI.getPackages(requestParams);
				setData(response);
			} catch (err) {
				setError(
					err instanceof Error
						? err.message
						: "Failed to fetch packages"
				);
			} finally {
				setLoading(false);
			}
		};

		fetchPackages();
	}, [page, perPage, type, sort, order, search]); // Use individual values as dependencies

	return {
		data,
		loading,
		error,
		refetch: () => {
			const fetchPackages = async () => {
				try {
					setLoading(true);
					setError(null);
					const response = await packagesAPI.getPackages(params);
					setData(response);
				} catch (err) {
					setError(
						err instanceof Error
							? err.message
							: "Failed to fetch packages"
					);
				} finally {
					setLoading(false);
				}
			};
			fetchPackages();
		},
	};
}

// Hook for fetching single package with simple caching
export function usePackage(slugOrId: string | number) {
	const [data, setData] = useState<PackageResponse | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const hasFetchedRef = useRef(false);

	useEffect(() => {
		// Skip if no slug/id provided
		if (!slugOrId) {
			setLoading(false);
			return;
		}

		// Check memory cache first (fastest)
		let cachedData = packageCache.get(slugOrId);

		// If not in memory, check localStorage (persistent)
		if (!cachedData) {
			const persistedData = getCachedPackage(slugOrId);
			if (persistedData) {
				cachedData = persistedData;
				// Restore to memory cache
				packageCache.set(slugOrId, persistedData);
			}
		}

		if (cachedData) {
			if (process.env.NODE_ENV === "development") {
				console.log("✅ Using cached package data");
			}
			setData(cachedData);
			setLoading(false);
			hasFetchedRef.current = true;
			return;
		}

		// Prevent double fetching with ref
		if (hasFetchedRef.current) {
			return;
		}

		hasFetchedRef.current = true;

		const fetchPackage = async () => {
			try {
				if (process.env.NODE_ENV === "development") {
					console.log("🔄 Fetching fresh package data from API");
				}
				setLoading(true);
				setError(null);

				const response = await packagesAPI.getPackage(slugOrId);

				// Cache the response in both memory and localStorage
				packageCache.set(slugOrId, response);
				setCachedPackage(slugOrId, response);
				setData(response);
			} catch (err) {
				setError(
					err instanceof Error
						? err.message
						: "Failed to fetch package"
				);
			} finally {
				setLoading(false);
			}
		};

		fetchPackage();

		// Reset the ref when component unmounts or slugOrId changes
		return () => {
			hasFetchedRef.current = false;
		};
	}, [slugOrId]);

	return { data, loading, error };
}
