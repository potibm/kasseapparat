import { expect, Locator, Page } from "@playwright/test";

const AUTH_RACE_ERROR = /checkAuth must be called first/;

export class AdminPage {
  readonly page: Page;
  readonly dashboardHeading: Locator;
  readonly listTable: Locator;
  readonly interopError: Locator;
  readonly authRaceErrors: string[] = [];

  constructor(page: Page) {
    this.page = page;

    this.dashboardHeading = page.getByRole("heading", { name: "Dashboard" });
    this.listTable = page.getByRole("table");
    this.interopError = page.getByText(/decodeComponent is not a function/);

    // Registered before any navigation so a failure during the first render is
    // still observed. react-admin calls getPermissions() during the first
    // render, which used to lose the race against the in-flight /auth/me
    // request and log this error.
    page.on("console", (message) => {
      if (AUTH_RACE_ERROR.test(message.text())) {
        this.authRaceErrors.push(message.text());
      }
    });
    page.on("pageerror", (error) => {
      if (AUTH_RACE_ERROR.test(error.message)) {
        this.authRaceErrors.push(error.message);
      }
    });
  }

  async goto() {
    await this.page.goto("/admin/");

    await this.page.waitForLoadState("networkidle");
  }

  async gotoPurchases() {
    await this.page.goto("/admin/purchases");

    await this.page.waitForLoadState("networkidle");
  }

  async gotoProducts() {
    await this.page.goto("/admin/products");

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
    await expect(this.listTable).toBeVisible();
  }

  /**
   * The products list is the only admin view that asks for permissions while
   * it is still mounting, so it is the only one that used to stall for a
   * second on a cold load while react-query retried the failed lookup.
   */
  async expectProductsListIsVisible() {
    await expect(this.listTable).toBeVisible();
  }

  async expectNoAuthRaceErrors() {
    expect(this.authRaceErrors).toEqual([]);
  }
}
