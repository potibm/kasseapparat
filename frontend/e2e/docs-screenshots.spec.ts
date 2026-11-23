import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect, Page } from "@playwright/test";
import { PosPage } from "./pages/PosPage";
import { resetDatabase } from "./helpers/db";
import { TEST_PRODUCTS } from "./fixtures/products";

/**
 * Screenshots for docs/manual.md, taken on the iPad the POS is made for. They are
 * written to frontend/docs-screenshots/; `mise run docs:screenshots` takes them and
 * copies them to docs/, and CI uploads them as an artifact.
 */
const OUTPUT_DIR = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "../docs-screenshots",
);

const IPAD_PROJECT = "iPad air (landscape)";

const shoot = async (page: Page, name: string) => {
  // Fonts must be loaded, otherwise the pictures differ from run to run.
  await page.evaluate(() => document.fonts.ready);
  await page.screenshot({
    path: path.join(OUTPUT_DIR, name),
    animations: "disabled",
    caret: "hide",
  });
};

test.describe("documentation screenshots", () => {
  test.beforeEach(async ({ page }, testInfo) => {
    test.skip(testInfo.project.name !== IPAD_PROJECT, "iPad screenshots only");

    await resetDatabase();
    await page.goto("/");
    await new PosPage(page).expectProductVisible(TEST_PRODUCTS.REGULAR_TICKET);
  });

  test("POS with a filled cart", async ({ page }) => {
    const pos = new PosPage(page);

    await pos.addProduct(TEST_PRODUCTS.REGULAR_TICKET);
    await page
      .getByRole("group", { name: "Male", exact: true })
      .getByTestId(`product-card-${TEST_PRODUCTS.MALE_SHIRT_L.id}`)
      .click();
    await pos.addProduct(TEST_PRODUCTS.COFFEE_MUG);
    await expect(pos.cartItems).toHaveCount(3);

    await shoot(page, "pos.png");
  });

  test("guest list with a search", async ({ page }) => {
    const pos = new PosPage(page);

    await pos.openGuestlistModal(TEST_PRODUCTS.FREE_TICKET);
    await page.getByRole("button", { name: "Add J to search term" }).click();
    await page.getByRole("button", { name: "Add E to search term" }).click();
    await expect(page.getByTestId("guestlist-result-table")).toBeVisible();

    await shoot(page, "pos-list.png");
  });

  test("change calculator for cash", async ({ page }) => {
    const pos = new PosPage(page);

    await pos.addProduct(TEST_PRODUCTS.REGULAR_TICKET);
    await page.getByTestId("checkout-button-CASH").click();
    await page.getByRole("button", { name: "50 €", exact: true }).click();
    await expect(page.getByTestId("change-calculator-change")).toBeVisible();

    await shoot(page, "pos-cash.png");
  });

  test("admin dashboard", async ({ page }) => {
    await page.goto("/admin");
    await expect(page.getByTestId("dashboard-kpis")).toBeVisible();

    await shoot(page, "admin.png");
  });
});
