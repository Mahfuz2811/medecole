import type { NextConfig } from "next";

const nextConfig: NextConfig = {
	reactStrictMode: false, // Disable strict mode to prevent double API calls in development
	output: "standalone", // Enable standalone output for Docker
	images: {
		remotePatterns: [
			{
				protocol: "https",
				hostname: "raw.githubusercontent.com",
				port: "",
				pathname: "/**",
			},
			{
				protocol: "http",
				hostname: "localhost",
				port: "8080",
				pathname: "/**",
			},
			{
				protocol: "https",
				hostname: "**", // Allow all HTTPS domains for flexibility
			},
		],
		unoptimized: false, // Enable optimization but allow fallback
		dangerouslyAllowSVG: true,
		contentDispositionType: "attachment",
		contentSecurityPolicy:
			"default-src 'self'; script-src 'none'; sandbox;",
	},
};

export default nextConfig;
