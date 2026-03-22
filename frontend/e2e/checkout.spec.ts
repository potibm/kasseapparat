import { test, expect } from "@playwright/test";
import { LoginPage } from "./pages/LoginPage";
import { resetDatabase } from "./helpers/db";

test.describe("Kassiervorgänge", () => {
  test.beforeEach(async ({ page }) => {
    resetDatabase();

    const loginPage = new LoginPage(page);
    await loginPage.loginSuccessfully("demo", "demo");
  });

  test("should perform a standard checkout", async ({ page }) => {
    // check that cart is empty
    const cartItems = page.getByTestId("cart-items");
    await expect(cartItems).toBeEmpty();

    // check that purchase history is empty
    const purchaseHistory = page.getByTestId("purchase-history");
    await expect(purchaseHistory).toBeEmpty();

    // add something to the cart
    await page
      .getByRole("button", { name: "Add 🎟️ Regular for 40 € to" })
      .click();
    await page
      .getByRole("button", { name: "Add ☕ Coffee Mug for 1 € to" })
      .click();

    // ckech that cart has the items
    await expect(cartItems).toContainText("Regular");
    await expect(cartItems).toContainText("Coffee Mug");

    // checkout
    await page.getByTestId("checkout-button-CASH").click();

    // check if purchase history has the new entry
    await expect(purchaseHistory).not.toBeEmpty();
  });
});
