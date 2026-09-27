// tests/e2e/pages/PosPage.ts
import { Page, Locator, expect } from "@playwright/test";
import { Product } from "../fixtures/products";

export class PosPage {
  readonly page: Page;
  readonly cartItems: Locator;
  readonly cartTable: Locator;
  readonly checkoutCashButton: Locator;

  constructor(page: Page) {
    this.page = page;
    this.cartItems = page.getByTestId(/^cart-product-/);
    this.cartTable = page.getByTestId("cart-table");
    this.checkoutCashButton = page.getByTestId("checkout-button-CASH");
  }

  /**
   * The product grid is rendered from a client-side fetch, so `page.goto`
   * resolving does not mean the POS is interactive yet. Gate on a rendered card
   * so a click is never issued against a half-loaded page.
   */
  async expectProductVisible(product: Product) {
    await expect(
      this.page.getByTestId(`product-card-${product.id}`),
    ).toBeVisible();
  }

  async addProductByName(name: string) {
    await this.page.getByRole("button", { name }).click();
  }

  async addProduct(product: Product) {
    await this.addProductByName(
      `Add ${product.name} for ${product.price} to cart`,
    );
  }

  async openGuestlistModalByName(productName: string) {
    const button = this.page.getByRole("button", {
      name: `Show guestlist for ${productName}`,
    });

    // Assert before clicking so a missing button fails here with a readable
    // reason instead of timing out inside the actionability check.
    await expect(button).toBeVisible();
    await button.click();

    await expect(this.page.getByTestId("guestlist-search-input")).toBeVisible();
  }

  async openGuestlistModal(product: Product) {
    await this.openGuestlistModalByName(product.name);
  }

  async checkout(method: "CASH" | "CC") {
    await this.page.getByTestId(`checkout-button-${method}`).click();
  }

  async expectEmptyCart() {
    await expect(this.cartItems).toHaveCount(0);
  }

  async expectProductInCartByProductId(productId: string | number) {
    const item = this.page.getByTestId(`cart-product-${productId}`);
    await expect(item).toBeVisible();
  }

  async expectProductInCart(product: Product) {
    return this.expectProductInCartByProductId(product.id);
  }

  async expectTotal(amount: string) {
    await expect(this.checkoutCashButton).toBeEnabled();
    await expect(this.checkoutCashButton).toContainText(amount);
  }

  async removeProductFromCart(name: string) {
    await this.page
      .getByRole("button", { name: `Remove ${name} from cart` })
      .click();
  }

  async removeAllProductsFromCart() {
    await this.page
      .getByRole("button", { name: `Remove all items from cart` })
      .click();
  }
}
