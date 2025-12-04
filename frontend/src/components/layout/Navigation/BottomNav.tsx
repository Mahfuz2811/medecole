"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

export const BottomNav = () => {
	const pathname = usePathname();

	return (
		<nav className="fixed bottom-0 left-0 right-0 z-50">
			<div className="w-full max-w-4xl mx-auto">
				<div className="w-full mx-auto bg-white border border-gray-300 shadow-lg flex justify-around items-center py-2 px-6">
					<Link
						href="/"
						className={`flex flex-col items-center justify-center py-1 px-2 min-w-0 rounded-md transition-colors ${
							pathname === "/"
								? "text-blue-600 bg-blue-50"
								: "text-primary hover:bg-gray-50"
						}`}
					>
						<span className="text-xl mb-1">🏠</span>
						<span className="text-xs font-medium">Home</span>
					</Link>
					<Link
						href="/dashboard"
						className={`flex flex-col items-center justify-center py-1 px-2 min-w-0 rounded-md transition-colors ${
							pathname === "/dashboard"
								? "text-blue-600 bg-blue-50"
								: "text-primary hover:bg-gray-50"
						}`}
					>
						<span className="text-xl mb-1">📚</span>
						<span className="text-xs font-medium">Dashboard</span>
					</Link>
					<Link
						href="/profile"
						className={`flex flex-col items-center justify-center py-1 px-2 min-w-0 rounded-md transition-colors ${
							pathname === "/profile"
								? "text-blue-600 bg-blue-50"
								: "text-primary hover:bg-gray-50"
						}`}
					>
						<span className="text-xl mb-1">👤</span>
						<span className="text-xs font-medium">Profile</span>
					</Link>
				</div>
			</div>
		</nav>
	);
};
