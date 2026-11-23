import React from "react";
import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  AdminContext,
  AuthProvider,
  ResourceContextProvider,
  testDataProvider,
} from "react-admin";
import { MemoryRouter } from "react-router";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";
import { ProductList } from "../../resources/products/ProductList";
import { PurchaseList } from "../../resources/purchases/PurchaseList";
import { GuestList } from "../../resources/guests/GuestList";
import { GuestlistList } from "../../resources/guestlists/GuestlistList";

vi.mock("./useIsSmallScreen", () => ({ useIsSmallScreen: () => true }));

const config = {
  currencyLocale: "en-US",
  currencyOptions: { style: "currency", currency: "USD" },
  dateLocale: "en-US",
  paymentMethods: [],
} as unknown as AppConfig;

const authProvider: AuthProvider = {
  login: () => Promise.resolve(),
  logout: () => Promise.resolve(),
  checkAuth: () => Promise.resolve(),
  checkError: () => Promise.resolve(),
  getPermissions: () => Promise.resolve("admin"),
};

const renderList = (
  resource: string,
  records: object[],
  list: React.ReactElement,
) =>
  render(
    <MemoryRouter>
      <AdminContext
        authProvider={authProvider}
        dataProvider={testDataProvider({
          getList: () =>
            Promise.resolve({ data: records, total: records.length }) as never,
        })}
      >
        <ConfigContext value={config}>
          <ResourceContextProvider value={resource}>
            {list}
          </ResourceContextProvider>
        </ConfigContext>
      </AdminContext>
    </MemoryRouter>,
  );

describe("admin lists on a phone", () => {
  it("lists products with price and stock", async () => {
    renderList(
      "products",
      [
        {
          id: 1,
          name: "Mug",
          grossPrice: 12,
          totalStock: 5,
          unitsSold: 3,
          soldOut: false,
          hidden: true,
        },
        { id: 2, name: "Mate", grossPrice: 3, unitsSold: 40, soldOut: true },
      ],
      <ProductList />,
    );

    expect(await screen.findByText("Mug")).toBeInTheDocument();
    expect(screen.getByText("$12.00 · 2 / 5 left")).toBeInTheDocument();
    expect(screen.getByText("Hidden")).toBeInTheDocument();
    expect(screen.getByText("$3.00 · 40 sold")).toBeInTheDocument();
    expect(screen.getByText("Sold out")).toBeInTheDocument();
  });

  it("lists purchases with total and items", async () => {
    renderList(
      "purchases",
      [
        {
          id: "p1",
          createdAt: "2026-11-23T18:42:00Z",
          totalGrossPrice: "41",
          purchaseItems: [{ quantity: 2, product: { name: "Mate" } }],
        },
      ],
      <PurchaseList />,
    );

    expect(await screen.findByText("$41.00")).toBeInTheDocument();
    expect(screen.getByText("2× Mate")).toBeInTheDocument();
  });

  it("lists guests with their list and arrival", async () => {
    renderList(
      "guests",
      [
        {
          id: 1,
          name: "Alex Müller",
          guestlist: { name: "Crew" },
          arrivedAt: "2026-11-23T18:42:00Z",
        },
        { id: 2, name: "Jana", guestlist: null },
      ],
      <GuestList />,
    );

    expect(await screen.findByText("Alex Müller")).toBeInTheDocument();
    expect(screen.getByText("Crew")).toBeInTheDocument();
    expect(screen.getByText(/^arrived /)).toBeInTheDocument();
    expect(screen.getByText("Jana")).toBeInTheDocument();
  });

  it("lists guest lists with their progress", async () => {
    renderList(
      "guestlists",
      [
        {
          id: 1,
          name: "Crew",
          product: { name: "Free" },
          entries: 20,
          arrivedEntries: 2,
        },
        { id: 2, name: "Empty", product: null },
      ],
      <GuestlistList />,
    );

    expect(await screen.findByText("Crew")).toBeInTheDocument();
    expect(screen.getByText("Free")).toBeInTheDocument();
    expect(screen.getByText("2 / 20 arrived")).toBeInTheDocument();
    expect(screen.getByText("No entries")).toBeInTheDocument();
  });
});
