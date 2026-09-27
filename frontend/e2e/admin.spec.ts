import { test } from "@playwright/test";
import { AdminPage } from "./pages/AdminPage";

test.describe("Admin", () => {
  test("should see the dashboard", async ({ page }) => {
    const dashboard = new AdminPage(page);

    await dashboard.goto();
    await dashboard.expectDashboardIsVisible();
  });

  test("should see the purchases list", async ({ page }) => {
    const dashboard = new AdminPage(page);

    await dashboard.gotoPurchases();
    await dashboard.expectPurchasesListIsVisible();
  });

  test("should keep the purchases list working after sorting", async ({
    page,
  }) => {
    const dashboard = new AdminPage(page);

    await dashboard.gotoPurchases();
    await dashboard.expectPurchasesListIsVisible();

    await dashboard.sortPurchasesById();
    await dashboard.expectPurchasesListIsVisible();
  });
});
