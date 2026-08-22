import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { AdminContext, DataProvider } from "react-admin";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";
import ProductStatsCard from "./ProductStatsCard";

const mockConfig: AppConfig = {
  version: "1.0.0",
  apiHost: "http://localhost:3001",
  apiBaseUrl: "http://localhost:3001/api/v3",
  websocketHost: "ws://localhost:3001",
  websocketBaseUrl: "ws://localhost:3001/api/v3",
  locale: "en",
  currencyCode: "USD",
  currencyLocale: "en-US",
  currency: new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }),
  currencyOptions: {},
  dateLocale: "en-US",
  dateOptions: {},
  vatRates: [],
  paymentMethods: [],
  sumupEnabled: false,
  authMode: "oidc",
};

const mockDataProvider = {
  getList: vi.fn(),
  getOne: vi.fn(),
  getMany: vi.fn(),
  getManyReference: vi.fn(),
  create: vi.fn(),
  update: vi.fn(),
  updateMany: vi.fn(),
  delete: vi.fn(),
  deleteMany: vi.fn(),
};

describe("ProductStatsCard", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders loading state initially", async () => {
    mockDataProvider.getList.mockImplementation(
      () => new Promise(() => {}), // Never resolves
    );

    render(
      <ConfigContext value={mockConfig}>
        <AdminContext
          dataProvider={mockDataProvider as unknown as DataProvider}
        >
          <ProductStatsCard />
        </AdminContext>
      </ConfigContext>,
    );

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("renders empty state when no data", async () => {
    mockDataProvider.getList.mockResolvedValue({ data: [] });

    render(
      <ConfigContext value={mockConfig}>
        <AdminContext
          dataProvider={mockDataProvider as unknown as DataProvider}
        >
          <ProductStatsCard />
        </AdminContext>
      </ConfigContext>,
    );

    await waitFor(() => {
      expect(screen.getByText("No products yet.")).toBeInTheDocument();
    });
  });

  it("renders product stats table with data", async () => {
    const mockData = [
      {
        id: 1,
        name: "Product A",
        soldItems: 10,
        totalNetPrice: "50.00",
        totalGrossPrice: "60.00",
      },
      {
        id: 2,
        name: "Product B",
        soldItems: 5,
        totalNetPrice: "25.00",
        totalGrossPrice: "30.00",
      },
    ];

    mockDataProvider.getList.mockResolvedValue({ data: mockData });

    render(
      <ConfigContext value={mockConfig}>
        <AdminContext
          dataProvider={mockDataProvider as unknown as DataProvider}
        >
          <ProductStatsCard />
        </AdminContext>
      </ConfigContext>,
    );

    await waitFor(() => {
      expect(screen.getByText("Product Sales Stats")).toBeInTheDocument();
    });

    expect(screen.getByText("Product A")).toBeInTheDocument();
    expect(screen.getByText("Product B")).toBeInTheDocument();
    expect(screen.getByText("10")).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
    expect(screen.getByText("Total")).toBeInTheDocument();
  });
});
