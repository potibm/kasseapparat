import { expect, test } from "@playwright/test";
import { AdminPage } from "./pages/AdminPage";

test.describe("Admin", () => {
  test("should see the dashboard", async ({ page }) => {
    const dashboard = new AdminPage(page);

    await dashboard.goto();
    await dashboard.expectDashboardIsVisible();
  });
});
