import { Page, expect } from "@playwright/test";

export class AdminPage {
  constructor(private readonly page: Page) {}

  async login(username = "demo", password = "demo") {
    await this.page.goto("https://localhost:4000/admin/");

    const profileInfo = this.page.getByRole("button", {
      name: `Logged in as ${username}`,
    });
    await expect(profileInfo).not.toBeVisible();

    await this.page.getByRole("textbox", { name: "Username" }).fill(username);
    await this.page.getByRole("textbox", { name: "Password" }).fill(password);
    await this.page.getByRole("button", { name: "Sign in" }).click();
  }

  async loginSuccessfully(username = "demo", password = "demo") {
    await this.login(username, password);

    await expect(this.page.getByText("Product Sales Stats")).toBeVisible();
    const profileInfo = this.page.getByRole("button", {
      name: `Profile`,
    });
    await expect(profileInfo.getByText(username)).toBeVisible();
  }
}
