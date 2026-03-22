import { test, expect } from "@playwright/test";
import { LoginPage } from "./pages/LoginPage";
import { PosPage } from "./pages/PosPage";
import { resetDatabase } from "./helpers/db";
import { TEST_PRODUCTS } from "./fixtures/products";
import { GuestlistPage } from "./pages/GuestlistPage";

test.describe("guestlist", () => {
  const freeTicketProduct = TEST_PRODUCTS.FREE_TICKET;

  test.describe("free product", () => {
    test.beforeEach(async ({ page }) => {
      resetDatabase();

      const loginPage = new LoginPage(page);
      await loginPage.loginSuccessfully("demo", "demo");

      const pos = new PosPage(page);
      await pos.openGuestlistModal(freeTicketProduct);
    });

    test("should open the guest list modal", async ({ page }) => {
      const guestlistPage = new GuestlistPage(page);

      await guestlistPage.expectListNotToBeEmpty();
    });

    test("should filter the guest list based on search input", async ({
      page,
    }) => {
      const guestlistPage = new GuestlistPage(page);

      await page.getByTestId("guestlist-search-input").fill("Jean");

      await guestlistPage.expectListNotToBeEmpty();
      await guestlistPage.expectListToContain("Jean Dupont");
    });
  });

  test.describe("using the onscreen keyboard", () => {
    test.beforeEach(async ({ page }) => {
      resetDatabase();

      const loginPage = new LoginPage(page);
      await loginPage.loginSuccessfully("demo", "demo");

      const pos = new PosPage(page);
      await pos.openGuestlistModal(freeTicketProduct);
    });

    test("should allow searching using the onscreen keyboard", async ({
      page,
    }) => {
      const guestlistPage = new GuestlistPage(page);
      expect(await guestlistPage.getSearchTerm()).toEqual("");

      await page.getByRole("button", { name: "Add J to search term" }).click();
      await page.getByRole("button", { name: "Add E to search term" }).click();
      await page.getByRole("button", { name: "Add A to search term" }).click();
      await page.getByRole("button", { name: "Add N to search term" }).click();
      expect(await guestlistPage.getSearchTerm()).toEqual("JEAN");

      await page
        .getByRole("button", { name: "Remove last character from search term" })
        .click();
      expect(await guestlistPage.getSearchTerm()).toEqual("JEA");

      await page.getByRole("button", { name: "Clear search term" }).click();
      expect(await guestlistPage.getSearchTerm()).toEqual("");
    });
  });
});
