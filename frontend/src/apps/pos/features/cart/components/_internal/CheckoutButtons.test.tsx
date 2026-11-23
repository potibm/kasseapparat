import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Decimal from "decimal.js";
import CheckoutButtons from "./CheckoutButtons";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import { Cart } from "../../services/Cart";
import {
  SumupReaderContext,
  SumupReaderContextValue,
} from "@pos/features/sumup-reader/context/SumupReaderContext";

vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: vi.fn(),
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

  it("books a free cart with cash right away", async () => {
    const { handleCheckoutCart } = renderButtons(0);

    await userEvent.click(screen.getByTestId("checkout-button-CASH"));

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

  describe("with SumUp", () => {
    const sumupConfig = () =>
      vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
        currency: { format: (value: number) => `${value} €` },
        paymentMethods: [
          { code: "SUMUP", name: "Card", changeCalculator: false },
        ],
      } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);

    const renderWithReader = (reader: Partial<SumupReaderContextValue>) => {
      sumupConfig();
      const handleCheckoutCart = vi.fn();

      render(
        <SumupReaderContext.Provider
          value={{ ...readerDefaults, enabled: true, ...reader }}
        >
          <CheckoutButtons
            cart={cart(20)}
            checkoutProcessing={null}
            handleCheckoutCart={handleCheckoutCart}
          />
        </SumupReaderContext.Provider>,
      );

      return { handleCheckoutCart };
    };

    it("pays with the selected reader", async () => {
      const { handleCheckoutCart } = renderWithReader({
        selectedReaderId: "rdr_1",
      });

      await userEvent.click(screen.getByTestId("checkout-button-SUMUP"));

      expect(handleCheckoutCart).toHaveBeenCalledWith("SUMUP", {
        type: "sumup",
        sumupReaderId: "rdr_1",
      });
      expect(screen.queryByTestId("sumup-reader-hint")).not.toBeInTheDocument();
    });

    it("is disabled and explains why when the reader has a problem", async () => {
      const openDialog = vi.fn();
      renderWithReader({
        problem: { kind: "none-selected" },
        openDialog,
      });

      expect(screen.getByTestId("checkout-button-SUMUP")).toBeDisabled();

      await userEvent.click(screen.getByTestId("sumup-reader-hint"));
      expect(openDialog).toHaveBeenCalled();
    });
  });
});

const readerDefaults: SumupReaderContextValue = {
  enabled: false,
  selectedReaderId: undefined,
  selectedReader: undefined,
  readers: null,
  loading: false,
  problem: null,
  selectReader: () => undefined,
  refresh: () => undefined,
  dialogOpen: false,
  openDialog: () => undefined,
  closeDialog: () => undefined,
};
