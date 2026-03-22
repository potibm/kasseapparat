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

  async addProductByName(name: string) {
    await this.page.getByRole("button", { name }).click();
  }

  async addProduct(product: Product) {
    await this.addProductByName(
      `Add ${product.name} for ${product.price} to cart`,
    );
  }

  async openGuestlistModalByName(productName: string) {
    await this.page
      .getByRole("button", { name: `Show guestlist for ${productName}` })
      .click();
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
