import { test, expect } from "@playwright/test";

test.describe("System", () => {
  test("should start the application with the login page", async ({ page }) => {
    await page.goto("https://localhost:4000/");

    await expect(page).toHaveTitle(/Kasseapparat/);

    const usernameInput = page.getByRole("textbox", { name: "Your username" });
    await expect(usernameInput).toBeVisible();
  });
});
