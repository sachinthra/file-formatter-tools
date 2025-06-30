import { defineConfig } from 'vite'
import preact from '@preact/preset-vite'

export default defineConfig({
  plugins: [preact()],
  base: './', // Ensure assets are loaded correctly
  server: {
    proxy: {
      '/api': {
        target: 'http://backend:8081',
        changeOrigin: true
      }
    }
  }
})