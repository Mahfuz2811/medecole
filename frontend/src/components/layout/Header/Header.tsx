"use client";

import { useAuth } from "@/contexts/AuthContext";
import { cn } from "@/styles/design-tokens";
import Image from "next/image";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";

interface HeaderProps {
	userName?: string;
	userId?: string;
	isLoggedIn?: boolean;
}

export const Header: React.FC<HeaderProps> = ({
	userName = "Guest User",
	userId = "ID-0000",
	isLoggedIn = false,
}) => {
	const { logout } = useAuth();
	const [isLoggingOut, setIsLoggingOut] = useState(false);
	const [isMenuOpen, setIsMenuOpen] = useState(false);
	const [isAnimating, setIsAnimating] = useState(false);
	const menuRef = useRef<HTMLDivElement>(null);

	const handleLogout = async () => {
		setIsLoggingOut(true);
		try {
			await logout();
		} catch (error) {
			console.error("Logout error:", error);
		} finally {
			setIsLoggingOut(false);
		}
	};

	const toggleMenu = () => {
		if (isMenuOpen) {
			// Start closing animation
			setIsAnimating(true);
			setTimeout(() => {
				setIsMenuOpen(false);
				setIsAnimating(false);
			}, 200); // Match animation duration
		} else {
			// Open menu
			setIsMenuOpen(true);
			setIsAnimating(true);
			setTimeout(() => {
				setIsAnimating(false);
			}, 200);
		}
	};

	// Close menu when clicking outside
	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			if (
				menuRef.current &&
				!menuRef.current.contains(event.target as Node) &&
				isMenuOpen
			) {
				// Start closing animation
				setIsAnimating(true);
				setTimeout(() => {
					setIsMenuOpen(false);
					setIsAnimating(false);
				}, 200);
			}
		};

		if (isMenuOpen) {
			document.addEventListener("mousedown", handleClickOutside);
		}

		return () => {
			document.removeEventListener("mousedown", handleClickOutside);
		};
	}, [isMenuOpen]);

	return (
		<header
			className={cn(
				"w-full bg-white shadow-sm border-b border-gray-100",
				"px-4 py-3 md:px-6 lg:px-8"
			)}
		>
			<div className="flex items-center justify-between max-w-7xl mx-auto">
				{/* Left side - Logo */}
				<Link
					href="/"
					className="flex items-center space-x-3 hover:opacity-80 transition-opacity duration-200"
				>
					<div className="flex items-center space-x-2">
						<div className="w-8 h-8 md:w-10 md:h-10 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
							<span className="text-white font-bold text-sm md:text-lg">
								Q
							</span>
						</div>
						<div className="hidden sm:block">
							<h1 className="text-lg md:text-xl font-bold text-gray-900">
								Medecole
							</h1>
							<p className="text-xs text-gray-500 hidden md:block">
								MCQ Practice Platform
							</p>
						</div>
					</div>
				</Link>

				{/* Right side - Authentication or Dashboard Menu */}
				{isLoggedIn ? (
					<div className="flex items-center space-x-4">
						{/* User Info - Hidden on small screens */}
						<div className="hidden md:flex items-center space-x-3">
							<div className="text-right">
								<p className="text-sm font-medium text-gray-900">
									{userName}
								</p>
								<p className="text-xs text-gray-500">
									{userId}
								</p>
							</div>
							<div className="w-8 h-8 rounded-full overflow-hidden">
								<Image
									src="/avatar.png"
									alt={userName}
									width={32}
									height={32}
									className="object-cover w-full h-full"
								/>
							</div>
						</div>

						{/* Menu Button */}
						<div className="relative" ref={menuRef}>
							<button
								onClick={toggleMenu}
								className="flex items-center space-x-1 p-2 rounded-lg hover:bg-gray-100 transition-all duration-200 ease-in-out"
								aria-label="Open menu"
							>
								<div className="w-6 h-6 md:hidden">
									<div className="w-6 h-6 rounded-full overflow-hidden">
										<Image
											src="/avatar.png"
											alt={userName}
											width={24}
											height={24}
											className="object-cover w-full h-full"
										/>
									</div>
								</div>
								<svg
									className="w-5 h-5 text-gray-600 transform transition-transform duration-200 ease-in-out"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
									style={{
										transform: isMenuOpen
											? "rotate(180deg)"
											: "rotate(0deg)",
									}}
								>
									<path
										strokeLinecap="round"
										strokeLinejoin="round"
										strokeWidth={2}
										d={
											isMenuOpen
												? "M6 18L18 6M6 6l12 12"
												: "M4 6h16M4 12h16M4 18h16"
										}
									/>
								</svg>
							</button>

							{/* Dropdown Menu */}
							{(isMenuOpen || isAnimating) && (
								<div
									className={`absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-50 transition-all duration-200 ease-in-out ${
										isMenuOpen && !isAnimating
											? "transform translate-x-0 opacity-100 scale-100"
											: "transform translate-x-4 opacity-0 scale-95"
									}`}
									style={{
										transformOrigin: "top right",
									}}
								>
									{/* User info */}
									<div className="px-4 py-3 border-b border-gray-100">
										<p className="text-sm font-medium text-gray-900">
											{userName}
										</p>
										<p className="text-xs text-gray-500">
											{userId}
										</p>
									</div>

									{/* Logout Button */}
									<div className="py-1">
										<button
											onClick={handleLogout}
											disabled={isLoggingOut}
											className="flex items-center w-full px-4 py-3 text-sm text-red-600 hover:bg-red-50 disabled:opacity-50 transition-colors duration-150 focus:outline-none focus:bg-red-50"
										>
											<svg
												className="w-4 h-4 mr-3"
												fill="none"
												stroke="currentColor"
												viewBox="0 0 24 24"
											>
												<path
													strokeLinecap="round"
													strokeLinejoin="round"
													strokeWidth={2}
													d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
												/>
											</svg>
											{isLoggingOut
												? "Logging out..."
												: "Logout"}
										</button>
									</div>
								</div>
							)}
						</div>
					</div>
				) : (
					/* Guest Login Button */
					<div className="flex items-center space-x-3">
						<Link href="/auth">
							<button
								className={cn(
									"bg-blue-600 hover:bg-blue-700 text-white",
									"px-4 py-2 rounded-lg text-sm font-medium",
									"transition-colors duration-200",
									"flex items-center space-x-2"
								)}
							>
								<svg
									className="w-4 h-4"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										strokeLinecap="round"
										strokeLinejoin="round"
										strokeWidth={2}
										d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
									/>
								</svg>
								<span>Login</span>
							</button>
						</Link>
					</div>
				)}
			</div>
		</header>
	);
};
