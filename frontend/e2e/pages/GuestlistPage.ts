import { Page, Locator, expect } from "@playwright/test";

export class GuestlistPage {
  readonly guestItems: Locator;
  readonly guestListTable: Locator;
  readonly searchInput: Locator;

  constructor(page: Page) {
    this.guestItems = page.getByTestId(/^guestlist-result-id/);
    this.guestListTable = page.getByTestId("guestlist-result-table");
    this.searchInput = page.getByTestId("guestlist-search-input");
  }

  async expectCountToBe(count: number) {
    await expect(this.guestItems).toHaveCount(count);
  }

  async expectListToBeEmpty() {
    await this.expectCountToBe(0);
  }

  async expectListNotToBeEmpty() {
    await expect(this.guestItems.first()).toBeVisible();
  }

  async expectListToContain(name: string) {
    await expect(this.guestListTable.getByText(name)).toBeVisible();
  }

  async getSearchTerm() {
    return this.searchInput.inputValue();
  }
}
