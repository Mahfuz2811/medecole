import { cn, designTokens } from "@/styles/design-tokens";
import { useState } from "react";
import { RegisterFormData } from "./types";
import {
	formatMSISDN,
	validateMSISDN,
	validatePassword,
	validatePasswordConfirmation,
} from "./utils";

interface RegisterFormProps {
	isLoading: boolean;
	onSubmit: (data: RegisterFormData) => Promise<void>;
	onError: (error: string) => void;
}

export default function RegisterForm({
	isLoading,
	onSubmit,
	onError,
}: RegisterFormProps) {
	const [formData, setFormData] = useState<RegisterFormData>({
		name: "",
		msisdn: "",
		password: "",
		confirmPassword: "",
	});

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		onError("");

		// Validate MSISDN format
		const msisdnError = validateMSISDN(formData.msisdn);
		if (msisdnError) {
			onError(msisdnError);
			return;
		}

		// Validate password
		const passwordError = validatePassword(formData.password);
		if (passwordError) {
			onError(passwordError);
			return;
		}

		// Validate password confirmation
		const confirmError = validatePasswordConfirmation(
			formData.password,
			formData.confirmPassword
		);
		if (confirmError) {
			onError(confirmError);
			return;
		}

		await onSubmit(formData);
	};

	const handleMSISDNChange = (value: string) => {
		const formatted = formatMSISDN(value);
		setFormData({ ...formData, msisdn: formatted });
	};

	return (
		<form onSubmit={handleSubmit} className="space-y-4">
			<div>
				<label className="block text-sm font-medium text-gray-700 mb-1">
					Full Name
				</label>
				<input
					type="text"
					required
					placeholder="Enter your full name"
					value={formData.name}
					onChange={(e) =>
						setFormData({
							...formData,
							name: e.target.value,
						})
					}
					className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors"
				/>
			</div>

			<div>
				<label className="block text-sm font-medium text-gray-700 mb-1">
					Phone Number (MSISDN)
				</label>
				<input
					type="tel"
					required
					placeholder="01XXXXXXXXX"
					value={formData.msisdn}
					onChange={(e) => handleMSISDNChange(e.target.value)}
					className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors"
				/>
			</div>

			<div>
				<label className="block text-sm font-medium text-gray-700 mb-1">
					Password
				</label>
				<input
					type="password"
					required
					placeholder="Create a password"
					value={formData.password}
					onChange={(e) =>
						setFormData({
							...formData,
							password: e.target.value,
						})
					}
					className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors"
				/>
				<p className="text-xs text-gray-500 mt-1">
					Password should be at least 6 characters long
				</p>
			</div>

			<div>
				<label className="block text-sm font-medium text-gray-700 mb-1">
					Confirm Password
				</label>
				<input
					type="password"
					required
					placeholder="Confirm your password"
					value={formData.confirmPassword}
					onChange={(e) =>
						setFormData({
							...formData,
							confirmPassword: e.target.value,
						})
					}
					className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors"
				/>
			</div>

			<div className="flex items-center">
				<input
					id="agree-terms"
					type="checkbox"
					required
					className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
				/>
				<label
					htmlFor="agree-terms"
					className="ml-2 block text-sm text-gray-700"
				>
					I agree to the{" "}
					<button
						type="button"
						className="text-blue-600 hover:text-blue-500"
					>
						Terms & Conditions
					</button>
				</label>
			</div>

			<button
				type="submit"
				disabled={isLoading}
				className={cn(
					designTokens.components.button.primary,
					"w-full",
					isLoading ? "opacity-50 cursor-not-allowed" : ""
				)}
			>
				{isLoading ? "Creating account..." : "Create Account"}
			</button>
		</form>
	);
}
