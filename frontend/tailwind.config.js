/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        dusk: {
          950: '#11100f',
          900: '#191817',
          800: '#252320',
          700: '#39342f',
          gold: '#c8a75d',
          moss: '#78866b',
          blood: '#8e3b46',
        },
      },
    },
  },
  plugins: [],
};
