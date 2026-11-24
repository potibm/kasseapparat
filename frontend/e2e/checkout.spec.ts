import { test, expect } from "@playwright/test";
import { PosPage } from "./pages/PosPage";
import { resetDatabase } from "./helpers/db";
import { TEST_PRODUCTS } from "./fixtures/products";

test.describe("checkout", () => {
  const regularTicketProduct = TEST_PRODUCTS.REGULAR_TICKET;
  const coffeeMugProduct = TEST_PRODUCTS.COFFEE_MUG;

  test.beforeEach(async ({ page }) => {
    await resetDatabase();

    await page.goto("/");

    const pos = new PosPage(page);
    await pos.expectProductVisible(regularTicketProduct);
  });

  test("should have an empty cart on start", async ({ page }) => {
    const pos = new PosPage(page);
    await pos.expectEmptyCart();
  });

  test("should add a shirt size from the size matrix", async ({ page }) => {
    const pos = new PosPage(page);
    const shirt = TEST_PRODUCTS.MALE_SHIRT_L;

    const maleRow = page.getByRole("group", { name: "Male", exact: true });
    await expect(maleRow).toBeVisible();

    await maleRow.getByTestId(`product-card-${shirt.id}`).click();

    await pos.expectProductInCart(shirt);
    await expect(pos.cartItems).toHaveCount(1);
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

  test("should show the change for cash", async ({ page }) => {
    const pos = new PosPage(page);

    await pos.addProduct(regularTicketProduct);
    await pos.expectTotal("40 €");

    await page.getByTestId("checkout-button-CASH").click();
    await page.getByRole("button", { name: "50 €", exact: true }).click();

    await expect(page.getByTestId("change-calculator-change")).toHaveText(
      "10 €",
    );

    await page.getByTestId("change-calculator-confirm").click();

    await pos.expectEmptyCart();
  });
});
