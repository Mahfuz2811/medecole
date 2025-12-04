// Design System Tokens
export const designTokens = {
	// Color System
	colors: {
		primary: {
			50: "bg-blue-50",
			100: "bg-blue-100",
			500: "bg-blue-500",
			600: "bg-blue-600",
			800: "bg-blue-800",
		},
		text: {
			primary: "text-gray-900",
			secondary: "text-gray-600",
			accent: "text-blue-600",
			inverse: "text-white",
		},
		background: {
			primary: "bg-white",
			secondary: "bg-gray-50",
			accent: "bg-blue-50",
			light: "bg-blue-50",
		},
	},

	// Spacing System
	spacing: {
		xs: "gap-1 p-1", // 4px
		sm: "gap-2 p-2", // 8px
		md: "gap-4 p-4", // 16px
		lg: "gap-6 p-6", // 24px
		xl: "gap-8 p-8", // 32px
		// Compact spacing for headers - 35% reduction from md
		compact: "gap-3 p-3", // ~10px (reduced from 16px)
	},

	// Border Radius
	radius: {
		none: "rounded-none",
		sm: "rounded-sm",
		md: "rounded-md",
		lg: "rounded-lg",
		full: "rounded-full",
	},

	// Shadow System
	shadows: {
		none: "shadow-none",
		sm: "shadow-sm",
		md: "shadow-md",
		lg: "shadow-lg",
		xl: "shadow-xl",
	},

	// Typography Scale
	typography: {
		xs: "text-xs",
		sm: "text-sm",
		base: "text-base",
		lg: "text-lg",
		xl: "text-xl",
		"2xl": "text-2xl",
	},

	// Component Variants
	components: {
		card: {
			base: "bg-white shadow-sm p-4 rounded-lg",
			interactive:
				"bg-white shadow-sm p-4 rounded-lg hover:shadow-md transition-shadow duration-200 relative z-10",
			bordered: "bg-white border border-gray-200 p-4 rounded-lg",
			subscription:
				"bg-white shadow-sm overflow-hidden hover:shadow-lg transition-all duration-300 relative z-10 border border-blue-100 hover:border-blue-200 transform hover:-translate-y-1",
		},
		button: {
			primary:
				"bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors",
			secondary:
				"bg-gray-100 text-gray-700 px-4 py-2 rounded-lg font-medium hover:bg-gray-200 transition-colors",
		},
		avatar: {
			xs: "w-6 h-6 rounded-full overflow-hidden bg-gray-100", // Very small
			sm: "w-8 h-8 rounded-full overflow-hidden bg-gray-100",
			md: "w-12 h-12 rounded-full overflow-hidden bg-gray-100",
			compact: "w-10 h-10 rounded-full overflow-hidden bg-gray-100", // Between sm and md
			lg: "w-16 h-16 rounded-full overflow-hidden bg-gray-100",
		},
	},
} as const;

// Helper functions for consistent styling
export const cn = (...classes: string[]) => classes.filter(Boolean).join(" ");

// Common layout patterns
export const layouts = {
	container: "w-full max-w-4xl mx-auto",
	containerFullWidth: "w-full", // Full width with no centering margins
	containerNoPadding: "w-full max-w-4xl mx-auto px-0", // Centered but no horizontal padding
	flexCenter: "flex items-center justify-center",
	flexBetween: "flex items-center justify-between",
	flexCol: "flex flex-col",
	grid2: "grid sm:grid-cols-2 gap-4",
	stickyHeader: "fixed top-0 left-0 right-0 z-50",
	pageContent: "pt-20 pb-12 space-y-4", // Reduced from pt-24 (96px to ~80px)
	contentWithGap: "pt-24 pb-12 space-y-4", // 112px for extra breathing room if needed
} as const;
