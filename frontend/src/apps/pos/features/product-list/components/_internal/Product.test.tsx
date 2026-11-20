import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Product, { LOW_STOCK_THRESHOLD } from "./Product";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import { createMockProduct } from "@pos/api/schemas.mocks";
import { OutOfStockBehavior } from "@core/config/types/config.types";

vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: vi.fn(),
}));

// The guestlist modal fetches through usePosApi, which needs a config provider this
// test does not wrap. The card's own affordance is what is under test here.
vi.mock("@pos/api/usePosApi", () => ({
  usePosApi: () => ({
    fetchGuestlistByProductId: vi.fn().mockResolvedValue([]),
  }),
}));

const mockCurrency = { format: (val: number) => `${val} €` };

const setup = (outOfStockBehavior: OutOfStockBehavior) => {
  vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
    currency: mockCurrency,
    outOfStockBehavior,
  } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
};

type MockProduct = ReturnType<typeof createMockProduct>;

/** The label the card renders when the product can still be added. */
const addLabel = (product: MockProduct) =>
  "Add " +
  product.name +
  " for " +
  mockCurrency.format(product.grossPrice.toNumber()) +
  " to cart";

/** The label the card renders once the cart holds the remaining stock. */
const noStockLabel = (product: MockProduct) => "No " + product.name + " left";

const renderProduct = (
  product: MockProduct,
  quantityInCart: number = 0,
): { addToCart: ReturnType<typeof vi.fn> } => {
  const addToCart = vi.fn();

  render(
    <Product
      product={product}
      addToCart={addToCart}
      hasListItem={vi.fn(() => false)}
      quantityByProductInCart={vi.fn(() => quantityInCart)}
      addProductInterest={vi.fn().mockResolvedValue(undefined)}
    />,
  );

  return { addToCart };
};

// A product with no guestlist, so the card shows the add-to-cart button.
const plainProduct = (overrides = {}) =>
  createMockProduct({
    totalStock: 1,
    unitsSold: 0,
    soldOut: false,
    guestlists: null,
    ...overrides,
  });

const guestlistProduct = (overrides = {}) =>
  createMockProduct({
    totalStock: 1,
    unitsSold: 0,
    soldOut: false,
    guestlists: [{ id: 1, name: "Guestlist", typeCode: false, productId: 1 }],
    ...overrides,
  });

describe("Product stock affordance", () => {
  beforeEach(() => {
    vi.resetAllMocks();
    setup("fail");
  });

  describe("when the cart already holds the remaining stock", () => {
    it("disables the add-to-cart button", () => {
      const product = plainProduct();

      renderProduct(product, 1);

      expect(screen.getByLabelText(noStockLabel(product))).toBeDisabled();
    });

    it("does not add to the cart when the card is clicked", async () => {
      const product = plainProduct();
      const { addToCart } = renderProduct(product, 1);

      await userEvent.click(screen.getByTestId("product-card-" + product.id));

      expect(addToCart).not.toHaveBeenCalled();
    });

    it("keeps the guestlist button enabled", () => {
      const product = guestlistProduct();

      renderProduct(product, 1);

      expect(
        screen.getByLabelText("Show guestlist for " + product.name),
      ).toBeEnabled();
    });

    it("keeps Register interest enabled for a sold-out product", () => {
      // A sold-out product has no stock left either, so gating the whole card would
      // have taken the interest flow down with it.
      const product = plainProduct({ soldOut: true });

      renderProduct(product, 0);

      expect(
        screen.getByLabelText("Register interest in " + product.name),
      ).toBeEnabled();
    });
  });

  describe("when stock remains", () => {
    it("enables the add-to-cart button and adds on click", async () => {
      const product = plainProduct({ totalStock: 5 });
      const { addToCart } = renderProduct(product, 1);

      expect(screen.getByLabelText(addLabel(product))).toBeEnabled();

      await userEvent.click(screen.getByTestId("product-card-" + product.id));

      expect(addToCart).toHaveBeenCalledTimes(1);
    });
  });

  describe("under the ignore behaviour", () => {
    it("never disables the add-to-cart button", () => {
      setup("ignore");
      const product = plainProduct();

      renderProduct(product, 9);

      expect(screen.getByLabelText(addLabel(product))).toBeEnabled();
    });
  });
});

