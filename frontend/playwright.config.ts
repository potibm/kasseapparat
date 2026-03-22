import { defineConfig, devices } from "@playwright/test";

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// import dotenv from 'dotenv';
// import path from 'path';
// dotenv.config({ path: path.resolve(__dirname, '.env') });

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false,
  workers: 1,
  reporter: "html",
  use: {
    baseURL: "https://localhost:4000",
    ignoreHTTPSErrors: true,

    trace: "on-first-retry",
  },

  projects: [
    { name: "chromium", use: { ...devices["Desktop Chrome"] } },
    {
      name: "ipad air (landscape)",
      use: {
        ...devices["ipad (gen 11) landscape"],
      },
    },
  ],

  webServer: [
    {
      command:
        'cd ../backend && go run ./cmd/main.go --port=4001 --log-level=debug --db-file "e2e-work"',
      port: 4001,
      reuseExistingServer: !process.env.CI,
      stdout: "pipe",
    },
    {
      command:
        "E2E_PORT=4000 E2E_API_TARGET=http://localhost:4001 corepack yarn dev",
      port: 4000,
      reuseExistingServer: !process.env.CI,
      ignoreHTTPSErrors: true,
      stdout: "pipe",
    },
  ],
});
