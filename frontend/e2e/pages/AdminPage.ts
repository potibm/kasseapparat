import { expect, Locator, Page } from "@playwright/test";

export class AdminPage {
  readonly page: Page;
  readonly dashboardHeading: Locator;
  readonly purchasesTable: Locator;
  readonly interopError: Locator;

  constructor(page: Page) {
    this.page = page;

    this.dashboardHeading = page.getByRole("heading", { name: "Dashboard" });
    this.purchasesTable = page.getByRole("table");
    this.interopError = page.getByText(/decodeComponent is not a function/);
  }

  async goto() {
    await this.page.goto("/admin/");

    await this.page.waitForLoadState("networkidle");
  }

  async gotoPurchases() {
    await this.page.goto("/admin/purchases");

    await this.page.waitForLoadState("networkidle");
  }

  async expectDashboardIsVisible() {
    await expect(this.dashboardHeading).toBeVisible();
  }

  /**
   * Sorts the purchases list, which makes react-admin write `sort`/`order` into
   * the query string and re-parse it through ra-core's `parse()` from
   * `query-string`. A bare list URL carries no query string, so it never
   * reaches that code path -- the dashboard never does either, which is why a
   * dashboard-only test cannot catch a broken CJS/ESM interop in that chain.
   */
  async sortPurchasesById() {
    await this.page.getByRole("button", { name: /Sort by iD/ }).click();
  }

  async expectPurchasesListIsVisible() {
    await expect(this.interopError).toHaveCount(0);
    await expect(this.purchasesTable).toBeVisible();
  }
}
