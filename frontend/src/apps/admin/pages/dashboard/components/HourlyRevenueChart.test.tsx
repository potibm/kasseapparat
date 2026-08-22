import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { AdminContext, DataProvider } from "react-admin";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";
import HourlyRevenueChart from "./HourlyRevenueChart";

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

describe("HourlyRevenueChart", () => {
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
          <HourlyRevenueChart />
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
          <HourlyRevenueChart />
        </AdminContext>
      </ConfigContext>,
    );

    await waitFor(() => {
      expect(screen.getByText("No data available.")).toBeInTheDocument();
    });
  });

  it("renders chart with data", async () => {
    const mockData = [
      {
        id: "2024-01-01 10:00_CASH",
        timeBucket: "2024-01-01 10:00",
        paymentMethod: "CASH",
        name: "Cash",
        totalGrossPrice: "100.00",
      },
      {
        id: "2024-01-01 11:00_CASH",
        timeBucket: "2024-01-01 11:00",
        paymentMethod: "CASH",
        name: "Cash",
        totalGrossPrice: "150.00",
      },
    ];

    mockDataProvider.getList.mockResolvedValue({ data: mockData });

    render(
      <ConfigContext value={mockConfig}>
        <AdminContext
          dataProvider={mockDataProvider as unknown as DataProvider}
        >
          <HourlyRevenueChart />
        </AdminContext>
      </ConfigContext>,
    );

    await waitFor(() => {
      expect(
        screen.getByText("Revenue by Payment Method (Last 3 Days)"),
      ).toBeInTheDocument();
    });
  });
});
