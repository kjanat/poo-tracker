/**
 * Branding utilities and constants for Poo Tracker
 */

// Brand colors
export const BRAND_COLORS = {
  primary: "#8B4513",
  secondary: "#B35F37",
  accent: "#CA7047",
  text: "#6B4423",
  light: "#DF8963",
  dark: "#37161B",
} as const;

// Logo sizes for different contexts
export const LOGO_SIZES = {
  small: 16,
  medium: 32,
  large: 64,
  hero: 96,
  icon: 24,
} as const;

// Available PNG logo assets
export const LOGO_ASSETS = {
  favicon: "/favicon.ico",
  svg: "/logo.svg",
  png16: "/logo_16x16.png",
  png32: "/logo_32x32.png",
  png64: "/logo_64x64.png",
  png128: "/logo_128x128.png",
  png256: "/logo_256x256.png",
} as const;

// App metadata
export const APP_META = {
  name: "Poo Tracker - Your Bowel Movement Journal",
  shortName: "Poo Tracker",
  description:
    "Track your bowel movements, analyze patterns, and optimize your digestive health with Poo Tracker.",
  themeColor: BRAND_COLORS.primary,
  backgroundColor: "#ffffff",
} as const;

// Logo component props helper
export const getLogoProps = (context: "navbar" | "hero" | "login" | "icon") => {
  switch (context) {
    case "navbar":
      return { size: LOGO_SIZES.medium, className: "text-poo-brown-600" };
    case "hero":
      return { size: LOGO_SIZES.hero, className: "text-poo-brown-700" };
    case "login":
      return { size: LOGO_SIZES.large, className: "text-poo-brown-600" };
    case "icon":
      return { size: LOGO_SIZES.icon, className: "text-poo-brown-500" };
    default:
      return { size: LOGO_SIZES.medium, className: "text-poo-brown-600" };
  }
};
