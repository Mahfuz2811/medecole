import { cn, designTokens } from "@/styles/design-tokens";
import { useState } from "react";
import { LoginFormData } from "./types";
import { formatMSISDN, validateMSISDN } from "./utils";

interface LoginFormProps {
	isLoading: boolean;
	onSubmit: (data: LoginFormData) => Promise<void>;
	onError: (error: string) => void;
}

export default function LoginForm({
	isLoading,
	onSubmit,
	onError,
}: LoginFormProps) {
	const [formData, setFormData] = useState<LoginFormData>({
		msisdn: "",
		password: "",
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
					placeholder="Enter your password"
					value={formData.password}
					onChange={(e) =>
						setFormData({
							...formData,
							password: e.target.value,
						})
					}
					className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors"
				/>
			</div>

			<div className="flex items-center justify-between">
				<div className="flex items-center">
					<input
						id="remember-me"
						type="checkbox"
						className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
					/>
					<label
						htmlFor="remember-me"
						className="ml-2 block text-sm text-gray-700"
					>
						Remember me
					</label>
				</div>
				<button
					type="button"
					className="text-sm text-blue-600 hover:text-blue-500"
				>
					Forgot password?
				</button>
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
				{isLoading ? "Signing in..." : "Sign In"}
			</button>
		</form>
	);
}
