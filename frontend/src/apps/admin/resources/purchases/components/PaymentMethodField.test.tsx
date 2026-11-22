import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { RecordContextProvider } from "react-admin";
import { PaymentMethodField } from "./PaymentMethodField";
import * as ConfigHookModule from "@core/config/hooks/useConfig";

vi.mock("@core/config/hooks/useConfig", () => ({ useConfig: vi.fn() }));

const renderField = (paymentMethod: string | undefined) => {
  vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
    paymentMethods: [{ code: "CASH", name: "Cash" }],
  } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);

  return render(
    <RecordContextProvider value={{ id: "p1", paymentMethod }}>
      <PaymentMethodField source="paymentMethod" />
    </RecordContextProvider>,
  );
};

describe("PaymentMethodField", () => {
  it("shows the configured name with the method's colour", () => {
    renderField("CASH");

    expect(screen.getByText("Cash")).toBeInTheDocument();
    expect(screen.getByTestId("payment-method-color")).toHaveStyle({
      backgroundColor: "#4caf50",
    });
  });

  it("falls back to the code for an unknown method", () => {
    renderField("GIFT");

    expect(screen.getByText("GIFT")).toBeInTheDocument();
  });

  it("renders nothing without a payment method", () => {
    const { container } = renderField(undefined);

    expect(container).toBeEmptyDOMElement();
  });
});
