import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { resolve } from "path";

// const configs = [
//   { name: "dev", port: 5173 },
//   { name: "preview", port: 4173 }
// ];

// // ES2024 Object.groupBy if needed
// const configsByType = Object.groupBy(configs, c => 
//   c.port > 5000 ? 'dev' : 'prod'
// );

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:3002",
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: "dist",
    sourcemap: true,
    target: "es2024" // Match TypeScript target
  },
  resolve: {
    alias: {
      "@": resolve(__dirname, "./src"),
      "@/components": resolve(__dirname, "./src/components"),
      "@/hooks": resolve(__dirname, "./src/hooks"),
      "@/utils": resolve(__dirname, "./src/utils"),
      "@/types": resolve(__dirname, "./src/types"),
      "@/api": resolve(__dirname, "./src/api")
    }
  }
});
