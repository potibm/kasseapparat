import { renderHook, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { AdminContext } from "react-admin";
import { useHourlyChartData } from "./useHourlyChartData";

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

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <AdminContext dataProvider={mockDataProvider as any}>{children}</AdminContext>
);

describe("useHourlyChartData", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("returns empty state when no data", async () => {
    mockDataProvider.getList.mockResolvedValue({ data: [] });

    const { result } = renderHook(
      () =>
        useHourlyChartData({
          resource: "testResource",
          dataKey: "category",
          aggregateFn: (acc, item: any) => acc + item.value,
        }),
      { wrapper },
    );

    await waitFor(() => {
      expect(result.current.data).toEqual([]);
    });

    expect(result.current.chartData).toEqual([]);
    expect(result.current.seriesKeys).toEqual([]);
  });

  it("processes data correctly", async () => {
    const mockData = [
      {
        id: "1",
        timeBucket: "2024-01-15 10:00",
        category: "A",
        value: 10,
      },
      {
        id: "2",
        timeBucket: "2024-01-15 11:00",
        category: "A",
        value: 20,
      },
      {
        id: "3",
        timeBucket: "2024-01-15 10:00",
        category: "B",
        value: 5,
      },
    ];

    mockDataProvider.getList.mockResolvedValue({ data: mockData });

    const { result } = renderHook(
      () =>
        useHourlyChartData({
          resource: "testResource",
          dataKey: "category",
          aggregateFn: (acc, item: any) => acc + item.value,
        }),
      { wrapper },
    );

    await waitFor(() => {
      expect(result.current.data).toHaveLength(3);
    });

    expect(result.current.seriesKeys).toHaveLength(2);
    expect(result.current.seriesKeys).toContain("A");
    expect(result.current.seriesKeys).toContain("B");
    expect(result.current.chartData.length).toBeGreaterThan(0);
  });

  it("handles fetch error", async () => {
    mockDataProvider.getList.mockRejectedValue(new Error("Fetch failed"));

    const { result } = renderHook(
      () =>
        useHourlyChartData({
          resource: "testResource",
          dataKey: "category",
          aggregateFn: (acc, item: any) => acc + item.value,
        }),
      { wrapper },
    );

    await waitFor(() => {
      expect(result.current.data).toEqual([]);
    });

    expect(result.current.chartData).toEqual([]);
  });
});
