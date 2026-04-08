import { renderHook } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { useConfig } from "./useConfig";
import { ConfigContext } from "../context/ConfigContext";
import { AppConfig } from "../types/config.types";
import { ReactNode } from "react";

describe("useConfig", () => {
  it("throws an error if used outside of ConfigProvider", () => {
    const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

    expect(() => renderHook(() => useConfig())).toThrow(
      "useConfig must be used within a ConfigProvider",
    );

    consoleSpy.mockRestore();
  });

  it("returns context value when used within ConfigProvider", () => {
    const mockConfig = {
      version: "1.2.3",
      apiHost: "http://localhost",
    } as unknown as AppConfig;

    const wrapper = ({ children }: { children: ReactNode }) => (
      <ConfigContext.Provider value={mockConfig}>
        {children}
      </ConfigContext.Provider>
    );

    const { result } = renderHook(() => useConfig(), { wrapper });

    expect(result.current).toEqual(mockConfig);
  });
});
