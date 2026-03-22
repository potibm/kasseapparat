/// <reference types="vitest/config" />
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import flowbiteReact from "flowbite-react/plugin/vite";
import basicSsl from "@vitejs/plugin-basic-ssl";
import path from "node:path";

const __dirname = path.resolve();

const frontendPort = process.env.E2E_PORT
  ? parseInt(process.env.E2E_PORT)
  : 3000;
const backendTarget = process.env.E2E_API_TARGET || "http://localhost:3001";

export default defineConfig({
  plugins: [react(), tailwindcss(), flowbiteReact(), basicSsl()],
  server: {
    port: frontendPort,
    strictPort: true,
    proxy: {
      "/api": {
        target: backendTarget,
        changeOrigin: true,
        ws: true,
        secure: false,
      },
    },
  },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: "./tests/setup.ts",
    teardownTimeout: 1000,
    pool: "threads",
    include: ["src/**/*.{test,spec}.{ts,mts,cts,jsx,tsx}"],
    coverage: {
      provider: "v8",
      reporter: ["text", "html", "lcov"],
    },
  },
  resolve: {
    alias: {
      "@core": path.resolve(__dirname, "./src/core"),
      "@admin": path.resolve(__dirname, "./src/apps/admin"),
      "@pos": path.resolve(__dirname, "./src/apps/pos"),
    },
  },
});
