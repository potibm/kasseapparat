import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { Table, TableBody } from "flowbite-react";
import Decimal from "decimal.js";
import CartRow from "./CartRow";
import { CartItem } from "../../types/cart.types";
import { createMockGuest, createMockProduct } from "@pos/api/schemas.mocks";
import { OutOfStockBehavior } from "@core/config/types/config.types";

const currency = new Intl.NumberFormat("de-DE", {
  style: "currency",
  currency: "EUR",
});

const createCartItem = (overrides: Partial<CartItem> = {}): CartItem => {
  const product = createMockProduct({
    id: 7,
    name: "Mate",
    grossPrice: new Decimal(3),
    netPrice: new Decimal(2.52),
    vatAmount: new Decimal(0.48),
    totalStock: 0,
    unitsSold: 0,
  });
  const quantity = overrides.quantity ?? 2;

  return {
    ...product,
    quantity,
    listItems: [],
    totalNetPrice: product.netPrice.mul(quantity),
    totalGrossPrice: product.grossPrice.mul(quantity),
    totalVatAmount: product.vatAmount.mul(quantity),
    ...overrides,
  };
};

const renderRow = (
  cartElement: CartItem,
  outOfStockBehavior: OutOfStockBehavior = "ignore",
) => {
  const handlers = {
    removeFromCart: vi.fn(),
    increaseInCart: vi.fn(),
    decreaseInCart: vi.fn(),
  };

  render(
    <Table>
      <TableBody>
        <CartRow
          cartElement={cartElement}
          currency={currency}
          outOfStockBehavior={outOfStockBehavior}
          {...handlers}
        />
      </TableBody>
    </Table>,
  );

  return handlers;
};

describe("CartRow", () => {
  it("shows name, line total, unit price and quantity", () => {
    renderRow(createCartItem({ quantity: 2 }));

    expect(screen.getByText("Mate")).toBeInTheDocument();
    expect(screen.getByText("6,00 €")).toBeInTheDocument();
    expect(screen.getByText(/× 3,00 €/)).toBeInTheDocument();
    expect(screen.getByTestId("cart-quantity-7")).toHaveTextContent("2");
  });

  it("decreases by one unit when the quantity is above one", () => {
    const item = createCartItem({ quantity: 2 });
    const { decreaseInCart, removeFromCart } = renderRow(item);

    fireEvent.click(screen.getByRole("button", { name: "Remove one Mate" }));

    expect(decreaseInCart).toHaveBeenCalledWith(item);
    expect(removeFromCart).not.toHaveBeenCalled();
    expect(
      screen.queryByRole("button", { name: "Remove Mate from cart" }),
    ).not.toBeInTheDocument();
  });

  it("turns the minus button into a remove button at quantity one", () => {
    const item = createCartItem({ quantity: 1 });
    const { decreaseInCart, removeFromCart } = renderRow(item);

    fireEvent.click(
      screen.getByRole("button", { name: "Remove Mate from cart" }),
    );

    expect(removeFromCart).toHaveBeenCalledWith(item);
    expect(decreaseInCart).not.toHaveBeenCalled();
  });

  it("adds one unit with the plus button", () => {
    const item = createCartItem({ quantity: 1 });
    const { increaseInCart } = renderRow(item);

    fireEvent.click(screen.getByRole("button", { name: "Add one Mate" }));

    expect(increaseInCart).toHaveBeenCalledWith(item);
  });

  it("disables plus when the stock is used up and the behaviour restricts it", () => {
    renderRow(
      createCartItem({ quantity: 2, totalStock: 5, unitsSold: 3 }),
      "fail",
    );

    expect(screen.getByRole("button", { name: "Add one Mate" })).toBeDisabled();
  });

  it("keeps plus enabled under the ignore behaviour even without stock", () => {
    renderRow(
      createCartItem({ quantity: 2, totalStock: 5, unitsSold: 3 }),
      "ignore",
    );

    expect(screen.getByRole("button", { name: "Add one Mate" })).toBeEnabled();
  });

  it("only allows removing the whole line when it carries guest-list entries", () => {
    const guest = createMockGuest({ id: 1, name: "Alex Müller", code: null });
    const item = createCartItem({ quantity: 2, listItems: [guest] });
    const { removeFromCart } = renderRow(item);

    expect(screen.getByText("Alex Müller")).toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: "Remove one Mate" }),
    ).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Add one Mate" })).toBeDisabled();

    fireEvent.click(
      screen.getByRole("button", { name: "Remove Mate from cart" }),
    );
    expect(removeFromCart).toHaveBeenCalledWith(item);
  });

  it("shows the code instead of the name for coded guest-list entries", () => {
    const guest = createMockGuest({ id: 2, name: "Ticket", code: "ABC123" });
    renderRow(createCartItem({ quantity: 1, listItems: [guest] }));

    expect(screen.getByText("ABC123")).toBeInTheDocument();
  });
});
