import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { RecordContextProvider } from "react-admin";
import { summarizeStock } from "./stock.utils";
import ProductStockField from "./ProductStockField";
import { Product } from "./types";

describe("summarizeStock", () => {
  it("treats a total of 0 as unlimited", () => {
    expect(summarizeStock({ totalStock: 0, unitsSold: 12 })).toEqual({
      left: null,
      sold: 12,
      total: null,
      tone: "unlimited",
    });
    expect(summarizeStock({})).toEqual({
      left: null,
      sold: 0,
      total: null,
      tone: "unlimited",
    });
  });

  it("counts the units left and rates them", () => {
    expect(summarizeStock({ totalStock: 30, unitsSold: 10 }).tone).toBe(
      "normal",
    );
    expect(summarizeStock({ totalStock: 30, unitsSold: 27 })).toEqual({
      left: 3,
      sold: 27,
      total: 30,
      tone: "low",
    });
    expect(summarizeStock({ totalStock: 30, unitsSold: 31 })).toEqual({
      left: 0,
      sold: 31,
      total: 30,
      tone: "out",
    });
  });
});

const product = (overrides: Partial<Product>): Product => ({
  id: 1,
  name: "Mug",
  vatRate: 19,
  netPrice: 1,
  grossPrice: 1.19,
  pos: 1,
  sectionId: 1,
  soldOut: false,
  hidden: false,
  ...overrides,
});

const renderField = (record: Product) =>
  render(
    <RecordContextProvider value={record}>
      <ProductStockField />
    </RecordContextProvider>,
  );

describe("ProductStockField", () => {
  it("shows the units left of a limited product", () => {
    renderField(product({ totalStock: 5, unitsSold: 3 }));

    expect(screen.getByTestId("product-stock")).toHaveTextContent("2 / 5 left");
  });

  it("shows the units sold of an unlimited product and the interest", () => {
    renderField(product({ unitsSold: 42, soldOutRequestCount: 4 }));

    expect(screen.getByText("42 sold")).toBeInTheDocument();
    expect(screen.getByText("4 interested")).toBeInTheDocument();
    expect(screen.queryByTestId("product-stock")).not.toBeInTheDocument();
  });

  it("renders nothing without a record", () => {
    const { container } = render(<ProductStockField />);

    expect(container).toBeEmptyDOMElement();
  });
});
