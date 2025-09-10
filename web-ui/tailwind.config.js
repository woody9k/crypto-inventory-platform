/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class', // Enable class-based dark mode switching
  theme: {
    extend: {
      // Custom color palette: Black, Gold, and Red theme
      colors: {
        // Primary color: Gold theme for main actions and branding
        primary: {
          50: '#fffbeb',   // Lightest gold
          100: '#fef3c7',  // Very light gold
          200: '#fde68a',  // Light gold
          300: '#fcd34d',  // Medium light gold
          400: '#fbbf24',  // Medium gold
          500: '#f59e0b',  // Main gold - primary brand color
          600: '#d97706',  // Dark gold - hover states
          700: '#b45309',  // Darker gold
          800: '#92400e',  // Very dark gold
          900: '#78350f',  // Darkest gold
        },
        // Secondary color: Black and gray theme for backgrounds and text
        secondary: {
          50: '#f8fafc',   // Lightest gray
          100: '#f1f5f9',  // Very light gray
          200: '#e2e8f0',  // Light gray
          300: '#cbd5e1',  // Medium light gray
          400: '#94a3b8',  // Medium gray
          500: '#64748b',  // Medium dark gray
          600: '#475569',  // Dark gray
          700: '#334155',  // Darker gray
          800: '#1e293b',  // Very dark gray
          900: '#0f172a',  // Deep black - primary background
        },
        // Accent color: Red theme for highlights, alerts, and important actions
        accent: {
          50: '#fef2f2',   // Lightest red
          100: '#fee2e2',  // Very light red
          200: '#fecaca',  // Light red
          300: '#fca5a5',  // Medium light red
          400: '#f87171',  // Medium red
          500: '#ef4444',  // Main red - primary accent color
          600: '#dc2626',  // Dark red - hover states
          700: '#b91c1c',  // Darker red
          800: '#991b1b',  // Very dark red
          900: '#7f1d1d',  // Darkest red
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
