import { test, expect } from "@playwright/test";

test.describe("System", () => {
  test("should start the application and show the authenticated user", async ({
    page,
  }) => {
    await page.goto("/pos");

    await expect(page).toHaveTitle(/Kasseapparat/);

    const userNameButton = page.getByRole("button", { name: /e2e-demo-user/i });
    await expect(userNameButton).toBeVisible();
  });

  test("offers the POS and the admin on the start page", async ({ page }) => {
    await page.goto("/");

    await expect(page.getByRole("link", { name: /Admin/ })).toHaveAttribute(
      "href",
      "/admin",
    );

    await page.getByRole("link", { name: /POS/ }).click();

    await expect(page).toHaveURL(/\/pos$/);
    await expect(
      page.getByRole("button", { name: /e2e-demo-user/i }),
    ).toBeVisible();
  });
});
