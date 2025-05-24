import { describe, it, expect, vi, afterEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";

const mockUseConfigVersion = (v) => {
  vi.doMock("../provider/ConfigProvider", () => ({
    useConfig: () => ({ version: v }),
  }));
};

describe("Version", () => {
  afterEach(() => {
    vi.resetModules();
  });

  it("shows link for normal version", async () => {
    mockUseConfigVersion("1.2.3");
    const { Version } = await import("./Version");

    render(
      <MemoryRouter>
        <Version />
      </MemoryRouter>,
    );

    expect(screen.getByText(/Version 1\.2\.3/)).toBeInTheDocument();
    expect(screen.getByRole("link")).toHaveAttribute(
      "href",
      expect.stringContaining("v1.2.3"),
    );
  });

  it("marks version as outdated if GitHub version is newer", async () => {
    // mock local version
    vi.doMock("../provider/ConfigProvider", () => ({
      useConfig: () => ({ version: "1.2.3" }),
    }));

    // faking sessionStorage
    vi.stubGlobal("sessionStorage", {
      getItem: vi.fn(() => "1.3.0"),
      setItem: vi.fn(),
    });

    const { Version } = await import("./Version");

    render(
      <MemoryRouter>
        <Version />
      </MemoryRouter>,
    );

    const link = screen.getByRole("link");
    expect(link).toHaveClass("text-red-600");
    expect(link).toHaveAttribute(
      "title",
      expect.stringContaining("A newer version"),
    );
  });

  it("has the same version as the local one", async () => {
    // mock local version
    vi.doMock("../provider/ConfigProvider", () => ({
      useConfig: () => ({ version: "1.2.3" }),
    }));

    // faking sessionStorage
    vi.stubGlobal("sessionStorage", {
      getItem: vi.fn(() => "1.2.3"),
      setItem: vi.fn(),
    });

    const { Version } = await import("./Version");

    render(
      <MemoryRouter>
        <Version />
      </MemoryRouter>,
    );

    const link = screen.getByRole("link");
    expect(link).not.toHaveClass("text-red-600");
    expect(link).not.toHaveAttribute(
      "title",
      expect.stringContaining("A newer version"),
    );
  });

  it("marks version as outdated if GitHub version is newer via fetch", async () => {
    mockUseConfigVersion("1.2.3");

    // simulate no cached version in sessionStorage
    vi.stubGlobal("sessionStorage", {
      getItem: vi.fn(() => null),
      setItem: vi.fn(),
    });

    // setting up fetch mock
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve({
          json: () =>
            Promise.resolve({
              tag_name: "v1.3.0",
            }),
        }),
      ),
    );

    const { Version } = await import("./Version");

    render(
      <MemoryRouter>
        <Version />
      </MemoryRouter>,
    );

    // wait for dom to update
    await waitFor(() =>
      expect(screen.getByRole("link")).toHaveClass("text-red-600"),
    );

    expect(screen.getByRole("link")).toHaveAttribute(
      "title",
      expect.stringContaining("1.3.0"),
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
