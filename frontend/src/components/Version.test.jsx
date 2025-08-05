import { describe, it, expect, vi, afterEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";

const mockUseConfigVersion = (v) => {
  return vi.doMock("../provider/ConfigProvider", () => ({
    useConfig: () => ({ version: v }),
  }));
};

const renderVersionWithMocks = async ({
  localVersion = "1.2.3",
  cachedRemoteVersion = null,
  fetchRemoteVersion = null,
}) => {
  vi.doMock("../provider/ConfigProvider", () => ({
    useConfig: () => ({ version: localVersion }),
  }));

  vi.stubGlobal("sessionStorage", {
    getItem: vi.fn(() => cachedRemoteVersion),
    setItem: vi.fn(),
  });

  if (fetchRemoteVersion) {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve({
          json: () => Promise.resolve({ tag_name: fetchRemoteVersion }),
        }),
      ),
    );
  }

  const { Version } = await import("./Version");

  render(
    <MemoryRouter>
      <Version />
    </MemoryRouter>,
  );

  return screen.getByRole("link");
};

describe("Version", () => {
  afterEach(() => {
    vi.resetModules();
    vi.restoreAllMocks();
  });

  it("shows link for normal version", async () => {
    const link = await renderVersionWithMocks({});
    expect(screen.getByText(/Version 1\.2\.3/)).toBeInTheDocument();
    expect(link).toHaveAttribute("href", expect.stringContaining("1.2.3"));
  });

  it("marks version as outdated if GitHub version is newer", async () => {
    const link = await renderVersionWithMocks({
      localVersion: "1.2.3",
      cachedRemoteVersion: "1.4.0",
    });

    expect(link).toHaveClass("text-red-600");
    expect(link).toHaveAttribute("title", expect.stringContaining("1.4.0"));
    expect(link).toHaveAttribute(
      "title",
      expect.stringContaining("A newer version"),
    );
  });

  it("has the same version as the local one", async () => {
    const link = await renderVersionWithMocks({
      localVersion: "1.2.3",
      cachedRemoteVersion: "1.2.3",
    });

    expect(link).not.toHaveClass("text-red-600");
    expect(link).not.toHaveAttribute(
      "title",
      expect.stringContaining("A newer version"),
    );
  });

  it("marks version as outdated if GitHub version is newer via fetch", async () => {
    const link = await renderVersionWithMocks({
      localVersion: "1.2.3",
      fetchRemoteVersion: "v1.3.0",
    });

    await waitFor(() => {
      expect(link).toHaveClass("text-red-600");
    });
    expect(link).toHaveAttribute("title", expect.stringContaining("1.3.0"));
    expect(link).toHaveAttribute(
      "title",
      expect.stringContaining("A newer version"),
    );
  });

  it("shows plain text for internal version (starts with 0)", async () => {
    mockUseConfigVersion("0.1.0");
    const { Version } = await import("./Version");
    render(<Version />);
    expect(screen.getByText(/Version 0\.1\.0/)).toBeInTheDocument();
    expect(screen.queryByRole("link")).not.toBeInTheDocument();
  });

  it("handles missing version", async () => {
    mockUseConfigVersion(null);
    const { Version } = await import("./Version");
    render(<Version />);
    expect(screen.getByText(/Version/)).toBeInTheDocument();
  });
});
