import { test, expect } from "@playwright/test";
import { LoginPage } from "./pages/LoginPage";
import { PosPage } from "./pages/PosPage";
import { resetDatabase } from "./helpers/db";
import { TEST_PRODUCTS } from "./fixtures/products";

test.describe("checkout", () => {
  const regularTicketProduct = TEST_PRODUCTS.REGULAR_TICKET;
  const coffeeMugProduct = TEST_PRODUCTS.COFFEE_MUG;

  test.beforeEach(async ({ page }) => {
    resetDatabase();

    const loginPage = new LoginPage(page);
    await loginPage.loginSuccessfully("demo", "demo");
  });

  test("should have an empty cart on start", async ({ page }) => {
    const pos = new PosPage(page);
    await pos.expectEmptyCart();
  });

  test("should add a product to the cart and remove it", async ({ page }) => {
    const pos = new PosPage(page);

    await pos.expectEmptyCart();
    await pos.addProduct(regularTicketProduct);

    await pos.expectProductInCart(regularTicketProduct);
    await expect(pos.cartItems).toHaveCount(1);

    await pos.removeProductFromCart(regularTicketProduct.name);

    await pos.expectEmptyCart();
  });

  test("should add multiple products to the cart and remove all", async ({
    page,
  }) => {
    const pos = new PosPage(page);

    await pos.expectEmptyCart();
    await pos.addProduct(regularTicketProduct);
    await pos.addProduct(coffeeMugProduct);

    await expect(pos.cartItems).toHaveCount(2);

    await pos.removeAllProductsFromCart();

    await pos.expectEmptyCart();
  });

  test("should add products to the cart and checkout", async ({ page }) => {
    const pos = new PosPage(page);

    await pos.expectEmptyCart();

    await pos.addProduct(regularTicketProduct);
    await pos.addProduct(coffeeMugProduct);

    await pos.expectProductInCart(regularTicketProduct);
    await pos.expectProductInCart(coffeeMugProduct);
    await expect(pos.cartItems).toHaveCount(2);
    await pos.expectTotal("41 €");

    await pos.checkout("CASH");

    await pos.expectEmptyCart();
  });
});
