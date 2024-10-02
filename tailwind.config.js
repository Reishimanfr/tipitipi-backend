/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html","./src/frontend/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      animation: {
        'slide-in': 'slide-in 0.5s ease-in-out',
        'slide-out': 'slide-out 0.5s ease-in-out',
        'slideInLeft': 'slideInFromLeft 1s ease-out',
        'slideInRight': 'slideInFromRight 1s ease-out'
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
            transform: 'translateX(100%)',
            //opacity: 0,
          },
        },
        slideInFromLeft: {
          '0%': { transform: 'translateX(-10%)', opacity: '0' },
          '100%': { transform: 'translateX(0)', opacity: '1' },
        },
        slideInFromRight: {
          '0%': { transform: 'translateX(10%)', opacity: '0' },
          '100%': { transform: 'translateX(0)', opacity: '1' },
        }
      },
    },
  },
  plugins: [
    function ({ addUtilities }) {
      const newUtilities = {
        '.fill-forwards': {
          'animation-fill-mode': 'forwards',
        },
        // Dodaj więcej niestandardowych właściwości CSS, jeśli to konieczne
      };
      
      addUtilities(newUtilities, ['responsive', 'hover']);
    },
    require('tailwindcss-animated')
  ],
};