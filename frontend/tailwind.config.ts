import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx}", // Make sure Tailwind scans your app files
  ],
  theme: {
    extend: {},
    container: {
      center: true,
      padding: "1rem", // equals px-4 (16px)
      screens: {
        sm: "600px",
        md: "768px",
        lg: "1024px",
        xl: "1280px",
      },
    },
  },
  plugins: [],
};

export default config;