import { nextui } from "@nextui-org/theme";
import type { Config } from "tailwindcss";

export default {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./node_modules/@nextui-org/theme/dist/components/(accordion|alert|autocomplete|avatar|badge|breadcrumbs|button|calendar|card|checkbox|chip|code|date-input|date-picker|divider|dropdown|form|image|input|kbd|link|listbox|menu|modal|navbar|pagination|popover|progress|radio|ripple|scroll-shadow|select|skeleton|slider|snippet|spacer|spinner|toggle|table|tabs|user).js",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--background)",
        foreground: "var(--foreground)",
      },
    },
  },
  plugins: [nextui()],
} satisfies Config;
