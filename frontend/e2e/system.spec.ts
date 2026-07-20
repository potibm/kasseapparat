import { test, expect } from "@playwright/test";

test.describe("System", () => {
  test("should start the application and show the authenticated user", async ({
    page,
  }) => {
    await page.goto("/");
    await page.pause();

    await expect(page).toHaveTitle(/Kasseapparat/);

    const userNameButton = page.getByRole("button", { name: /e2e-demo-user/i });
    await expect(userNameButton).toBeVisible();
  });
});
