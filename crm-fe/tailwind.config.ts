/** @type {import('tailwindcss').Config} */
import type { Config } from 'tailwindcss'

export default <Partial<Config>> {
  
  content: [],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'], // Example font family
        serif: ['Georgia', 'serif'],
        mono: ['Menlo', 'monospace'],
      },
      aspectRatio: {
        auto: 'auto',
        square: '1 / 1',
        video: '16 / 9'
      }
    },
  },
  plugins: [],
}


