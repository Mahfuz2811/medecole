"use client";

import { useAuth } from "@/contexts/AuthContext";
import { Header } from "./Header";

export const StickyHeader = () => {
	const { user, isAuthenticated } = useAuth();

	return (
		<div className="sticky top-0 left-0 right-0 z-50">
			<Header
				userName={user?.name || "Guest"}
				userId={user?.id ? `ID-${user.id}` : ""}
				isLoggedIn={isAuthenticated}
			/>
		</div>
	);
};
