import { Page, expect } from "@playwright/test";

export class LoginPage {
  constructor(private readonly page: Page) {}

  async login(username = "demo", password = "demo") {
    await this.page.goto("https://localhost:4000/");

    const profileInfo = this.page.getByRole("button", {
      name: `Logged in as ${username}`,
    });
    await expect(profileInfo).not.toBeVisible();

    await this.page
      .getByRole("textbox", { name: "Your username" })
      .fill(username);
    await this.page
      .getByRole("textbox", { name: "Your password" })
      .fill(password);
    await this.page.getByRole("button", { name: "Login" }).click();
  }

  async loginSuccessfully(username = "demo", password = "demo") {
    await this.login(username, password);

    const profileInfo = this.page.getByRole("button", {
      name: `Logged in as ${username}`,
    });
    await expect(profileInfo).toBeVisible();
  }
}
