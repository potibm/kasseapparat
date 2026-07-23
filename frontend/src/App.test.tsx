import { render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import App from "./App";
import "@testing-library/jest-dom";

vi.stubGlobal(
  "fetch",
  vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({}),
  } as Response),
);

describe("App", () => {
  beforeEach(() => {
    vi.resetAllMocks();
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.resetAllMocks();
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it("renders the App component while auth is loading", async () => {
    render(<App />);

    await waitFor(
      () => {
        const element = screen.getByText(/loading kasseapparat/i);
        expect(element).toBeInTheDocument();
      },
      { timeout: 2000 },
    );
  });

  it("renders the App component without crashing", () => {
    render(<App />);
    expect(document.body).toBeInTheDocument();
  });
});
