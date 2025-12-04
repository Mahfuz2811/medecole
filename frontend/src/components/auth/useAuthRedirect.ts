import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

interface UseAuthRedirectProps {
	isAuthenticated: boolean;
}

export function useAuthRedirect({ isAuthenticated }: UseAuthRedirectProps) {
	const router = useRouter();
	const searchParams = useSearchParams();

	useEffect(() => {
		if (isAuthenticated) {
			const redirectUrl = searchParams.get("redirect");
			const action = searchParams.get("action");

			if (redirectUrl) {
				// If there's a specific action, append it to the redirect URL
				const finalUrl = action
					? `${redirectUrl}?action=${action}`
					: redirectUrl;
				router.push(finalUrl);
			} else {
				router.push("/dashboard");
			}
		}
	}, [isAuthenticated, router, searchParams]);
}
