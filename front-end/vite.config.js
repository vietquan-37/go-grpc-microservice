import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import dotenv from 'dotenv'

// Load environment variables from .env file
dotenv.config()

// https://vite.dev/config/
export default defineConfig({
    plugins: [react()],
    server: {
        port: 3000
    },
    define: {
        'process.env.VITE_BACK_END': JSON.stringify(process.env.VITE_BACK_END)
    }
})
