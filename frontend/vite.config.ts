/// <reference types="vitest/config" />
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import flowbiteReact from "flowbite-react/plugin/vite";
import basicSsl from "@vitejs/plugin-basic-ssl";
import path from "node:path";

const __dirname = path.resolve();

const frontendPort = process.env.E2E_PORT
  ? Number.parseInt(process.env.E2E_PORT)
  : 3000;
const backendTarget = process.env.E2E_API_TARGET || "http://127.0.0.1:3001";

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
        configure: (proxy, _options) => {
          proxy.on("error", (err, _req, _res) => {
            //eslint-disable-next-line no-console
            console.log("proxy error", err);
          });
          proxy.on("proxyReq", (proxyReq, req, _res) => {
            //eslint-disable-next-line no-console
            console.log(
              "Vite proxy forwards this request:",
              req.method,
              req.url,
            );
          });
        },
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
