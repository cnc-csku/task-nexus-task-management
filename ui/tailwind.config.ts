import { heroui } from "@heroui/theme";
import type { Config } from "tailwindcss";

export default {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./node_modules/@heroui/theme/dist/components/(accordion|alert|autocomplete|avatar|badge|breadcrumbs|button|calendar|card|checkbox|chip|code|date-input|date-picker|divider|dropdown|form|image|input|kbd|link|listbox|menu|modal|navbar|pagination|popover|progress|radio|ripple|scroll-shadow|select|skeleton|slider|snippet|spacer|spinner|toggle|table|tabs|user).js",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#333333", 
          100: "#E6E6E6", // Very light gray
          200: "#CCCCCC", // Light gray
          300: "#B3B3B3", // Softer gray
          400: "#999999", // Medium gray
          500: "#808080", // Neutral gray
          600: "#666666", // Darker gray
          700: "#4D4D4D", // Dark gray
          800: "#333333", // Very dark gray
          900: "#1A1A1A", // Minimal black
        },
      },
    },
  },
  plugins: [heroui()],
} satisfies Config;
