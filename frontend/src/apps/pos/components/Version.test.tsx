import "@testing-library/jest-dom";
import { describe, it, expect, vi, afterEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { Version } from "./Version";
import * as ConfigModule from "../../../core/config/providers/ConfigProvider";

vi.mock("../../../core/config/providers/ConfigProvider", () => ({
  useConfig: vi.fn(),
}));

const renderVersionWithMocks = async ({
  localVersion = "1.2.3",
  cachedRemoteVersion = null as string | null,
  fetchRemoteVersion = null as string | null,
}) => {
  // Config Mock
  vi.mocked(ConfigModule.useConfig).mockReturnValue({
    version: localVersion,
  } as any);

  // Storage Stubbing
  vi.stubGlobal("sessionStorage", {
    getItem: vi.fn().mockReturnValue(cachedRemoteVersion),
    setItem: vi.fn(),
  });

  // Fetch Stubbing
  if (fetchRemoteVersion) {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ tag_name: fetchRemoteVersion }),
      }),
    );
  } else {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({}),
      }),
    );
  }

  const utils = render(
    <MemoryRouter>
      <Version />
    </MemoryRouter>,
  );

  return utils;
};

describe("Version", () => {
  afterEach(() => {
    vi.clearAllMocks();
    vi.resetAllMocks();
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it("marks version as outdated if GitHub version is newer via fetch", async () => {
    await renderVersionWithMocks({
      localVersion: "1.2.3",
      fetchRemoteVersion: "v1.3.0",
    });

    const link = await screen.findByRole("link", { name: /Version 1.2.3/i });

    await waitFor(
      () => {
        expect(link).toHaveClass("text-red-600");
      },
      { timeout: 2000 },
    );

    expect(link).toHaveAttribute("title", expect.stringContaining("1.3.0"));
  });

  it("does not flash or update if version is already cached", async () => {
    await renderVersionWithMocks({
      localVersion: "1.2.3",
      cachedRemoteVersion: "1.2.3",
    });

    const link = screen.getByRole("link");
    expect(link).not.toHaveClass("text-red-600");
  });

  it("shows plain text for internal version (no act warning here)", async () => {
    vi.mocked(ConfigModule.useConfig).mockReturnValue({
      version: "0.1.0",
    } as any);

    render(
      <MemoryRouter>
        <Version />
      </MemoryRouter>,
    );

    expect(screen.getByText(/Version 0\.1\.0/)).toBeInTheDocument();
    expect(screen.queryByRole("link")).not.toBeInTheDocument();
  });
});
