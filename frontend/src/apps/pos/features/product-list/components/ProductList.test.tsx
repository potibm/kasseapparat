import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, within } from "@testing-library/react";
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

  it("shows the products under their section headings", () => {
    const product = (id: number, sectionId: number) => ({
      id,
      name: "Product " + id,
      totalStock: 0,
      unitsSold: 0,
      soldOut: false,
      sectionId,
      guestlists: null,
      grossPrice: { toNumber: () => 10 },
    });
    const products = [
      product(1, 2),
      product(2, 1),
      product(3, 3),
    ] as unknown as Parameters<typeof ProductList>[0]["products"];

    render(
      <ProductList
        products={products}
        sections={[
          { id: 1, name: "Tickets", pos: 1 },
          { id: 2, name: "Merch", pos: 2 },
          { id: 3, name: "", pos: 3 },
          { id: 4, name: "Empty", pos: 4 },
        ]}
        addToCart={baseProps.addToCart}
        hasListItem={baseProps.hasListItem}
        quantityByProductInCart={baseProps.quantityByProductInCart}
        addProductInterest={baseProps.addProductInterest}
      />,
    );

    const headings = screen.getAllByRole("heading").map((h) => h.textContent);
    expect(headings).toEqual(["Tickets", "Merch"]);

    const tickets = screen.getByRole("region", { name: "Tickets" });
    expect(within(tickets).getByTestId("product-card-2")).toBeInTheDocument();

    const merch = screen.getByRole("region", { name: "Merch" });
    expect(within(merch).getByTestId("product-card-1")).toBeInTheDocument();

    // the unnamed section is a plain row break: own grid, no heading
    const rowBreak = screen.getByTestId("section-3");
    expect(within(rowBreak).queryByRole("heading")).not.toBeInTheDocument();
    expect(within(rowBreak).getByTestId("product-card-3")).toBeInTheDocument();

    expect(screen.queryByTestId("section-4")).not.toBeInTheDocument();
  });

  it("shows all products without headings when no sections are known", () => {
    const products = [
      { id: 1, name: "Product 1", sectionId: 1 },
    ].map((p) => ({
      ...p,
      totalStock: 0,
      unitsSold: 0,
      soldOut: false,
      guestlists: null,
      grossPrice: { toNumber: () => 10 },
    })) as unknown as Parameters<typeof ProductList>[0]["products"];

    render(
      <ProductList
        products={products}
        addToCart={baseProps.addToCart}
        hasListItem={baseProps.hasListItem}
        quantityByProductInCart={baseProps.quantityByProductInCart}
        addProductInterest={baseProps.addProductInterest}
      />,
    );

    expect(screen.getByTestId("product-card-1")).toBeInTheDocument();
    expect(screen.queryByRole("heading")).not.toBeInTheDocument();
  });
});
