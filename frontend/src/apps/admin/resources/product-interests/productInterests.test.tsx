import { describe, expect, it } from "vitest";
import { render, screen, within } from "@testing-library/react";
import { AdminContext, testDataProvider } from "react-admin";
import { productsByInterest } from "./interestSummary.utils";
import ProductInterestSummary from "./ProductInterestSummary";

const products = [
  { id: 1, name: "Mug", soldOutRequestCount: 2 },
  { id: 2, name: "Mate", soldOutRequestCount: 0 },
  { id: 3, name: "Hoodie", soldOutRequestCount: 5 },
  { id: 4, name: "Cap", soldOutRequestCount: 2 },
  { id: 5, name: "Shirt" },
];

describe("productsByInterest", () => {
  it("keeps requested products, most requested first", () => {
    expect(productsByInterest(products).map((p) => p.name)).toEqual([
      "Hoodie",
      "Cap",
      "Mug",
    ]);
  });
});

const renderSummary = (data: object[]) =>
  render(
    <AdminContext
      dataProvider={testDataProvider({
        getList: () => Promise.resolve({ data, total: data.length }) as never,
      })}
    >
      <ProductInterestSummary />
    </AdminContext>,
  );

describe("ProductInterestSummary", () => {
  it("lists the interest per product", async () => {
    renderSummary(products);

    const summary = await screen.findByTestId("interest-summary");
    const hoodie = await within(summary).findByText("Hoodie");
    expect(hoodie.closest("li")).toHaveTextContent("Hoodie5");
    expect(within(summary).queryByText("Mate")).not.toBeInTheDocument();
  });

  it("says when nobody asked yet", async () => {
    renderSummary([{ id: 2, name: "Mate", soldOutRequestCount: 0 }]);

    expect(
      await screen.findByText("Nobody asked for a sold-out product yet."),
    ).toBeInTheDocument();
  });
});
