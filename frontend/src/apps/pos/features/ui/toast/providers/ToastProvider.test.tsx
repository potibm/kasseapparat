import { render, screen, act, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { ToastProvider } from "./ToastProvider";
import { useToast } from "../hooks/useToast"; // Pfad ggf. anpassen
import "@testing-library/jest-dom";

vi.stubGlobal(
  "fetch",
  vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({}),
  }),
);

const mockUndo = vi.fn();

const TestComponent = () => {
  const { showToast } = useToast();

  return (
    <div>
      <button
        onClick={() => showToast({ severity: "success", message: "YAY!" })}
      >
        Normal Toast
      </button>
      <button
        onClick={() =>
          showToast({ severity: "error", message: "CRITICAL!", blocking: true })
        }
      >
        Blocking Toast
      </button>
      <button
        onClick={() =>
          showToast({
            severity: "info",
            message: "Cart cleared",
            action: { label: "Undo", onClick: mockUndo },
          })
        }
      >
        Action Toast
      </button>
    </div>
  );
};

const renderWithProvider = () => {
  return render(
    <ToastProvider>
      <TestComponent />
    </ToastProvider>,
  );
};

describe("Toast System Integration", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.clearAllTimers();
    vi.useRealTimers();
    vi.unstubAllGlobals();
  });

  it("render a toast when showToast is called", () => {
    renderWithProvider();

    fireEvent.click(screen.getByText("Normal Toast"));

    expect(screen.getByText("YAY!")).toBeInTheDocument();
  });

  it("removes the toast automatically after the time has elapsed (Auto-Close)", () => {
    renderWithProvider();

    fireEvent.click(screen.getByText("Normal Toast"));
    expect(screen.getByText("YAY!")).toBeInTheDocument();

    // fast forward time by 10 seconds (10000 ms)
    act(() => {
      vi.advanceTimersByTime(10000);
    });

    expect(screen.queryByText("YAY!")).not.toBeInTheDocument();
  });

  it("renders an action button that runs the action and dismisses the toast", () => {
    renderWithProvider();

    fireEvent.click(screen.getByText("Action Toast"));
    expect(screen.getByText("Cart cleared")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Undo" }));

    expect(mockUndo).toHaveBeenCalledTimes(1);
    expect(screen.queryByText("Cart cleared")).not.toBeInTheDocument();
  });

  it("renders no action button for toasts without an action", () => {
    renderWithProvider();

    fireEvent.click(screen.getByText("Normal Toast"));

    expect(
      screen.queryByRole("button", { name: "Undo" }),
    ).not.toBeInTheDocument();
  });

  it("renders a backdrop when a blocking toast is active", () => {
    renderWithProvider();

    fireEvent.click(screen.getByText("Blocking Toast"));

    const backdrop = document.querySelector(".z-9998");
    expect(backdrop).toBeInTheDocument();
  });
});
