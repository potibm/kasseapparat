import { defineConfig, devices } from "@playwright/test";

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false,
  workers: 1,
  reporter: "html",
  use: {
    baseURL: "http://localhost:4000",
    trace: "retain-on-failure",
    screenshot: "only-on-failure",

    locale: "da-DK",
    timezoneId: "Europe/Berlin",

    extraHTTPHeaders: {
      "X-Remote-User": "e2e-demo-user",
    },
  },

  projects: [
    { name: "chromium", use: { ...devices["Desktop Chrome"] } },
    {
      name: "iPad air (landscape)",
      use: {
        ...devices["iPad (gen 11) landscape"],
      },
    },
  ],

  webServer: [
    {
      command:
        'cd ../backend && go run . serve --port=4001 --log-level=debug --db-file "e2e-work"',
      port: 4001,
      reuseExistingServer: !process.env.CI,
      stdout: "pipe",
      env: {
        APP_CORS_ALLOW_ORIGINS: "http://localhost:4000",
        APP_REDIS_URL: "",
        FORMAT_CURRENCY_LOCALE: "de-DE",
        FORMAT_CURRENCY_CODE: "EUR",
        FORMAT_CURRENCY_FRACTION_DIGITS_MAX: "2",
        FORMAT_CURRENCY_FRACTION_DIGITS_MIN: "0",
        FORMAT_DATE_LOCALE: "de-DE",
        AUTH_MODE: "proxy",
      },
    },
    {
      command: "E2E_PORT=4000 E2E_API_TARGET=http://localhost:4001 npm run dev",
      port: 4000,
      reuseExistingServer: !process.env.CI,
      ignoreHTTPSErrors: true,
      stdout: "pipe",
    },
  ],
});
