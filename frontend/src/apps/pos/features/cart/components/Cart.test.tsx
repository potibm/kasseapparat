import { render, screen, fireEvent } from "@testing-library/react";
import { describe, expect, it, vi, beforeEach } from "vitest";
import Cart from "./Cart";
import { Cart as CartObject } from "../services/Cart";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import { Decimal } from "decimal.js";
import { Product } from "@pos/utils/api.schemas";

// Mock für useConfig
vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: vi.fn(),
}));

describe("Cart Component", () => {
  const mockCurrency = new Intl.NumberFormat("de-DE", {
    style: "currency",
    currency: "EUR",
  });

  const mockPaymentMethods = [
    { code: "CASH", name: "Cash" },
    { code: "CARD", name: "CC" },
  ];

  const mockProduct = {
    id: 1,
    name: "Test Coffee",
    netPrice: new Decimal(2.5),
    grossPrice: new Decimal(3),
    vatAmount: new Decimal(0.5),
  };

  const setupMocks = () => {
    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      currency: mockCurrency,
      paymentMethods: mockPaymentMethods,
    } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
  };

  beforeEach(() => {
    vi.resetAllMocks();
    setupMocks();
  });

  it("renders an empty cart correctly", () => {
    const emptyCart = new CartObject();

    render(
      <Cart
        cart={emptyCart}
        removeFromCart={vi.fn()}
        removeAllFromCart={vi.fn()}
        checkoutCart={vi.fn()}
        checkoutProcessing={null}
      />,
    );

    expect(screen.getByRole("cell", { name: /^total$/i })).toBeInTheDocument();
    // Total should show 0.00 when cart is empty
    expect(screen.getByTestId("cart-total-value")).toHaveTextContent(/0,00/);

    // The "Remove All" button should be disabled when the cart is empty
    const removeAllBtn = screen.getByRole("button", {
      name: /remove all items/i,
    });
    expect(removeAllBtn).toBeDisabled();
  });

  it("renders cart items and handles removal", () => {
    let fullCart = new CartObject();
    fullCart = fullCart.add(mockProduct as unknown as Product, 2); // 2x 3.00 = 6.00

    const mockRemoveOne = vi.fn();
    const mockRemoveAll = vi.fn();

    render(
      <Cart
        cart={fullCart}
        removeFromCart={mockRemoveOne}
        removeAllFromCart={mockRemoveAll}
        checkoutCart={vi.fn()}
        checkoutProcessing={null}
      />,
    );

    // Product name
    expect(screen.getByText(/Test Coffee/i)).toBeInTheDocument();
    // Tocal price for 2 items
    expect(screen.getByTestId("cart-total-value")).toHaveTextContent(/6,00/);

    // Click "remove all items" button
    const removeAllBtn = screen.getByRole("button", {
      name: /remove all items/i,
    });
    fireEvent.click(removeAllBtn);
    expect(mockRemoveAll).toHaveBeenCalledTimes(1);
  });

  it("shows processing state on checkout buttons", () => {
    const cart = new CartObject();
    cart.add(mockProduct as unknown as Product, 1);

    render(
      <Cart
        cart={cart}
        removeFromCart={vi.fn()}
        removeAllFromCart={vi.fn()}
        checkoutCart={vi.fn()}
        checkoutProcessing="cash" // simulates running checkout
      />,
    );
  });

  it("triggers the flash animation when cart quantity changes", async () => {
    vi.useFakeTimers();

    // rAF Mock
    vi.stubGlobal("requestAnimationFrame", (cb: FrameRequestCallback) =>
      setTimeout(() => cb(Date.now()), 16),
    );

    const cart = new CartObject();
    const { rerender } = render(
      <Cart
        cart={cart}
        removeFromCart={vi.fn()}
        removeAllFromCart={vi.fn()}
        checkoutCart={vi.fn()}
        checkoutProcessing={null}
      />,
    );

    const updatedCart = cart.add(mockProduct as unknown as Product, 1);
    expect(updatedCart.totalQuantity).toBeGreaterThan(cart.totalQuantity);

    rerender(
      <Cart
        cart={updatedCart}
        removeFromCart={vi.fn()}
        removeAllFromCart={vi.fn()}
        checkoutCart={vi.fn()}
        checkoutProcessing={null}
      />,
    );

    vi.advanceTimersByTime(20);

    const table = screen.getByTestId("cart-table");

    await vi.waitFor(() => {
      if (!table.className.includes("animate__pulse")) {
        throw new Error("Class not yet applied");
      }
    });

    expect(table).toHaveClass("animate__pulse");

    vi.advanceTimersByTime(500);

    await vi.waitFor(() => {
      if (table.className.includes("animate__pulse")) {
        throw new Error("Class still applied");
      }
    });

    vi.useRealTimers();
    vi.unstubAllGlobals();
  });
});
