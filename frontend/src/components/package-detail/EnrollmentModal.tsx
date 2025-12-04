import { useEnrollment } from "@/hooks/useEnrollment";
import type { CouponValidationResponse } from "@/lib/api/enrollment-types";
import type { PackageResponse } from "@/lib/api/types";
import { cn } from "@/styles/design-tokens";
import {
	AlertCircle,
	Calculator,
	CheckCircle,
	Loader2,
	Tag,
	X,
} from "lucide-react";
import { useEffect, useState } from "react";

interface EnrollmentModalProps {
	package: PackageResponse;
	isOpen: boolean;
	onClose: () => void;
	onSuccess?: (enrollmentId: number) => void;
}

export function EnrollmentModal({
	package: pkg,
	isOpen,
	onClose,
	onSuccess,
}: EnrollmentModalProps) {
	const { enrollInPackage, validateCoupon, loading, error, clearError } =
		useEnrollment();
	const [couponCode, setCouponCode] = useState(pkg.coupon_code || "");
	const [couponValidation, setCouponValidation] =
		useState<CouponValidationResponse | null>(null);
	const [validatingCoupon, setValidatingCoupon] = useState(false);
	const [enrollmentStep, setEnrollmentStep] = useState<
		"input" | "processing" | "success" | "error"
	>("input");

	// Reset state when modal opens/closes
	useEffect(() => {
		if (isOpen) {
			setCouponCode(pkg.coupon_code || "");
			setCouponValidation(null);
			setValidatingCoupon(false);
			setEnrollmentStep("input");
			clearError();
		}
	}, [isOpen, pkg.coupon_code, clearError]);

	// Validate coupon when user stops typing
	useEffect(() => {
		if (!couponCode.trim()) {
			setCouponValidation(null);
			return;
		}

		const timeoutId = setTimeout(async () => {
			setValidatingCoupon(true);
			const validation = await validateCoupon({
				coupon_code: couponCode,
				package_id: pkg.id,
			});
			setCouponValidation(validation);
			setValidatingCoupon(false);
		}, 500);

		return () => clearTimeout(timeoutId);
	}, [couponCode, pkg.id, validateCoupon]);

	const handleEnroll = async () => {
		setEnrollmentStep("processing");
		clearError();

		const enrollmentRequest = {
			package_id: pkg.id,
			...(couponCode.trim() && { coupon_code: couponCode.trim() }),
		};

		const response = await enrollInPackage(enrollmentRequest);

		if (response) {
			setEnrollmentStep("success");
			setTimeout(() => {
				onSuccess?.(response.id);
				onClose();
			}, 2000);
		} else {
			setEnrollmentStep("error");
		}
	};

	const calculateFinalPrice = () => {
		if (couponValidation?.valid && couponValidation.price_calculation) {
			return couponValidation.price_calculation.final_price;
		}
		return pkg.price;
	};

	if (!isOpen) return null;

	return (
		<div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
			<div className="bg-white rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-y-auto">
				{/* Header */}
				<div className="flex items-center justify-between p-6 border-b">
					<h2 className="text-xl font-semibold text-gray-900">
						Enroll in {pkg.name}
					</h2>
					<button
						onClick={onClose}
						disabled={enrollmentStep === "processing"}
						className="text-gray-400 hover:text-gray-600 disabled:opacity-50"
					>
						<X size={24} />
					</button>
				</div>

				{/* Content */}
				<div className="p-6">
					{enrollmentStep === "input" && (
						<>
							{/* Package Info */}
							<div className="mb-6">
								<div className="flex items-center justify-between mb-2">
									<span className="text-gray-600">
										Package Type:
									</span>
									<span
										className={cn(
											"px-2 py-1 rounded text-sm font-medium",
											pkg.package_type === "FREE"
												? "bg-green-100 text-green-800"
												: "bg-blue-100 text-blue-800"
										)}
									>
										{pkg.package_type}
									</span>
								</div>

								<div className="flex items-center justify-between mb-2">
									<span className="text-gray-600">
										Original Price:
									</span>
									<span className="font-semibold">
										{pkg.package_type === "FREE"
											? "Free"
											: `৳${pkg.price.toLocaleString()}`}
									</span>
								</div>

								{couponValidation?.valid &&
									couponValidation.price_calculation && (
										<div className="mt-4 p-3 bg-green-50 rounded-lg">
											<div className="flex items-center gap-2 mb-2">
												<Calculator
													size={16}
													className="text-green-600"
												/>
												<span className="text-sm font-medium text-green-800">
													Coupon Applied!
												</span>
											</div>
											<div className="space-y-1 text-sm">
												<div className="flex justify-between">
													<span>Discount:</span>
													<span className="text-green-600">
														-৳
														{couponValidation.price_calculation.discount_amount.toLocaleString()}
														(
														{
															couponValidation
																.price_calculation
																.discount_percentage
														}
														%)
													</span>
												</div>
												<div className="flex justify-between font-semibold">
													<span>Final Price:</span>
													<span className="text-green-600">
														৳
														{couponValidation.price_calculation.final_price.toLocaleString()}
													</span>
												</div>
											</div>
										</div>
									)}
							</div>

							{/* Coupon Input */}
							{pkg.package_type === "PREMIUM" && (
								<div className="mb-6">
									<label className="block text-sm font-medium text-gray-700 mb-2">
										Coupon Code (Optional)
									</label>
									<div className="relative">
										<div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
											<Tag
												size={16}
												className="text-gray-400"
											/>
										</div>
										<input
											type="text"
											value={couponCode}
											onChange={(e) =>
												setCouponCode(
													e.target.value.toUpperCase()
												)
											}
											placeholder="Enter coupon code"
											className="block w-full pl-10 pr-10 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
											disabled={
												loading || validatingCoupon
											}
										/>
										{validatingCoupon && (
											<div className="absolute inset-y-0 right-0 pr-3 flex items-center">
												<Loader2
													size={16}
													className="text-gray-400 animate-spin"
												/>
											</div>
										)}
									</div>

									{/* Coupon Validation Message */}
									{couponValidation && !validatingCoupon && (
										<div
											className={cn(
												"mt-2 p-2 rounded text-sm flex items-center gap-2",
												couponValidation.valid
													? "bg-green-50 text-green-700"
													: "bg-red-50 text-red-700"
											)}
										>
											{couponValidation.valid ? (
												<CheckCircle size={16} />
											) : (
												<AlertCircle size={16} />
											)}
											{couponValidation.message}
										</div>
									)}
								</div>
							)}

							{/* Error Display */}
							{error && (
								<div className="mb-4 p-3 bg-red-50 rounded-lg flex items-center gap-2">
									<AlertCircle
										size={16}
										className="text-red-600"
									/>
									<span className="text-sm text-red-700">
										{error}
									</span>
								</div>
							)}

							{/* Action Buttons */}
							<div className="flex gap-3">
								<button
									onClick={onClose}
									className="flex-1 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors"
									disabled={loading}
								>
									Cancel
								</button>
								<button
									onClick={handleEnroll}
									disabled={
										loading ||
										validatingCoupon ||
										(couponCode.trim() !== "" &&
											!couponValidation?.valid)
									}
									className={cn(
										"flex-1 px-4 py-2 rounded-md font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed",
										pkg.package_type === "FREE"
											? "bg-green-600 hover:bg-green-700 text-white"
											: "bg-blue-600 hover:bg-blue-700 text-white"
									)}
								>
									{loading ? (
										<div className="flex items-center justify-center gap-2">
											<Loader2
												size={16}
												className="animate-spin"
											/>
											Enrolling...
										</div>
									) : (
										`Enroll${
											pkg.package_type === "PREMIUM"
												? ` - ৳${calculateFinalPrice().toLocaleString()}`
												: " Free"
										}`
									)}
								</button>
							</div>
						</>
					)}

					{enrollmentStep === "processing" && (
						<div className="text-center py-8">
							<Loader2
								size={48}
								className="mx-auto mb-4 text-blue-600 animate-spin"
							/>
							<h3 className="text-lg font-medium text-gray-900 mb-2">
								Processing Enrollment...
							</h3>
							<p className="text-gray-600">
								Please wait while we process your enrollment.
							</p>
						</div>
					)}

					{enrollmentStep === "success" && (
						<div className="text-center py-8">
							<CheckCircle
								size={48}
								className="mx-auto mb-4 text-green-600"
							/>
							<h3 className="text-lg font-medium text-gray-900 mb-2">
								Enrollment Successful!
							</h3>
							<p className="text-gray-600">
								You have been successfully enrolled in{" "}
								{pkg.name}.
								{pkg.package_type === "PREMIUM" && (
									<>
										{" "}
										You will be redirected to complete
										payment.
									</>
								)}
							</p>
						</div>
					)}

					{enrollmentStep === "error" && (
						<div className="text-center py-8">
							<AlertCircle
								size={48}
								className="mx-auto mb-4 text-red-600"
							/>
							<h3 className="text-lg font-medium text-gray-900 mb-2">
								Enrollment Failed
							</h3>
							<p className="text-gray-600 mb-4">
								{error ||
									"Something went wrong. Please try again."}
							</p>
							<button
								onClick={() => setEnrollmentStep("input")}
								className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
							>
								Try Again
							</button>
						</div>
					)}
				</div>
			</div>
		</div>
	);
}