describe("Product card", () => {
  beforeEach(() => {
    vi.resetAllMocks();
    setup("fail");
  });

  it("is a single button showing name and price", () => {
    const product = plainProduct({ name: "Mate", totalStock: 0 });
    renderProduct(product);

    const card = screen.getByRole("button", { name: addLabel(product) });
    expect(card).toHaveAttribute("data-testid", "product-card-" + product.id);
    expect(card).toHaveTextContent("Mate");
    expect(card).toHaveTextContent(
      mockCurrency.format(product.grossPrice.toNumber()),
    );
  });

  describe("cart badge", () => {
    it("shows how many units are already in the cart", () => {
      const product = plainProduct({ totalStock: 0 });
      renderProduct(product, 2);

      expect(screen.getByTestId("cart-badge-" + product.id)).toHaveTextContent(
        "2",
      );
    });

    it("is hidden while the product is not in the cart", () => {
      const product = plainProduct({ totalStock: 0 });
      renderProduct(product, 0);

      expect(
        screen.queryByTestId("cart-badge-" + product.id),
      ).not.toBeInTheDocument();
    });
  });

  describe("stock pill", () => {
    it("is not shown for unlimited products", () => {
      const product = plainProduct({ totalStock: 0 });
      renderProduct(product);

      expect(
        screen.queryByTestId("stock-pill-" + product.id),
      ).not.toBeInTheDocument();
    });

    it("shows the units left, minus the ones in the cart", () => {
      const product = plainProduct({ totalStock: 30, unitsSold: 10 });
      renderProduct(product, 2);

      const pill = screen.getByTestId("stock-pill-" + product.id);
      expect(pill).toHaveTextContent("18 left");
      expect(pill).toHaveAttribute("data-tone", "normal");
    });

    it(`turns low at ${LOW_STOCK_THRESHOLD} units or less`, () => {
      const product = plainProduct({
        totalStock: 10,
        unitsSold: 10 - LOW_STOCK_THRESHOLD,
      });
      renderProduct(product);

      expect(screen.getByTestId("stock-pill-" + product.id)).toHaveAttribute(
        "data-tone",
        "low",
      );
    });

    it("turns red when nothing is left", () => {
      const product = plainProduct({ totalStock: 5, unitsSold: 5 });
      renderProduct(product);

      const pill = screen.getByTestId("stock-pill-" + product.id);
      expect(pill).toHaveTextContent("0 left");
      expect(pill).toHaveAttribute("data-tone", "none");
    });
  });

  it("marks guest-list products", () => {
    const product = guestlistProduct({ totalStock: 0 });
    renderProduct(product);

    expect(screen.getByTestId("product-card-" + product.id)).toHaveTextContent(
      "List",
    );
  });

  it("shows sold-out products with the number of interested people", () => {
    const product = plainProduct({ soldOut: true, soldOutRequestCount: 7 });
    renderProduct(product);

    expect(
      screen.getByRole("button", { name: "Register interest in " + product.name }),
    ).toHaveTextContent("Sold out · 7 interested");
    expect(
      screen.queryByTestId("stock-pill-" + product.id),
    ).not.toBeInTheDocument();
  });

  it("opens the interest dialog instead of adding a sold-out product", async () => {
    const product = plainProduct({ soldOut: true });
    const { addToCart } = renderProduct(product);

    await userEvent.click(screen.getByTestId("product-card-" + product.id));

    expect(addToCart).not.toHaveBeenCalled();
    expect(await screen.findByText(/currently sold out/i)).toBeInTheDocument();
  });

  it("applies extra classes, e.g. to start a new grid row", () => {
    const product = plainProduct({ totalStock: 0 });
    const addToCart = vi.fn();
    render(
      <Product
        product={product}
        addToCart={addToCart}
        hasListItem={vi.fn(() => false)}
        quantityByProductInCart={vi.fn(() => 0)}
        addProductInterest={vi.fn().mockResolvedValue(undefined)}
        className="col-start-1"
      />,
    );

    expect(screen.getByTestId("product-card-" + product.id)).toHaveClass(
      "col-start-1",
    );
  });
});
