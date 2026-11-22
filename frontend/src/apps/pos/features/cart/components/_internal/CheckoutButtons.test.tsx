import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Decimal from "decimal.js";
import CheckoutButtons from "./CheckoutButtons";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import { Cart } from "../../services/Cart";

vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: vi.fn(),
}));

vi.mock("@core/localstorage/helper/local-storage-reader", () => ({
  getCurrentReaderId: () => undefined,
}));

const cart = (total: number, isEmpty = false) =>
  ({ totalGross: new Decimal(total), isEmpty }) as unknown as Cart;

const renderButtons = (total = 37, isEmpty = false) => {
  const handleCheckoutCart = vi.fn();

  render(
    <CheckoutButtons
      cart={cart(total, isEmpty)}
      checkoutProcessing={null}
      handleCheckoutCart={handleCheckoutCart}
    />,
  );

  return { handleCheckoutCart };
};

describe("CheckoutButtons", () => {
  beforeEach(() => {
    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      currency: { format: (value: number) => `${value} €` },
      paymentMethods: [
        { code: "CASH", name: "Cash", changeCalculator: true },
        { code: "CC", name: "Card", changeCalculator: false },
      ],
    } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
  });

  it("books a payment method without calculator right away", async () => {
    const { handleCheckoutCart } = renderButtons();

    await userEvent.click(screen.getByTestId("checkout-button-CC"));

    expect(handleCheckoutCart).toHaveBeenCalledWith("CC", { type: "empty" });
    expect(
      screen.queryByTestId("change-calculator-confirm"),
    ).not.toBeInTheDocument();
  });

  it("asks for the amount received before booking cash", async () => {
    const { handleCheckoutCart } = renderButtons();

    await userEvent.click(screen.getByTestId("checkout-button-CASH"));

    expect(handleCheckoutCart).not.toHaveBeenCalled();
    expect(
      await screen.findByTestId("change-calculator-total"),
    ).toHaveTextContent("37 €");

    await userEvent.click(screen.getByTestId("change-calculator-confirm"));

    expect(handleCheckoutCart).toHaveBeenCalledWith("CASH", { type: "empty" });
    expect(
      screen.queryByTestId("change-calculator-confirm"),
    ).not.toBeInTheDocument();
  });

  it("does not book when the calculator is cancelled", async () => {
    const { handleCheckoutCart } = renderButtons();

    await userEvent.click(screen.getByTestId("checkout-button-CASH"));
    await userEvent.click(
      await screen.findByTestId("change-calculator-cancel"),
    );

    expect(handleCheckoutCart).not.toHaveBeenCalled();
    expect(
      screen.queryByTestId("change-calculator-cancel"),
    ).not.toBeInTheDocument();
  });

  it("disables the buttons for an empty cart", () => {
    renderButtons(0, true);

    expect(screen.getByTestId("checkout-button-CASH")).toBeDisabled();
    expect(
      screen.getByRole("button", { name: "Pay with Card" }),
    ).toBeDisabled();
  });
});
