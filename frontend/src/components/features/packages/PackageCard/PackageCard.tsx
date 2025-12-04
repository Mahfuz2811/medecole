import { OptimizedImage } from "@/components/ui/Image";
import { PackageResponse } from "@/lib/api/types";
import { cn, designTokens } from "@/styles/design-tokens";
import Link from "next/link";

export interface PackageCardProps {
	package: PackageResponse;
}

export const PackageCard = ({ package: pkg }: PackageCardProps) => {
	return (
		<Link href={`/packages/${pkg.slug}`} className="block group">
			<div className={designTokens.components.card.subscription}>
				{/* Image First - Full width, no padding */}
				<div className="relative w-full h-48 overflow-hidden bg-gradient-to-br from-blue-50 to-blue-100">
					<OptimizedImage
						images={pkg.images}
						fill
						variant="tablet"
						objectFit="cover"
						className="group-hover:scale-105 transition-transform duration-300"
						priority={false}
						sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
					/>
					{/* Package indicator overlay */}
					<div className="absolute top-3 right-3 bg-blue-600 text-white text-xs px-2 py-1 rounded-full font-medium z-10">
						{pkg.package_type === "FREE" ? "Free" : "Premium"}
					</div>
				</div>

				{/* Title Below Image */}
				<div className="p-4">
					<h3
						className={cn(
							designTokens.typography.lg,
							"font-semibold text-gray-900 line-clamp-2 group-hover:text-blue-600 transition-colors mb-3"
						)}
					>
						{pkg.name}
					</h3>

					{/* Package Info - Compact Design */}
					<div className="flex items-center justify-between mt-3 pt-3 border-t border-gray-100">
						<div className="flex items-center gap-3 text-xs text-gray-500">
							<span className="flex items-center gap-1">
								<div className="w-2 h-2 bg-green-500 rounded-full"></div>
								{pkg.total_exams} Exams
							</span>
							{pkg.validity_days && (
								<span className="flex items-center gap-1">
									<div className="w-2 h-2 bg-blue-500 rounded-full"></div>
									{pkg.validity_days} Days
								</span>
							)}
							{pkg.package_type === "PREMIUM" && (
								<span className="flex items-center gap-1">
									<div className="w-2 h-2 bg-purple-500 rounded-full"></div>
									${pkg.price}
								</span>
							)}
						</div>
						<div className="text-xs text-blue-600 font-medium">
							View Details →
						</div>
					</div>
				</div>
			</div>
		</Link>
	);
};
