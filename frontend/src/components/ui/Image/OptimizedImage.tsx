"use client";

import { PackageImageResponse } from "@/lib/api/types";
import { cn } from "@/styles/design-tokens";
import Image from "next/image";
import { useState } from "react";

interface OptimizedImageProps {
	images: PackageImageResponse;
	className?: string;
	alt?: string;
	fill?: boolean;
	width?: number;
	height?: number;
	priority?: boolean;
	sizes?: string;
	variant?: "thumbnail" | "mobile" | "tablet" | "desktop" | "original";
	objectFit?: "cover" | "contain" | "fill" | "none" | "scale-down";
}

export function OptimizedImage({
	images,
	className,
	alt,
	fill = false,
	width,
	height,
	priority = false,
	sizes = "(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw",
	objectFit = "contain",
}: OptimizedImageProps) {
	const [imageError, setImageError] = useState(false);
	const [isLoading, setIsLoading] = useState(true);

	// Simple URL selection with fallback
	const getImageUrl = () => {
		if (imageError) {
			return "/images/banner.png";
		}

		// Always use the original URL for now to test
		return images.original || "/images/banner.png";
	};

	const handleImageError = () => {
		console.log("❌ Image failed:", getImageUrl());
		setImageError(true);
		setIsLoading(false);
	};

	const handleImageLoad = () => {
		console.log("✅ Image loaded:", getImageUrl());
		setIsLoading(false);
	};

	const imageUrl = getImageUrl();
	const imageAlt = alt || images.alt_text || "Package image";

	return (
		<div
			className={cn(
				"relative overflow-hidden w-full h-full",
				className || ""
			)}
		>
			{isLoading && (
				<div className="absolute inset-0 bg-gray-200 animate-pulse z-10" />
			)}

			<Image
				src={imageUrl}
				alt={imageAlt}
				fill={fill}
				width={!fill ? width : undefined}
				height={!fill ? height : undefined}
				priority={priority}
				sizes={sizes}
				className={cn(
					"transition-opacity duration-300",
					isLoading ? "opacity-0" : "opacity-100",
					fill ? `object-${objectFit}` : ""
				)}
				onError={handleImageError}
				onLoad={handleImageLoad}
				unoptimized={true}
			/>

			{imageError && (
				<div className="absolute inset-0 bg-blue-100 flex items-center justify-center z-20">
					<p className="text-xs text-gray-500">No image</p>
				</div>
			)}
		</div>
	);
}
