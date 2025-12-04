import { useEnrollment } from "@/hooks/useEnrollment";
import { cn } from "@/styles/design-tokens";
import { AlertCircle, CheckCircle, Clock, Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { PremiumEnrollmentModal } from "./PremiumEnrollmentModal";
import { PackageActionProps } from "./types";

export function PackagePurchaseCard({
	package: pkg,
	isAuthenticated,
	onShowLogin,
}: PackageActionProps) {
	const router = useRouter();
	const {
		checkEnrollmentStatus,
		enrollmentStatus,
		loading: enrollmentLoading,
		enrollInPackage,
		error,
		clearError,
	} = useEnrollment(isAuthenticated ? pkg.id : undefined);

	// UI State
	const [showPremiumModal, setShowPremiumModal] = useState(false);
	const [enrollmentStep, setEnrollmentStep] = useState<
		"form" | "processing" | "success" | "error"
	>("form");

	// Check enrollment status when component mounts and user is authenticated
	useEffect(() => {
		if (isAuthenticated && enrollmentStatus === null) {
			// Only check if we have a valid token
			const token = localStorage.getItem("authToken");
			if (!token) {
				console.warn(
					"No auth token found, skipping enrollment status check"
				);
				return;
			}

			// Make API call since we don't have cached data
			const checkStatus = async () => {
				try {
					await checkEnrollmentStatus(pkg.id);
				} catch (error) {
					// Silently fail - enrollment status is optional for page display
					console.warn("Failed to check enrollment status:", error);
				}
			};
			checkStatus();
		}
	}, [isAuthenticated, pkg.id, checkEnrollmentStatus, enrollmentStatus]);

	const handleEnrollClick = () => {
		if (!isAuthenticated) {
			onShowLogin();
			return;
		}

		// Direct enrollment for FREE packages (one-click)
		if (pkg.package_type === "FREE") {
			handleDirectEnroll();
		} else {
			// Show modal for PREMIUM packages (multi-step)
			setShowPremiumModal(true);
		}
	};

	const handleDirectEnroll = async () => {
		setEnrollmentStep("processing");
		clearError();

		const enrollmentRequest = {
			package_id: pkg.id,
		};

		const response = await enrollInPackage(enrollmentRequest);

		if (response) {
			setEnrollmentStep("success");
			setTimeout(() => {
				handleEnrollmentSuccess(response.id);
				setEnrollmentStep("form"); // Reset to normal state
			}, 2000);
		} else {
			setEnrollmentStep("error");
		}
	};

	const handleEnrollmentSuccess = (enrollmentId: number) => {
		// Refresh enrollment status after successful enrollment
		checkEnrollmentStatus(pkg.id);
		console.log("Enrollment successful:", enrollmentId);
	};

	// Determine the current enrollment status
	const hasActiveEnrollment = enrollmentStatus?.has_active_enrollment;
	const enrollment = enrollmentStatus?.enrollment;
	const canAccessContent = enrollment?.can_access_content;
	const effectiveStatus = enrollment?.effective_status;

	// Render enrollment status message
	const renderEnrollmentStatus = () => {
		if (!isAuthenticated || !hasActiveEnrollment || !enrollment) {
			return null;
		}

		switch (effectiveStatus) {
			case "PAID_ACTIVE":
				return (
					<div className="flex items-center gap-2 p-3 bg-green-50 rounded-lg mb-4">
						<CheckCircle size={16} className="text-green-600" />
						<span className="text-sm text-green-800">
							You are enrolled and can access all content
						</span>
					</div>
				);
			case "PENDING_PAYMENT":
				return (
					<div className="flex items-center gap-2 p-3 bg-yellow-50 rounded-lg mb-4">
						<Clock size={16} className="text-yellow-600" />
						<span className="text-sm text-yellow-800">
							Enrollment pending payment
						</span>
					</div>
				);
			case "EXPIRED":
				return (
					<div className="flex items-center gap-2 p-3 bg-red-50 rounded-lg mb-4">
						<AlertCircle size={16} className="text-red-600" />
						<span className="text-sm text-red-800">
							Your enrollment has expired
						</span>
					</div>
				);
			case "TRIAL_ACTIVE":
				return (
					<div className="flex items-center gap-2 p-3 bg-blue-50 rounded-lg mb-4">
						<Clock size={16} className="text-blue-600" />
						<span className="text-sm text-blue-800">
							Trial access active
						</span>
					</div>
				);
			case "TRIAL_EXPIRED_PAYMENT_REQUIRED":
				return (
					<div className="flex items-center gap-2 p-3 bg-orange-50 rounded-lg mb-4">
						<AlertCircle size={16} className="text-orange-600" />
						<span className="text-sm text-orange-800">
							Trial expired - Payment required for access
						</span>
					</div>
				);
			default:
				return (
					<div className="flex items-center gap-2 p-3 bg-green-50 rounded-lg mb-4">
						<CheckCircle size={16} className="text-green-600" />
						<span className="text-sm text-green-800">
							Active enrollment
						</span>
					</div>
				);
		}
	};

	// Determine button text and behavior
	const getButtonConfig = () => {
		if (!isAuthenticated) {
			return {
				text:
					pkg.package_type === "FREE"
						? "Free! Enroll Now"
						: "Enroll Now",
				disabled: false,
				onClick: handleEnrollClick,
			};
		}

		// Show loading state while checking enrollment status
		// This prevents the blink from "Enroll Now" -> "Access Content"
		if (enrollmentLoading) {
			return {
				text: "Checking status...",
				disabled: true,
				onClick: () => {},
			};
		}

		// For FREE packages - show processing state during direct enrollment
		if (pkg.package_type === "FREE" && enrollmentStep === "processing") {
			return {
				text: "Processing...",
				disabled: true,
				onClick: () => {},
			};
		}

		if (hasActiveEnrollment && canAccessContent) {
			return {
				text: "Access Content",
				disabled: false,
				onClick: () => {
					// Navigate to package exams page
					router.push(`/exams/${pkg.slug}`);
				},
			};
		}

		if (hasActiveEnrollment && effectiveStatus === "PENDING_PAYMENT") {
			return {
				text: "Complete Payment",
				disabled: false,
				onClick: () => {
					// TODO: Navigate to payment
					console.log("Navigate to payment");
				},
			};
		}

		if (
			hasActiveEnrollment &&
			(effectiveStatus === "EXPIRED" ||
				effectiveStatus === "TRIAL_EXPIRED_PAYMENT_REQUIRED")
		) {
			return {
				text: "Renew Enrollment",
				disabled: false,
				onClick: handleEnrollClick,
			};
		}

		// Default: Not enrolled or no active enrollment
		return {
			text:
				pkg.package_type === "FREE" ? "Free! Enroll Now" : "Enroll Now",
			disabled: false,
			onClick: handleEnrollClick,
		};
	};

	const buttonConfig = getButtonConfig();

	return (
		<div className="bg-white shadow-sm sticky top-6">
			{/* Main Purchase Card */}
			<div className="p-6">
				{/* Enrollment Status */}
				{renderEnrollmentStatus()}

				{/* Price Display */}
				<div className="text-center mb-6">
					{pkg.package_type === "FREE" ? (
						<div>
							<span className="text-3xl font-bold text-green-600">
								Free
							</span>
							<p className="text-sm text-gray-600 mt-1">
								No cost, full access
							</p>
						</div>
					) : (
						<div>
							<span className="text-3xl font-bold text-gray-900">
								৳{pkg.price.toLocaleString()}
							</span>
							{pkg.coupon_code && (
								<p className="text-sm text-green-600 mt-1">
									Use code:{" "}
									<span className="font-mono font-bold">
										{pkg.coupon_code}
									</span>
								</p>
							)}
						</div>
					)}
				</div>

				{/* Direct Enrollment Feedback for FREE packages */}
				{pkg.package_type === "FREE" &&
					enrollmentStep === "success" && (
						<div className="mb-4 p-3 bg-green-50 rounded-lg flex items-center gap-2 border border-green-200 animate-in fade-in-0 duration-200">
							<CheckCircle size={16} className="text-green-600" />
							<span className="text-sm text-green-800 font-medium">
								Successfully enrolled! Welcome aboard! 🎉
							</span>
						</div>
					)}

				{pkg.package_type === "FREE" && enrollmentStep === "error" && (
					<div className="mb-4 p-3 bg-red-50 rounded-lg flex items-center gap-2 border border-red-200 animate-in fade-in-0 duration-200">
						<AlertCircle size={16} className="text-red-600" />
						<div className="flex-1">
							<span className="text-sm text-red-800 font-medium block">
								Enrollment failed
							</span>
							<span className="text-sm text-red-700">
								{error || "Please try again"}
							</span>
						</div>
						<button
							onClick={() => setEnrollmentStep("form")}
							className="text-sm text-red-700 hover:text-red-800 font-medium underline"
						>
							Retry
						</button>
					</div>
				)}

				{/* Action Button */}
				<button
					onClick={buttonConfig.onClick}
					disabled={buttonConfig.disabled}
					className={cn(
						"w-full py-3 px-4 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed",
						// If not authenticated, show default colors based on package type
						!isAuthenticated
							? pkg.package_type === "FREE"
								? "bg-green-600 hover:bg-green-700 text-white"
								: "bg-blue-600 hover:bg-blue-700 text-white"
							: // If authenticated, use enrollment-based colors
							hasActiveEnrollment && canAccessContent
							? "bg-gray-600 hover:bg-gray-700 text-white"
							: hasActiveEnrollment &&
							  effectiveStatus === "PENDING_PAYMENT"
							? "bg-orange-600 hover:bg-orange-700 text-white"
							: pkg.package_type === "FREE"
							? "bg-green-600 hover:bg-green-700 text-white"
							: "bg-blue-600 hover:bg-blue-700 text-white"
					)}
				>
					{pkg.package_type === "FREE" &&
					enrollmentStep === "processing" ? (
						<div className="flex items-center justify-center gap-2">
							<Loader2 size={16} className="animate-spin" />
							Enrolling...
						</div>
					) : (
						buttonConfig.text
					)}
				</button>

				{/* Additional Info */}
				{pkg.package_type === "PREMIUM" && !hasActiveEnrollment && (
					<p className="text-xs text-gray-500 text-center mt-3">
						30-day money-back guarantee
					</p>
				)}

				{/* Enrollment Details */}
				{hasActiveEnrollment && enrollment && (
					<div className="mt-4 p-3 bg-gray-50 rounded-lg">
						<h4 className="text-sm font-medium text-gray-900 mb-2">
							Enrollment Details
						</h4>
						<div className="space-y-1 text-xs text-gray-600">
							<div className="flex justify-between">
								<span>Enrolled:</span>
								<span>
									{new Date(
										enrollment.enrolled_at
									).toLocaleDateString()}
								</span>
							</div>
							{enrollment.expires_at && (
								<div className="flex justify-between">
									<span>Expires:</span>
									<span>
										{new Date(
											enrollment.expires_at
										).toLocaleDateString()}
									</span>
								</div>
							)}
							{enrollment.enrolled_price > 0 && (
								<div className="flex justify-between">
									<span>Price Paid:</span>
									<span>
										৳
										{enrollment.enrolled_price.toLocaleString()}
									</span>
								</div>
							)}
						</div>
					</div>
				)}
			</div>

			{/* Premium Enrollment Modal */}
			<PremiumEnrollmentModal
				package={pkg}
				isOpen={showPremiumModal}
				onClose={() => setShowPremiumModal(false)}
				onSuccess={handleEnrollmentSuccess}
			/>
		</div>
	);
}
