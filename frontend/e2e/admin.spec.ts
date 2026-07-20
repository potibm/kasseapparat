import { test } from "@playwright/test";

test.describe("Admin", () => {
  /*
  test("should allow a user to log in with valid credentials", async ({
    page,
  }) => {
    const loginPage = new AdminPage(page);
    await loginPage.loginSuccessfully("demo", "demo");
  });

  test("should not log in with invalid credentials", async ({ page }) => {
    const loginPage = new AdminPage(page);
    await loginPage.login("invalidUser", "invalidPass");

    const profileInfo = page.getByRole("button", {
      name: `Role`,
    });
    await expect(profileInfo).not.toBeVisible();
    await expect(
      page.getByText("incorrect Username or Password"),
    ).toBeVisible();
  });
  */
});
