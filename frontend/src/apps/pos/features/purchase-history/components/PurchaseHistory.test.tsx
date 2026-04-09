import { describe, it, expect, beforeEach, vi } from "vitest";
import { render, screen, waitFor, act } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import PurchaseHistory from "./PurchaseHistory";
import Decimal from "decimal.js";
import { Purchase } from "../../../utils/api.schemas";

// --- 1. MOCKS ---

vi.mock("@core/logger/logger", () => ({
  createLogger: vi.fn(() => ({
    info: vi.fn(),
    debug: vi.fn(),
    error: vi.fn(),
  })),
}));

vi.mock("@core/config/hooks/useConfig", () => ({
  useConfig: () => ({
    currency: new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }),
    dateLocale: "en-US",
    dateOptions: { dateStyle: "short", timeStyle: "short" },
  }),
}));

const mockShowToast = vi.fn();
vi.mock("@pos/features/ui/toast/hooks/useToast", () => ({
  useToast: () => ({
    showToast: mockShowToast,
  }),
}));

vi.mock("./_internal/RefundModal", () => {
  type Purchase = import("../../../utils/api.schemas").Purchase;

  interface MockRefundModalProps {
    show: boolean;
    purchase: Purchase | null;
    processing: boolean;
    onClose: () => void;
    onConfirm: (purchase: Purchase) => Promise<void> | void;
  }

  return {
    RefundModal: ({
      show,
      onClose,
      onConfirm,
      purchase,
    }: MockRefundModalProps) => {
      if (!show) return null;
      return (
        <div data-testid="mock-refund-modal">
          <button data-testid="modal-close" onClick={onClose}>
            Close
          </button>
          <button
            data-testid="modal-confirm"
            onClick={() => purchase && onConfirm(purchase)}
          >
            Confirm
          </button>
        </div>
      );
    },
  };
});

// --- 2. FIXTURES ---

const mockRemoveFromHistory = vi.fn();
const mockResumePolling = vi.fn();

const createMockPurchase = (
  id: string,
  status: "confirmed" | "pending" = "confirmed",
): Purchase =>
  ({
    id,
    status,
    createdAt: "2026-04-08T12:00:00Z", // Festes Datum für konsistente Tests
    totalGrossPrice: new Decimal(25.5),
  }) as unknown as Purchase;

const defaultProps = {
  history: [],
  loading: false,
  removeFromPurchaseHistory: mockRemoveFromHistory,
  resumePolling: mockResumePolling,
  cartEmpty: true,
};

// --- 3. TESTS ---
describe("PurchaseHistory Component", () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Rendering States", () => {
    it("should render loading spinner when loading is true and no history exists", () => {
      render(<PurchaseHistory {...defaultProps} loading={true} />);
      expect(screen.getByText(/Loading.../i)).toBeInTheDocument();
    });

    it("should render empty state when not loading and history is empty", () => {
      render(<PurchaseHistory {...defaultProps} />);
      expect(screen.getByText("No purchases yet.")).toBeInTheDocument();
    });

    it("should render a maximum of 3 items, even if more are provided", () => {
      const history = [
        createMockPurchase("1"),
        createMockPurchase("2"),
        createMockPurchase("3"),
        createMockPurchase("4"), // should not be rendered
      ];
      render(<PurchaseHistory {...defaultProps} history={history} />);

      const rows = screen.getAllByRole("row");
      // 1 Header-Row + 3 Daten-Rows = 4 Rows total
      expect(rows).toHaveLength(4);
    });
  });

  describe("Resume Flow (Pending Purchases)", () => {
    it("should call resumePolling when resume button is clicked and cart is empty", async () => {
      const pendingPurchase = createMockPurchase("pending-1", "pending");
      render(
        <PurchaseHistory
          {...defaultProps}
          history={[pendingPurchase]}
          cartEmpty={true}
        />,
      );

      const resumeButton = screen.getByLabelText(/Resume purchase/i);
      await user.click(resumeButton);

      expect(mockResumePolling).toHaveBeenCalledWith(pendingPurchase);
      expect(mockShowToast).not.toHaveBeenCalled();
    });

    it("should show warning toast and NOT call resumePolling if cart is NOT empty", async () => {
      const pendingPurchase = createMockPurchase("pending-1", "pending");
      render(
        <PurchaseHistory
          {...defaultProps}
          history={[pendingPurchase]}
          cartEmpty={false}
        />,
      );

      const resumeButton = screen.getByLabelText(/Resume purchase/i);
      await user.click(resumeButton);

      expect(mockResumePolling).not.toHaveBeenCalled();
      expect(mockShowToast).toHaveBeenCalledWith({
        severity: "warning",
        message:
          "Please complete or clear the current cart before resuming this purchase.",
      });
    });
  });

  describe("Refund Flow (Confirmed Purchases)", () => {
    it("should open and close the refund modal", async () => {
      const confirmedPurchase = createMockPurchase("conf-1", "confirmed");
      render(
        <PurchaseHistory {...defaultProps} history={[confirmedPurchase]} />,
      );

      // modal should not be visible initially
      expect(screen.queryByTestId("mock-refund-modal")).not.toBeInTheDocument();

      // click on refund button should open modal
      const refundButton = screen.getByTestId("refund-purchase-conf-1");
      await user.click(refundButton);
      expect(screen.getByTestId("mock-refund-modal")).toBeInTheDocument();

      // click on Close in the modal should close it
      await user.click(screen.getByTestId("modal-close"));
      expect(screen.queryByTestId("mock-refund-modal")).not.toBeInTheDocument();
    });

    it("should process the refund and close modal on confirm", async () => {
      const confirmedPurchase = createMockPurchase("conf-1", "confirmed");
      mockRemoveFromHistory.mockResolvedValueOnce(undefined);

      render(
        <PurchaseHistory {...defaultProps} history={[confirmedPurchase]} />,
      );

      // Open modal
      await user.click(screen.getByTestId("refund-purchase-conf-1"));

      // click confirm in the modal
      await act(async () => {
        await user.click(screen.getByTestId("modal-confirm"));
      });

      expect(mockRemoveFromHistory).toHaveBeenCalledWith(confirmedPurchase);

      // after the api call, the modal should close
      await waitFor(() => {
        expect(
          screen.queryByTestId("mock-refund-modal"),
        ).not.toBeInTheDocument();
      });
    });
  });

  describe("Animation Effects", () => {
    it("should trigger flash animation when a new item is added to the top", () => {
      vi.useFakeTimers();

      const { rerender } = render(
        <PurchaseHistory
          {...defaultProps}
          history={[createMockPurchase("old-1")]}
        />,
      );

      const table = screen.getByTestId("purchase-history-table");
      expect(table.className).not.toContain("animate__pulse");

      // rerender triggers useEffect and schedules requestAnimationFrame
      rerender(
        <PurchaseHistory
          {...defaultProps}
          history={[createMockPurchase("new-2"), createMockPurchase("old-1")]}
        />,
      );

      act(() => {
        vi.advanceTimersByTime(20);
      });

      expect(table.className).toContain("animate__pulse");

      act(() => {
        vi.advanceTimersByTime(500);
      });

      expect(table.className).not.toContain("animate__pulse");

      vi.useRealTimers();
    });
  });
});
