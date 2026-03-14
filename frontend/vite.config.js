import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import flowbiteReact from "flowbite-react/plugin/vite";
import basicSsl from "@vitejs/plugin-basic-ssl";
import path from "path";
const __dirname = import.meta.dirname;

// https://vite.dev/config/
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
