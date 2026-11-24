import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { AppConfig } from "@core/config/types/config.types";
import { ConnectionStatus } from "../hooks/useConnectionStatus";

vi.mock("@core/config/hooks/useConfig");

const mockUseConnectionStatus = vi.fn<(url: string) => ConnectionStatus>();
vi.mock("../hooks/useConnectionStatus", () => ({
  useConnectionStatus: (url: string) => mockUseConnectionStatus(url),
}));

import { useConfig } from "@core/config/hooks/useConfig";
import ConnectionIndicator from "./ConnectionIndicator";

describe("ConnectionIndicator", () => {
  beforeEach(() => {
    vi.mocked(useConfig).mockReturnValue({
      apiBaseUrl: "http://localhost:3001/api/v3",
    } as unknown as AppConfig);
  });

  it("checks the public config endpoint of the API", () => {
    mockUseConnectionStatus.mockReturnValue("online");

    render(<ConnectionIndicator />);

    expect(mockUseConnectionStatus).toHaveBeenCalledWith(
      "http://localhost:3001/api/v3/config",
    );
  });

  it.each([
    ["online", "Online"],
    ["slow", "Slow"],
    ["offline", "Offline"],
  ] as const)("shows %s as %s", (status, label) => {
    mockUseConnectionStatus.mockReturnValue(status);

    render(<ConnectionIndicator />);

    const indicator = screen.getByRole("status", {
      name: `Connection: ${label}`,
    });
    expect(indicator).toHaveTextContent(label);
    expect(indicator).toHaveAttribute("data-status", status);
  });
});
