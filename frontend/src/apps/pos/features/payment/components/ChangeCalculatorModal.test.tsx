import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Decimal from "decimal.js";
import ChangeCalculatorModal from "./ChangeCalculatorModal";
import * as ConfigHookModule from "@core/config/hooks/useConfig";

vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: vi.fn(),
}));

const renderModal = (total = 37) => {
  const onConfirm = vi.fn();
  const onClose = vi.fn();

  render(
    <ChangeCalculatorModal
      total={new Decimal(total)}
      paymentMethodName="Cash"
      onConfirm={onConfirm}
      onClose={onClose}
    />,
  );

  return { onConfirm, onClose };
};

describe("ChangeCalculatorModal", () => {
  beforeEach(() => {
    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      currency: { format: (value: number) => `${value} €` },
    } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
  });

  it("starts with an exact payment that can be booked right away", async () => {
    const { onConfirm } = renderModal();

    expect(screen.getByTestId("change-calculator-total")).toHaveTextContent(
      "37 €",
    );
    expect(screen.getByTestId("change-calculator-change")).toHaveTextContent(
      "0 €",
    );

    await userEvent.click(screen.getByTestId("change-calculator-confirm"));

    expect(onConfirm).toHaveBeenCalledTimes(1);
  });

  it("shows the change for a bank note", async () => {
    renderModal();

    await userEvent.click(screen.getByRole("button", { name: "50 €" }));

    expect(screen.getByTestId("change-calculator-change")).toHaveTextContent(
      "13 €",
    );
    expect(screen.getByTestId("change-calculator-confirm")).toBeEnabled();
  });

  it("shows what is missing and blocks booking", () => {
    renderModal(20);

    fireEvent.change(screen.getByLabelText("Other amount received"), {
      target: { value: "12,5" },
    });

    expect(screen.getByTestId("change-calculator-missing")).toHaveTextContent(
      "7.5 €",
    );
    expect(screen.getByTestId("change-calculator-confirm")).toBeDisabled();
  });

  it("asks for a valid amount", () => {
    renderModal(20);

    fireEvent.change(screen.getByLabelText("Other amount received"), {
      target: { value: "abc" },
    });

    expect(screen.getByText("Enter a valid amount")).toBeInTheDocument();
    expect(screen.getByTestId("change-calculator-confirm")).toBeDisabled();
  });

  it("goes back to the exact amount when the input is cleared", async () => {
    renderModal(20);

    const input = screen.getByLabelText("Other amount received");
    fireEvent.change(input, { target: { value: "50" } });
    expect(screen.getByTestId("change-calculator-change")).toHaveTextContent(
      "30 €",
    );

    fireEvent.change(input, { target: { value: "" } });
    expect(screen.getByTestId("change-calculator-change")).toHaveTextContent(
      "0 €",
    );

    await userEvent.click(screen.getByRole("button", { name: "Exact" }));
    expect(screen.getByTestId("change-calculator-change")).toHaveTextContent(
      "0 €",
    );
  });

  it("can be cancelled", async () => {
    const { onClose, onConfirm } = renderModal();

    await userEvent.click(screen.getByTestId("change-calculator-cancel"));

    expect(onClose).toHaveBeenCalled();
    expect(onConfirm).not.toHaveBeenCalled();
  });
});
