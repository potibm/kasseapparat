import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { ConfigProvider } from "./ConfigProvider";
import { useConfig } from "../hooks/useConfig";
import "@testing-library/jest-dom";

const TestChild = () => {
  const config = useConfig();
  return <div data-testid="config-version">Version: {config.version}</div>;
};

describe("ConfigProvider", () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it("renders loading state initially", () => {
    vi.stubGlobal("fetch", vi.fn().mockReturnValue(new Promise(() => {})));

    render(
      <ConfigProvider>
        <TestChild />
      </ConfigProvider>,
    );

    expect(screen.getByText("⏳ Loading config...")).toBeInTheDocument();
  });

  it("renders children with config on successful fetch", async () => {
    const mockRawConfig = {
      version: "2.0.0",
      dateOptions: {
        year: "numeric",
      },
      vatRates: [{ name: "STANDARD", rate: 20 }],
      paymentMethods: [{ code: "CASH", name: "Cash" }],
    };

    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockRawConfig),
      }),
    );

    render(
      <ConfigProvider>
        <TestChild />
      </ConfigProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("config-version")).toHaveTextContent(
        "Version: 2.0.0",
      );
    });
  });

  it("renders error state on fetch failure (e.g. 500 Internal Server Error)", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
      }),
    );

    render(
      <ConfigProvider>
        <TestChild />
      </ConfigProvider>,
    );

    await waitFor(() => {
      expect(screen.getByText(/Error loading config:/i)).toBeInTheDocument();
      expect(screen.getByText("HTTP error! status: 500")).toBeInTheDocument();
    });
  });

  it("renders error state on schema validation failure", async () => {
    const invalidRawConfig = {
      version: 12345,
      paymentMethods: "No Array",
      dateOptions: "{ defect json...",
    };

    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(invalidRawConfig),
      }),
    );

    const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

    render(
      <ConfigProvider>
        <TestChild />
      </ConfigProvider>,
    );

    await waitFor(() => {
      expect(screen.getByText(/Error loading config:/i)).toBeInTheDocument();
    });

    consoleSpy.mockRestore();
  });
});
