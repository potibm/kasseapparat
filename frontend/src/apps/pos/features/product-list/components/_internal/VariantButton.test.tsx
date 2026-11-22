import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import VariantButton from "./VariantButton";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import { createMockProduct } from "@pos/api/schemas.mocks";
import { OutOfStockBehavior } from "@core/config/types/config.types";

vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: vi.fn(),
}));

vi.mock("@pos/api/usePosApi", () => ({
  usePosApi: () => ({
    fetchGuestlistByProductId: vi.fn().mockResolvedValue([]),
  }),
}));

const mockCurrency = { format: (val: number) => `${val} €` };

const setup = (outOfStockBehavior: OutOfStockBehavior = "ignore") => {
  vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
    currency: mockCurrency,
    outOfStockBehavior,
  } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
};

const shirt = (overrides = {}) =>
  createMockProduct({
    id: 7,
    name: "👕 Male L",
    variant: "L",
    productGroupId: 10,
    soldOut: false,
    totalStock: 0,
    unitsSold: 0,
    guestlists: null,
    ...overrides,
  });

const renderVariant = (
  product: ReturnType<typeof createMockProduct>,
  {
    quantityInCart = 0,
    showPrice = false,
    vanished = false,
  }: { quantityInCart?: number; showPrice?: boolean; vanished?: boolean } = {},
) => {
  const addToCart = vi.fn();

  render(
    <VariantButton
      product={product}
      addToCart={addToCart}
      hasListItem={vi.fn(() => false)}
      quantityByProductInCart={vi.fn(() => quantityInCart)}
      addProductInterest={vi.fn().mockResolvedValue(undefined)}
      showPrice={showPrice}
      vanished={vanished}
    />,
  );

  return { addToCart };
};

describe("VariantButton", () => {
  beforeEach(() => setup());

  it("shows the variant and adds the product on tap", async () => {
    const product = shirt();
    const { addToCart } = renderVariant(product);

    const button = screen.getByTestId("product-card-7");
    expect(button).toHaveTextContent("L");
    expect(button).not.toHaveTextContent("€");

    await userEvent.click(button);

    expect(addToCart).toHaveBeenCalledWith(product, 1, null);
  });

  it("falls back to the product name without a variant", () => {
    renderVariant(shirt({ variant: "" }));

    expect(screen.getByTestId("product-card-7")).toHaveTextContent("👕 Male L");
  });

  it("shows the price when the row has no shared price", () => {
    const product = shirt();
    renderVariant(product, { showPrice: true });

    expect(screen.getByTestId("product-card-7")).toHaveTextContent(
      mockCurrency.format(product.grossPrice.toNumber()),
    );
  });

  it("shows the stock left and the cart quantity", () => {
    renderVariant(shirt({ totalStock: 10, unitsSold: 7 }), {
      quantityInCart: 1,
    });

    expect(screen.getByTestId("variant-status-7")).toHaveTextContent("2 left");
    expect(screen.getByTestId("cart-badge-7")).toHaveTextContent("1");
  });

  it("is disabled once the stock is in the cart under fail", () => {
    setup("fail");
    renderVariant(shirt({ totalStock: 2, unitsSold: 0 }), {
      quantityInCart: 2,
    });

    const button = screen.getByRole("button", { name: "No 👕 Male L left" });
    expect(button).toBeDisabled();
  });

  it("offers to register interest when sold out", () => {
    renderVariant(shirt({ soldOut: true }));

    expect(
      screen.getByRole("button", { name: "Register interest in 👕 Male L" }),
    ).toBeEnabled();
    expect(screen.getByTestId("variant-status-7")).toHaveTextContent(
      "sold out",
    );
  });

  it("keeps a vanished variant as a disabled placeholder", () => {
    renderVariant(shirt(), { vanished: true });

    const button = screen.getByRole("button", {
      name: "👕 Male L is sold out",
    });
    expect(button).toBeDisabled();
    expect(button).toHaveAttribute("data-vanished", "true");
  });

  it("marks a product with a guest list", () => {
    renderVariant(
      shirt({
        guestlists: [{ id: 1, name: "List", typeCode: false, productId: 7 }],
      }),
    );

    expect(screen.getByTestId("variant-status-7")).toHaveTextContent("List");
  });
});
