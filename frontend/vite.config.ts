import { defineConfig } from 'vite'
import preact from '@preact/preset-vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [preact(),
    tailwindcss()
  ],
  base: './', // Ensure assets are loaded correctly
  server: {
    proxy: {
      '/api': {
        target: 'http://192.168.1.4:8081',
        // target: 'http://backend:8081',
        changeOrigin: true
      }
    }
  }
})