/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'poo-brown': {
          50: '#F7F3F0',
          100: '#E8DDD6',
          200: '#D0BBA6',
          300: '#B89976',
          400: '#A07746',
          500: '#8B5A2B',
          600: '#6F4622',
          700: '#533318',
          800: '#37200F',
          900: '#1C0E05'
        }
      },
      fontFamily: {
        'comic': ['Comic Sans MS', 'cursive']
      }
    },
  },
  plugins: [],
}
