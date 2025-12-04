import { cn } from "@/styles/design-tokens";
import { FormType } from "./types";

interface FormToggleProps {
	formType: FormType;
	onToggle: (type: FormType) => void;
	onError: (error: string) => void;
}

export default function FormToggle({
	formType,
	onToggle,
	onError,
}: FormToggleProps) {
	const handleToggle = (type: FormType) => {
		onToggle(type);
		onError(""); // Clear errors when switching forms
	};

	return (
		<div className="flex bg-gray-100 rounded-lg p-1 mb-6">
			<button
				onClick={() => handleToggle("login")}
				className={cn(
					"flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all",
					formType === "login"
						? "bg-white text-blue-600 shadow-sm"
						: "text-gray-600 hover:text-gray-900"
				)}
			>
				Login
			</button>
			<button
				onClick={() => handleToggle("register")}
				className={cn(
					"flex-1 py-2 px-4 rounded-md text-sm font-medium transition-all",
					formType === "register"
						? "bg-white text-blue-600 shadow-sm"
						: "text-gray-600 hover:text-gray-900"
				)}
			>
				Register
			</button>
		</div>
	);
}
