import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom";
import ProductList from "./ProductList";
import * as ConfigHookModule from "@core/config/hooks/useConfig";

vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: vi.fn(),
}));

vi.mock("@pos/api/usePosApi", () => ({
  usePosApi: () => ({
    fetchGuestlistByProductId: vi.fn().mockResolvedValue([]),
  }),
}));

const mockCurrency = { format: (val: number) => `${val} €` };

beforeEach(() => {
  vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
    currency: mockCurrency,
  } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
});

const baseProps = {
  addToCart: vi.fn(),
  hasListItem: vi.fn(() => false),
  quantityByProductInCart: vi.fn(() => 0),
  addProductInterest: vi.fn().mockResolvedValue(undefined),
};

describe("ProductList", () => {
  it("shows a spinner while products are null and loading", () => {
    render(
      <ProductList
        products={null}
        loading
        addToCart={baseProps.addToCart}
        hasListItem={baseProps.hasListItem}
        quantityByProductInCart={baseProps.quantityByProductInCart}
        addProductInterest={baseProps.addProductInterest}
      />,
    );
    expect(screen.getByText(/Loading products/i)).toBeInTheDocument();
  });

  it("shows an error message when products are null and not loading", () => {
    render(
      <ProductList
        products={null}
        loading={false}
        addToCart={baseProps.addToCart}
        hasListItem={baseProps.hasListItem}
        quantityByProductInCart={baseProps.quantityByProductInCart}
        addProductInterest={baseProps.addProductInterest}
      />,
    );
    expect(screen.getByText(/could not be loaded/i)).toBeInTheDocument();
  });

  it("shows 'No products available.' when the list is empty", () => {
    render(
      <ProductList
        products={[]}
        addToCart={baseProps.addToCart}
        hasListItem={baseProps.hasListItem}
        quantityByProductInCart={baseProps.quantityByProductInCart}
        addProductInterest={baseProps.addProductInterest}
      />,
    );
    expect(screen.getByText(/No products available/i)).toBeInTheDocument();
  });

  it("renders one product card per product", () => {
    const products = [
      {
        id: 1,
        name: "Alpha",
        totalStock: 0,
        unitsSold: 0,
        soldOut: false,
        guestlists: null,
        grossPrice: { toNumber: () => 10 },
      },
      {
        id: 2,
        name: "Beta",
        totalStock: 0,
        unitsSold: 0,
        soldOut: false,
        guestlists: null,
        grossPrice: { toNumber: () => 20 },
      },
    ] as unknown as Parameters<typeof ProductList>[0]["products"];
    render(
      <ProductList
        products={products}
        addToCart={baseProps.addToCart}
        hasListItem={baseProps.hasListItem}
        quantityByProductInCart={baseProps.quantityByProductInCart}
        addProductInterest={baseProps.addProductInterest}
      />,
    );
    expect(screen.getAllByTestId(/product-card-/)).toHaveLength(2);
  });

  it("renders vanished products as placeholders", () => {
    const products = [
      {
        id: 1,
        name: "Ticket",
        totalStock: 5,
        unitsSold: 5,
        soldOut: false,
        vanished: true,
        guestlists: null,
        grossPrice: { toNumber: () => 10 },
      },
    ] as unknown as Parameters<typeof ProductList>[0]["products"];

    render(
      <ProductList
        products={products}
        addToCart={baseProps.addToCart}
        hasListItem={baseProps.hasListItem}
        quantityByProductInCart={baseProps.quantityByProductInCart}
        addProductInterest={baseProps.addProductInterest}
      />,
    );

    expect(screen.getByTestId("product-card-1")).toHaveAttribute(
      "data-vanished",
      "true",
    );
  });

  it("starts a new grid row after a product with wrapAfter", () => {
    const product = (id: number, wrapAfter: boolean) => ({
      id,
      name: "Product " + id,
      totalStock: 0,
      unitsSold: 0,
      soldOut: false,
      wrapAfter,
      guestlists: null,
      grossPrice: { toNumber: () => 10 },
    });
    const products = [
      product(1, false),
      product(2, true),
      product(3, false),
    ] as unknown as Parameters<typeof ProductList>[0]["products"];

    render(
      <ProductList
        products={products}
        addToCart={baseProps.addToCart}
        hasListItem={baseProps.hasListItem}
        quantityByProductInCart={baseProps.quantityByProductInCart}
        addProductInterest={baseProps.addProductInterest}
      />,
    );

    expect(screen.getByTestId("product-card-1")).not.toHaveClass("col-start-1");
    expect(screen.getByTestId("product-card-2")).not.toHaveClass("col-start-1");
    expect(screen.getByTestId("product-card-3")).toHaveClass("col-start-1");
  });
});
