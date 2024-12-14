import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: "localhost",
    port: 3034,
    proxy: {
      "/api": {
        target: "http://localhost:3033", // Backend API Path
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
