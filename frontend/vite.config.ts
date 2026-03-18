/// <reference types="vitest/config" />
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import flowbiteReact from "flowbite-react/plugin/vite";
import basicSsl from "@vitejs/plugin-basic-ssl";
import path from "node:path";

const __dirname = path.resolve();

export default defineConfig({
  plugins: [react(), tailwindcss(), flowbiteReact(), basicSsl()],
  server: {
    port: 3000,
    proxy: {
      "/api": {
        target: "http://localhost:3001",
        changeOrigin: true,
        ws: true,
      },
    },
  },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: "./tests/setup.ts",
    teardownTimeout: 1000,
    pool: "threads",
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
