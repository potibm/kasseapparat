import { beforeEach, describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import GuestlistModal from "./GuestlistModal";
import * as PosApiModule from "@pos/api/usePosApi";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import { createMockProduct } from "@pos/api/schemas.mocks";
import { Guest } from "@pos/api/schemas";

vi.mock("@pos/api/usePosApi", () => ({ usePosApi: vi.fn() }));
vi.mock("@core/config/hooks/useConfig", () => ({ useConfig: vi.fn() }));

const ticket = (overrides: Partial<Guest>): Guest => ({
  id: 1,
  name: "Jan Jansen",
  code: "ABCDEFGHI",
  listName: "Deine Tickets",
  additionalGuests: 0,
  arrivalNote: null,
  attendedGuests: 0,
  ...overrides,
});

const product = createMockProduct({
  id: 4,
  name: "Prepaid",
  totalStock: 0,
  unitsSold: 0,
  soldOut: false,
  guestlists: [{ id: 3, name: "Deine Tickets", typeCode: true, productId: 4 }],
});

const fetchGuests = vi.fn();

const renderModal = (inCart: (guest: Guest) => boolean = () => false) => {
  const onClose = vi.fn();
  const addToCart = vi.fn();

  render(
    <GuestlistModal
      isOpen
      onClose={onClose}
      product={product}
      addToCart={addToCart}
      hasListItem={inCart}
    />,
  );

  return { onClose, addToCart };
};

const scan = (code: string) => {
  const input = screen.getByTestId("guestlist-search-input");
  fireEvent.change(input, { target: { value: code } });
  fireEvent.keyDown(input, { key: "Enter" });
};

const returnFor = (code: string, guest: Guest) =>
  fetchGuests.mockImplementation(async (_id: number, query: string) =>
    query === code ? [guest] : [],
  );

describe("GuestlistModal with a barcode scanner", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    fetchGuests.mockResolvedValue([]);
    vi.mocked(PosApiModule.usePosApi).mockReturnValue({
      fetchGuestlistByProductId: fetchGuests,
    } as unknown as ReturnType<typeof PosApiModule.usePosApi>);
    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      outOfStockBehavior: "ignore",
    } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
  });

  it("shows when the scanner is ready", async () => {
    renderModal();

    // the search field gets the focus for code lists
    expect(await screen.findByText("Ready to scan")).toBeInTheDocument();

    fireEvent.blur(screen.getByTestId("guestlist-search-input"));
    const notReady = screen.getByRole("button", {
      name: /Not ready to scan/,
    });

    fireEvent.click(notReady);
    expect(screen.getByTestId("guestlist-search-input")).toHaveFocus();
  });

  it("checks in a scanned ticket and stays open for the next one", async () => {
    const guest = ticket({});
    returnFor("ABCDEFGHI", guest);
    const { addToCart, onClose } = renderModal();

    scan("ABCDEFGHI");

    await waitFor(() =>
      expect(addToCart).toHaveBeenCalledWith(product, 1, guest),
    );
    expect(onClose).not.toHaveBeenCalled();
    expect(screen.getByTestId("guestlist-scan-status")).toHaveTextContent(
      "ABCDEFGHI added – scan the next ticket",
    );
    expect(screen.getByTestId("guestlist-search-input")).toHaveValue("");
  });

  it("leaves the number of guests to the cashier", async () => {
    returnFor("ABCDEFGHI", ticket({ additionalGuests: 2 }));
    const { addToCart } = renderModal();

    scan("ABCDEFGHI");

    expect(
      await screen.findByText(/choose the number of guests/),
    ).toBeVisible();
    expect(addToCart).not.toHaveBeenCalled();
  });

  it("does not add a ticket twice", async () => {
    returnFor("ABCDEFGHI", ticket({}));
    const { addToCart } = renderModal(() => true);

    scan("ABCDEFGHI");

    expect(
      await screen.findByText("ABCDEFGHI is already in the cart"),
    ).toBeInTheDocument();
    expect(addToCart).not.toHaveBeenCalled();
  });

  it("refuses a ticket when none are left", async () => {
    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      outOfStockBehavior: "fail",
    } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
    returnFor("ABCDEFGHI", ticket({}));
    const { addToCart } = renderSoldOutModal();

    scan("ABCDEFGHI");

    expect(await screen.findByText("ABCDEFGHI: no tickets left")).toBeVisible();
    expect(addToCart).not.toHaveBeenCalled();
  });

  it("shows the arrival note of a scanned ticket and continues scanning", async () => {
    returnFor("ABCDEFGHI", ticket({ arrivalNote: "VIP wristband" }));
    const { addToCart, onClose } = renderModal();

    scan("ABCDEFGHI");

    expect(
      await screen.findByTestId("guestlist-arrival-note"),
    ).toHaveTextContent("VIP wristband");
    expect(addToCart).toHaveBeenCalled();

    fireEvent.click(screen.getByRole("button", { name: "Got it" }));

    expect(onClose).not.toHaveBeenCalled();
    expect(screen.queryByTestId("guestlist-arrival-note")).toBeNull();
    expect(screen.getByTestId("guestlist-search-input")).toHaveValue("");
  });
});

/** A product whose stock is used up: 2 of 2 sold. */
function renderSoldOutModal() {
  const onClose = vi.fn();
  const addToCart = vi.fn();
  const soldOut = { ...product, totalStock: 2, unitsSold: 2 };

  render(
    <GuestlistModal
      isOpen
      onClose={onClose}
      product={soldOut}
      addToCart={addToCart}
      hasListItem={() => false}
    />,
  );

  return { onClose, addToCart };
}
