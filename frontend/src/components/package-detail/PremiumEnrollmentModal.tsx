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

interface PremiumEnrollmentModalProps {
	package: PackageResponse;
	isOpen: boolean;
	onClose: () => void;
	onSuccess?: (enrollmentId: number) => void;
}

export function PremiumEnrollmentModal({
	package: pkg,
	isOpen,
	onClose,
	onSuccess,
}: PremiumEnrollmentModalProps) {
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
		<>
			{/* Contextual Overlay - Lighter background, not full black */}
			<div className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center p-4 z-50">
				{/* Modal - Contextual size, not full screen */}
				<div className="bg-white rounded-xl shadow-2xl max-w-lg w-full max-h-[90vh] overflow-y-auto animate-in fade-in-0 zoom-in-95 duration-200">
					{/* Header */}
					<div className="flex items-center justify-between p-6 border-b border-gray-100">
						<div>
							<h2 className="text-xl font-semibold text-gray-900">
								Enroll in Premium Package
							</h2>
							<p className="text-sm text-gray-600 mt-1">
								{pkg.name}
							</p>
						</div>
						<button
							onClick={onClose}
							disabled={enrollmentStep === "processing"}
							className="text-gray-400 hover:text-gray-600 disabled:opacity-50 p-2 hover:bg-gray-100 rounded-lg transition-colors"
						>
							<X size={20} />
						</button>
					</div>

					{/* Content */}
					<div className="p-6">
						{enrollmentStep === "input" && (
							<>
								{/* Package Info Card */}
								<div className="mb-6 p-4 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg border border-blue-100">
									<div className="flex items-center justify-between mb-3">
										<span className="text-sm font-medium text-blue-900">
											Package Details
										</span>
										<span className="px-3 py-1 bg-blue-600 text-white text-xs font-medium rounded-full">
											PREMIUM
										</span>
									</div>

									<div className="space-y-2">
										<div className="flex items-center justify-between">
											<span className="text-sm text-blue-800">
												Original Price:
											</span>
											<span className="font-semibold text-blue-900">
												৳{pkg.price.toLocaleString()}
											</span>
										</div>

										{couponValidation?.valid &&
											couponValidation.price_calculation && (
												<div className="mt-3 p-3 bg-green-50 rounded-lg border border-green-200">
													<div className="flex items-center gap-2 mb-2">
														<Calculator
															size={16}
															className="text-green-600"
														/>
														<span className="text-sm font-medium text-green-800">
															Coupon Applied
															Successfully!
														</span>
													</div>
													<div className="space-y-1 text-sm">
														<div className="flex justify-between">
															<span className="text-green-700">
																Discount:
															</span>
															<span className="text-green-700 font-medium">
																-৳
																{couponValidation.price_calculation.discount_amount.toLocaleString()}{" "}
																(
																{
																	couponValidation
																		.price_calculation
																		.discount_percentage
																}
																%)
															</span>
														</div>
														<div className="flex justify-between font-semibold text-lg border-t border-green-200 pt-2">
															<span className="text-green-800">
																Final Price:
															</span>
															<span className="text-green-800">
																৳
																{couponValidation.price_calculation.final_price.toLocaleString()}
															</span>
														</div>
													</div>
												</div>
											)}
									</div>
								</div>

								{/* Coupon Input */}
								<div className="mb-6">
									<label className="block text-sm font-medium text-gray-700 mb-3">
										Have a coupon code?{" "}
										<span className="text-gray-500">
											(Optional)
										</span>
									</label>
									<div className="relative">
										<div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
											<Tag
												size={18}
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
											placeholder="Enter your coupon code"
											className="block w-full pl-12 pr-12 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm placeholder-gray-400"
											disabled={
												loading || validatingCoupon
											}
										/>
										{validatingCoupon && (
											<div className="absolute inset-y-0 right-0 pr-4 flex items-center">
												<Loader2
													size={18}
													className="text-blue-500 animate-spin"
												/>
											</div>
										)}
									</div>

									{/* Coupon Validation Message */}
									{couponValidation && !validatingCoupon && (
										<div
											className={cn(
												"mt-3 p-3 rounded-lg text-sm flex items-center gap-3",
												couponValidation.valid
													? "bg-green-50 text-green-700 border border-green-200"
													: "bg-red-50 text-red-700 border border-red-200"
											)}
										>
											{couponValidation.valid ? (
												<CheckCircle
													size={18}
													className="text-green-600"
												/>
											) : (
												<AlertCircle
													size={18}
													className="text-red-600"
												/>
											)}
											<span>
												{couponValidation.message}
											</span>
										</div>
									)}
								</div>

								{/* Error Display */}
								{error && (
									<div className="mb-6 p-4 bg-red-50 rounded-lg flex items-center gap-3 border border-red-200">
										<AlertCircle
											size={18}
											className="text-red-600"
										/>
										<div>
											<p className="text-sm font-medium text-red-800">
												Enrollment Failed
											</p>
											<p className="text-sm text-red-700">
												{error}
											</p>
										</div>
									</div>
								)}

								{/* Benefits Section */}
								<div className="mb-6 p-4 bg-gray-50 rounded-lg">
									<h4 className="text-sm font-medium text-gray-900 mb-2">
										What you&apos;ll get:
									</h4>
									<ul className="text-sm text-gray-600 space-y-1">
										<li className="flex items-center gap-2">
											<CheckCircle
												size={14}
												className="text-green-500"
											/>
											Full access to all premium content
										</li>
										<li className="flex items-center gap-2">
											<CheckCircle
												size={14}
												className="text-green-500"
											/>
											Priority support
										</li>
										<li className="flex items-center gap-2">
											<CheckCircle
												size={14}
												className="text-green-500"
											/>
											30-day money-back guarantee
										</li>
									</ul>
								</div>

								{/* Action Buttons */}
								<div className="flex gap-3">
									<button
										onClick={onClose}
										className="flex-1 px-6 py-3 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors font-medium"
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
										className="flex-2 px-8 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
									>
										{loading ? (
											<div className="flex items-center justify-center gap-2">
												<Loader2
													size={18}
													className="animate-spin"
												/>
												Processing...
											</div>
										) : (
											`Enroll Now - ৳${calculateFinalPrice().toLocaleString()}`
										)}
									</button>
								</div>
							</>
						)}

						{enrollmentStep === "processing" && (
							<div className="text-center py-12">
								<Loader2
									size={48}
									className="mx-auto mb-4 text-blue-600 animate-spin"
								/>
								<h3 className="text-lg font-medium text-gray-900 mb-2">
									Processing Your Enrollment...
								</h3>
								<p className="text-gray-600">
									Please wait while we set up your premium
									access.
								</p>
							</div>
						)}

						{enrollmentStep === "success" && (
							<div className="text-center py-12">
								<div className="mx-auto mb-4 w-16 h-16 bg-green-100 rounded-full flex items-center justify-center">
									<CheckCircle
										size={32}
										className="text-green-600"
									/>
								</div>
								<h3 className="text-lg font-medium text-gray-900 mb-2">
									Enrollment Successful!
								</h3>
								<p className="text-gray-600 mb-4">
									Welcome to <strong>{pkg.name}</strong>!
									You&apos;ll be redirected to complete your
									payment.
								</p>
								<div className="inline-flex items-center gap-2 px-4 py-2 bg-green-50 text-green-700 rounded-lg text-sm">
									<CheckCircle size={16} />
									Premium access activated
								</div>
							</div>
						)}

						{enrollmentStep === "error" && (
							<div className="text-center py-12">
								<div className="mx-auto mb-4 w-16 h-16 bg-red-100 rounded-full flex items-center justify-center">
									<AlertCircle
										size={32}
										className="text-red-600"
									/>
								</div>
								<h3 className="text-lg font-medium text-gray-900 mb-2">
									Enrollment Failed
								</h3>
								<p className="text-gray-600 mb-6">
									{error ||
										"Something went wrong. Please try again."}
								</p>
								<div className="flex gap-3 justify-center">
									<button
										onClick={onClose}
										className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
									>
										Close
									</button>
									<button
										onClick={() =>
											setEnrollmentStep("input")
										}
										className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
									>
										Try Again
									</button>
								</div>
							</div>
						)}
					</div>
				</div>
			</div>
		</>
	);
}
