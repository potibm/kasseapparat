import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { PollingModal } from "./PollingModal";
import * as PaymentWSModule from "../hooks/usePaymentWebSocket";
import * as ToastHookModule from "@pos/features/ui/toast/hooks/useToast";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import * as LocalStorageModule from "@core/localstorage/helper/local-storage-reader";
import "@testing-library/jest-dom";
import { Purchase as PurchaseType } from "../../../utils/api.schemas";
import { AppConfig } from "@core/config/types/config.types";
import { UsePaymentWebSocketReturn } from "../types/payment.types";

// --- Mocks Setup ---
vi.mock("../hooks/usePaymentWebSocket");
vi.mock("@pos/features/ui/toast/hooks/useToast");
vi.mock("@core/config/hooks/useConfig");
vi.mock("@core/localstorage/helper/local-storage-reader");

describe("PollingModal", () => {
  const mockOnComplete = vi.fn();
  const mockShowToast = vi.fn();
  const mockCancel = vi.fn();

  // A representative purchase object for testing
  const mockPurchase = {
    id: "purchase-123",
    paymentMethod: "SUMUP",
    totalGrossPrice: { toNumber: () => 40 },
  } as unknown as PurchaseType;

  beforeEach(() => {
    vi.resetAllMocks();

    vi.mocked(ToastHookModule.useToast).mockReturnValue({
      showToast: mockShowToast,
    });

    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      currency: { format: (val: number) => `€${val.toFixed(2)}` },
      paymentMethods: [{ code: "SUMUP", name: "SumUp" }],
    } as unknown as AppConfig);

    vi.mocked(LocalStorageModule.getCurrentReaderId).mockReturnValue(
      "reader-xyz",
    );
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders pending state correctly and allows cancellation", () => {
    // WebSocket sends "pending" status
    vi.mocked(PaymentWSModule.usePaymentWebSocket).mockReturnValue({
      status: "pending",
      error: null,
      lastMessageAt: Date.now(),
      cancel: mockCancel,
      isConnected: true,
    } as unknown as UsePaymentWebSocketReturn);

    render(
      <PollingModal purchase={mockPurchase} onComplete={mockOnComplete} />,
    );

    // UI Checks
    expect(screen.getByText(/Waiting for Terminal.../i)).toBeInTheDocument();

    // Click cancel button
    const abortBtn = screen.getByRole("button", { name: /Abort Purchase/i });
    expect(abortBtn).not.toBeDisabled();

    fireEvent.click(abortBtn);

    // check that cancel function was called with the correct reader ID
    expect(mockCancel).toHaveBeenCalledWith("reader-xyz");
    expect(screen.getByText(/Cancelling.../i)).toBeInTheDocument();
  });

  it("handles a successful payment (confirmed) correctly", () => {
    // WebSocket is "confirmed" directly
    vi.mocked(PaymentWSModule.usePaymentWebSocket).mockReturnValue({
      status: "confirmed",
      error: null,
      lastMessageAt: Date.now(),
      cancel: mockCancel,
      isConnected: true,
    } as unknown as UsePaymentWebSocketReturn);

    render(
      <PollingModal purchase={mockPurchase} onComplete={mockOnComplete} />,
    );

    // Toast needs to be called with success message
    expect(mockShowToast).toHaveBeenCalledWith({
      severity: "success",
      message: "Payment of €40.00 successful via SumUp!",
    });

    // onComplete needs to be called with success
    expect(mockOnComplete).toHaveBeenCalledWith(true);
  });

  it("handles a user cancellation correctly", () => {
    // WebSocket sends "cancelled"
    vi.mocked(PaymentWSModule.usePaymentWebSocket).mockReturnValue({
      status: "cancelled",
      error: null,
      lastMessageAt: Date.now(),
      cancel: mockCancel,
      isConnected: true,
    } as unknown as UsePaymentWebSocketReturn);

    render(
      <PollingModal purchase={mockPurchase} onComplete={mockOnComplete} />,
    );

    // Error toast for cancellation
    expect(mockShowToast).toHaveBeenCalledWith(
      expect.objectContaining({
        severity: "error",
        message: "Payment via SumUp was cancelled by the user.",
        autoClose: true,
        blocking: false,
      }),
    );

    // onComplete muss mit false (failure) aufgerufen werden
    expect(mockOnComplete).toHaveBeenCalledWith(false);
  });

  it("handles a terminal timeout correctly", () => {
    // WebSocket sends "timeout"
    vi.mocked(PaymentWSModule.usePaymentWebSocket).mockReturnValue({
      status: "timeout",
      error: null,
      lastMessageAt: Date.now(),
      cancel: mockCancel,
      isConnected: true,
    } as unknown as UsePaymentWebSocketReturn);

    render(
      <PollingModal purchase={mockPurchase} onComplete={mockOnComplete} />,
    );

    // Error toast for timeout
    expect(mockShowToast).toHaveBeenCalledWith(
      expect.objectContaining({
        severity: "error",
        message: "Timeout: No response from SumUp terminal.",
        autoClose: false, // Real error. needs to be acknowledged by the user.
        blocking: true,
      }),
    );

    expect(mockOnComplete).toHaveBeenCalledWith(false);
  });
});
