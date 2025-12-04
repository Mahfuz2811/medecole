import Link from "next/link";

export default function BackToHome() {
	return (
		<div className="text-center">
			<Link
				href="/"
				className="inline-flex items-center gap-2 text-sm font-medium text-gray-600 hover:text-indigo-600 transition-colors duration-200 group"
			>
				<svg
					className="w-4 h-4 transform group-hover:-translate-x-1 transition-transform duration-200"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						strokeLinecap="round"
						strokeLinejoin="round"
						strokeWidth={2}
						d="M10 19l-7-7m0 0l7-7m-7 7h18"
					/>
				</svg>
				Back to Home
			</Link>
		</div>
	);
}
