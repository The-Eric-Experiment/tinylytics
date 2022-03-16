import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 8006,
    host: "0.0.0.0",
    proxy: {
      "/api": {
        target: "http://192.168.1.60:8080/",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
