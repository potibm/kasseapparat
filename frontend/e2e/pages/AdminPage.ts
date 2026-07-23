import { expect, Locator, Page } from "@playwright/test";

export class AdminPage {
  readonly page: Page;
  readonly dashboardHeading: Locator;

  constructor(page: Page) {
    this.page = page;

    this.dashboardHeading = page.getByRole("heading", { name: "Dashboard" });
  }

  async goto() {
    await this.page.goto("/admin/");

    await this.page.waitForLoadState("networkidle");
  }

  async expectDashboardIsVisible() {
    await expect(this.dashboardHeading).toBeVisible();
  }
}
