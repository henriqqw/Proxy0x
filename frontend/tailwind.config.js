/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        surface: {
          0: '#191919',
          1: '#202020',
          2: '#252525',
          3: '#2f2f2f',
          4: '#373737',
          5: '#434343',
          6: '#505050',
        },
        border: {
          subtle: 'rgba(255,255,255,0.06)',
          default: 'rgba(255,255,255,0.09)',
          strong: 'rgba(255,255,255,0.14)',
        },
        text: {
          primary: '#ebebeb',
          secondary: '#999999',
          tertiary: '#6b6b6b',
        },
        accent: {
          DEFAULT: '#ffffff',
          muted: '#aaaaaa',
        },
        status: {
          online: '#4ac083',
          offline: '#6b6b6b',
          connecting: '#d4a024',
          error: '#cc4444',
        },
      },
      fontFamily: {
        sans: ['-apple-system', 'BlinkMacSystemFont', 'SF Pro Display', 'Segoe UI', 'sans-serif'],
        mono: ['SF Mono', 'Cascadia Code', 'Consolas', 'monospace'],
      },
      borderRadius: {
        DEFAULT: '8px',
        lg: '12px',
        xl: '16px',
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(8px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
