import { test, expect } from "@playwright/test";
import { LoginPage } from "./pages/LoginPage";

test.describe("Authentication", () => {
  test("should allow a user to log in with valid credentials", async ({
    page,
  }) => {
    const loginPage = new LoginPage(page);
    await loginPage.loginSuccessfully("demo", "demo");
  });

  test("should not log in with invalid credentials", async ({ page }) => {
    const loginPage = new LoginPage(page);
    await loginPage.login("invalidUser", "invalidPass");

    const profileInfo = page.getByRole("button", {
      name: `Logged in as invalidUser`,
    });
    await expect(profileInfo).not.toBeVisible();
  });
});
