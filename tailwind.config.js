/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html","./src/frontend/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        transparent: 'transparent',
        current: 'currentColor',
        'tipiOrange': '#EB9511',
        'tipiYellow': '#F0B93A',
        'tipiBrown': '#853A3F',
        'tipiPink': '#EC7B7F',
        'tipiWhite': '#F5F2E3'
      },
      fontFamily: {
        custom: ['Ubuntu', 'sans-serif'], // Nazwa "CustomFont" z CSS
      },

      animation: {
        'slide-in': 'slide-in 0.5s ease-in-out',
        'slide-out': 'slide-out 0.5s ease-in-out',
        // 'slideInLeft': 'slideInFromLeft 1s ease-out',
        // 'slideInRight': 'slideInFromRight 1s ease-out'
      },
      keyframes: {
        'slide-in': {
          '0%': {
            transform: 'translateX(100%)',
            //opacity: 0,
          },
          '100%': {
            transform: 'translateX(0)',
            //opacity: 1,
          },
        },
        'slide-out': {
          '0%': {
            transform: 'translateX(0)',
            //opacity: 1,
          },
          '100%': {
            transform: 'translateX(100%)'
            //opacity: 0,
          },
        },
      },
    },
  },
};