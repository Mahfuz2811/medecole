import { cn } from "@/styles/design-tokens";
import { BookOpen, X } from "lucide-react";
import Link from "next/link";
import { LoginModalProps } from "./types";

export function LoginModal({ package: pkg, onClose }: LoginModalProps) {
	return (
		<div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
			<div className="bg-white rounded-lg max-w-md w-full p-6 relative">
				{/* Close Button */}
				<button
					onClick={onClose}
					className="absolute top-4 right-4 text-gray-400 hover:text-gray-600"
				>
					<X size={24} />
				</button>

				{/* Modal Content */}
				<div className="text-center">
					<div className="mb-4">
						<div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
							<BookOpen className="w-8 h-8 text-blue-600" />
						</div>
						<h2 className="text-xl font-bold text-gray-900 mb-2">
							Ready to Enroll?
						</h2>
						<p className="text-gray-600 mb-4">
							Sign in to enroll in{" "}
							<span className="font-semibold">{pkg.name}</span>
						</p>
					</div>

					{/* Package Preview */}
					<div className="bg-gray-50 rounded-lg p-4 mb-6">
						<div className="flex items-center justify-between text-sm">
							<span className="text-gray-600">Package Type:</span>
							<span
								className={cn(
									"px-2 py-1 rounded text-xs font-medium",
									pkg.package_type === "FREE"
										? "bg-green-100 text-green-800"
										: "bg-blue-100 text-blue-800"
								)}
							>
								{pkg.package_type === "FREE"
									? "Free"
									: "Premium"}
							</span>
						</div>
						{pkg.package_type === "PREMIUM" && (
							<div className="flex items-center justify-between text-sm mt-2">
								<span className="text-gray-600">Price:</span>
								<span className="font-semibold">
									৳{pkg.price.toLocaleString()}
								</span>
							</div>
						)}
						<div className="flex items-center justify-between text-sm mt-2">
							<span className="text-gray-600">Total Exams:</span>
							<span className="font-semibold">
								{pkg.total_exams}
							</span>
						</div>
					</div>

					{/* Action Buttons */}
					<div className="space-y-3">
						<Link
							href={`/auth?redirect=/packages/${pkg.slug}&action=enroll`}
							className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-blue-700 transition-colors block text-center"
						>
							Sign In to Continue
						</Link>
						<p className="text-sm text-gray-500">
							Don&apos;t have an account?{" "}
							<Link
								href={`/auth?mode=register&redirect=/packages/${pkg.slug}&action=enroll`}
								className="text-blue-600 hover:text-blue-700 font-medium"
							>
								Create one now
							</Link>
						</p>
					</div>
				</div>
			</div>
		</div>
	);
}
