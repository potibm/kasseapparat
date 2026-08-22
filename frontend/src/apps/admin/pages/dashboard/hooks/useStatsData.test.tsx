import { renderHook, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { AdminContext } from "react-admin";
import { useStatsData } from "./useStatsData";

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

describe("useStatsData", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("fetches data successfully", async () => {
    const mockData = [{ id: 1, name: "Test" }];
    mockDataProvider.getList.mockResolvedValue({ data: mockData });

    const { result } = renderHook(() => useStatsData("testResource"), {
      wrapper,
    });

    expect(result.current.loading).toBe(true);
    expect(result.current.data).toBeNull();

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.data).toEqual(mockData);
    expect(result.current.error).toBe(false);
  });

  it("handles fetch error", async () => {
    mockDataProvider.getList.mockRejectedValue(new Error("Fetch failed"));

    const { result } = renderHook(() => useStatsData("testResource"), {
      wrapper,
    });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.data).toEqual([]);
    expect(result.current.error).toBe(true);
  });

  it("uses custom sort parameters", async () => {
    mockDataProvider.getList.mockResolvedValue({ data: [] });

    renderHook(() => useStatsData("testResource", "customField", "DESC"), {
      wrapper,
    });

    await waitFor(() => {
      expect(mockDataProvider.getList).toHaveBeenCalledWith("testResource", {
        pagination: { page: 1, perPage: 100 },
        sort: { field: "customField", order: "DESC" },
        filter: {},
      });
    });
  });
});
