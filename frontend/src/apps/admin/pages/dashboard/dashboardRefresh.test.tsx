import {
  act,
  fireEvent,
  render,
  renderHook,
  screen,
  waitFor,
} from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { AdminContext, DataProvider } from "react-admin";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";
import { AUTO_REFRESH_MS, useAutoRefresh } from "./hooks/useAutoRefresh";
import { useStatsData } from "./hooks/useStatsData";
import { DashboardRefreshContext } from "./context/DashboardRefreshContext";
import KpiCards from "./components/KpiCards";
import DashboardToolbar from "./components/DashboardToolbar";

const mockConfig = {
  currency: new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }),
  dateLocale: "en-US",
} as unknown as AppConfig;

const getList = vi.fn();
const dataProvider = { getList } as unknown as DataProvider;

describe("useAutoRefresh", () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => vi.useRealTimers());

  it("ticks every minute while enabled", () => {
    const { result } = renderHook(() => useAutoRefresh());

    expect(result.current.tick).toBe(0);

    act(() => {
      vi.advanceTimersByTime(AUTO_REFRESH_MS);
    });
    expect(result.current.tick).toBe(1);

    act(() => result.current.setEnabled(false));
    act(() => {
      vi.advanceTimersByTime(AUTO_REFRESH_MS * 3);
    });
    expect(result.current.tick).toBe(1);

    act(() => result.current.refreshNow());
    expect(result.current.tick).toBe(2);
  });
});

describe("useStatsData", () => {
  it("reloads on every refresh tick", async () => {
    getList.mockResolvedValue({ data: [{ id: 1 }] });
    let tick = 0;

    const { rerender } = renderHook(() => useStatsData("productStats"), {
      wrapper: ({ children }) => (
        <AdminContext dataProvider={dataProvider}>
          <DashboardRefreshContext.Provider value={tick}>
            {children}
          </DashboardRefreshContext.Provider>
        </AdminContext>
      ),
    });

    await waitFor(() => expect(getList).toHaveBeenCalledTimes(1));

    tick = 1;
    rerender();

    await waitFor(() => expect(getList).toHaveBeenCalledTimes(2));
  });
});

describe("KpiCards", () => {
  it("shows revenue, purchases and the average purchase", async () => {
    getList.mockResolvedValue({
      data: [
        { id: "CASH", purchaseCount: 3, totalGrossPrice: "30" },
        { id: "CC", purchaseCount: 1, totalGrossPrice: 50 },
      ],
    });

    render(
      <AdminContext dataProvider={dataProvider}>
        <ConfigContext.Provider value={mockConfig}>
          <KpiCards />
        </ConfigContext.Provider>
      </AdminContext>,
    );

    const kpis = await screen.findByTestId("dashboard-kpis");
    expect(kpis).toHaveTextContent("Revenue$80.00");
    expect(kpis).toHaveTextContent("Purchases4");
    expect(kpis).toHaveTextContent("Average purchase$20.00");
  });

  it("copes with no purchases", async () => {
    getList.mockResolvedValue({ data: [] });

    render(
      <AdminContext dataProvider={dataProvider}>
        <ConfigContext.Provider value={mockConfig}>
          <KpiCards />
        </ConfigContext.Provider>
      </AdminContext>,
    );

    expect(await screen.findByTestId("dashboard-kpis")).toHaveTextContent(
      "Average purchase$0.00",
    );
  });
});

describe("DashboardToolbar", () => {
  it("switches auto-refresh and refreshes on demand", () => {
    const setEnabled = vi.fn();
    const refreshNow = vi.fn();

    render(
      <ConfigContext.Provider value={mockConfig}>
        <DashboardToolbar
          enabled
          setEnabled={setEnabled}
          refreshNow={refreshNow}
          lastRefresh={new Date("2026-11-23T09:05:00")}
        />
      </ConfigContext.Provider>,
    );

    fireEvent.click(screen.getByLabelText("Auto-refresh every minute"));
    expect(setEnabled).toHaveBeenCalledWith(false);

    fireEvent.click(screen.getByRole("button", { name: "Refresh now" }));
    expect(refreshNow).toHaveBeenCalled();

    expect(screen.getByText(/Updated/)).toHaveTextContent("09:05");
  });
});
