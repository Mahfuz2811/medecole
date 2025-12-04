import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { FormType } from "./types";

export function useFormType() {
	const searchParams = useSearchParams();

	// Get initial form type from URL parameter
	const initialFormType =
		searchParams.get("mode") === "register" ? "register" : "login";
	const [formType, setFormType] = useState<FormType>(initialFormType);

	// Update form type when URL parameters change
	useEffect(() => {
		const mode = searchParams.get("mode");
		if (mode === "register") {
			setFormType("register");
		} else if (mode === "login" || !mode) {
			setFormType("login");
		}
	}, [searchParams]);

	return { formType, setFormType };
}
