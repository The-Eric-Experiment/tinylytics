import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  root: "./",
  build: {
    outDir: "dist",
  },
  publicDir: "public",
  plugins: [react()],
  server: {
    port: 8006,
    host: "0.0.0.0",
    proxy: {
      "/api": {
        // target: "http://192.168.1.60:8008/",
        target: "https://analytics.ericexperiment.com/",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
