import { cn, designTokens } from "@/styles/design-tokens";
import { PackageDetailProps } from "./types";

export function PackageDescription({ package: pkg }: PackageDetailProps) {
	return (
		<div className="bg-white shadow-sm p-6">
			<h2
				className={cn(
					designTokens.typography.lg,
					"font-semibold text-gray-900 mb-4"
				)}
			>
				About This Package
			</h2>
			<div
				className={cn(
					designTokens.typography.base,
					"text-gray-700 leading-relaxed"
				)}
			>
				{pkg.description || (
					<div>
						<p className="mb-4">
							This comprehensive exam package is designed to help
							you master the subject matter through carefully
							curated questions and detailed explanations.
						</p>
						<p className="mb-4">
							Each exam in this package has been developed by
							subject matter experts and follows the latest
							curriculum standards to ensure you&apos;re getting
							the most relevant and up-to-date content.
						</p>
						<p>
							Perfect for students, professionals, and anyone
							looking to test their knowledge and improve their
							understanding of the subject.
						</p>
					</div>
				)}
			</div>
		</div>
	);
}
