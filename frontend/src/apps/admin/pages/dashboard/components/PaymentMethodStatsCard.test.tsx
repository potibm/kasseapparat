import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { AdminContext, DataProvider } from "react-admin";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";
import PaymentMethodStatsCard from "./PaymentMethodStatsCard";

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

describe("PaymentMethodStatsCard", () => {
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
          <PaymentMethodStatsCard />
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
          <PaymentMethodStatsCard />
        </AdminContext>
      </ConfigContext>,
    );

    await waitFor(() => {
      expect(screen.getByText("No purchases yet.")).toBeInTheDocument();
    });
  });

  it("renders payment method stats table with data", async () => {
    const mockData = [
      {
        id: "CASH",
        paymentMethod: "CASH",
        name: "Cash",
        purchaseCount: 5,
        totalNetPrice: "100.00",
        totalGrossPrice: "120.00",
      },
      {
        id: "CARD",
        paymentMethod: "CARD",
        name: "Credit Card",
        purchaseCount: 3,
        totalNetPrice: "200.00",
        totalGrossPrice: "240.00",
      },
    ];

    mockDataProvider.getList.mockResolvedValue({ data: mockData });

    render(
      <ConfigContext value={mockConfig}>
        <AdminContext
          dataProvider={mockDataProvider as unknown as DataProvider}
        >
          <PaymentMethodStatsCard />
        </AdminContext>
      </ConfigContext>,
    );

    await waitFor(() => {
      expect(screen.getByText("Payment Method Stats")).toBeInTheDocument();
    });

    expect(screen.getByText("Cash")).toBeInTheDocument();
    expect(screen.getByText("Credit Card")).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();
    expect(screen.getByText("Total")).toBeInTheDocument();
  });
});
