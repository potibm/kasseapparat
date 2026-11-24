import { beforeEach, describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import GuestlistModal from "./GuestlistModal";
import * as PosApiModule from "@pos/api/usePosApi";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import { createMockProduct } from "@pos/api/schemas.mocks";
import { Guest } from "@pos/api/schemas";

vi.mock("@pos/api/usePosApi", () => ({ usePosApi: vi.fn() }));
vi.mock("@core/config/hooks/useConfig", () => ({ useConfig: vi.fn() }));

const guests: Guest[] = [
  {
    id: 1,
    name: "Alex Müller",
    code: null,
    listName: "Crew",
    additionalGuests: 0,
    arrivalNote: null,
    attendedGuests: 0,
  },
  {
    id: 2,
    name: "Tim Mülheim",
    code: null,
    listName: "Press",
    additionalGuests: 0,
    arrivalNote: "Hand over the press badge",
    attendedGuests: 0,
  },
];

const product = createMockProduct({
  id: 3,
  name: "Free",
  totalStock: 0,
  unitsSold: 0,
  soldOut: false,
  guestlists: [{ id: 1, name: "Crew", typeCode: false, productId: 3 }],
});

const fetchGuests = vi.fn();

const renderModal = () => {
  const onClose = vi.fn();
  const addToCart = vi.fn();

  render(
    <GuestlistModal
      isOpen
      onClose={onClose}
      product={product}
      addToCart={addToCart}
      hasListItem={() => false}
    />,
  );

  return { onClose, addToCart };
};

describe("GuestlistModal", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    fetchGuests.mockResolvedValue(guests);
    vi.mocked(PosApiModule.usePosApi).mockReturnValue({
      fetchGuestlistByProductId: fetchGuests,
    } as unknown as ReturnType<typeof PosApiModule.usePosApi>);
    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      outOfStockBehavior: "ignore",
    } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
  });

  it("lists the guests and searches with the letter keys", async () => {
    renderModal();

    expect(await screen.findByText("Alex Müller")).toBeInTheDocument();
    expect(screen.getByText("List for Free")).toBeInTheDocument();
    expect(screen.getByLabelText("Search name")).toHaveValue("");

    fireEvent.click(
      screen.getByRole("button", { name: "Add M to search term" }),
    );
    fireEvent.click(
      screen.getByRole("button", { name: "Add U to search term" }),
    );

    expect(screen.getByTestId("guestlist-search-input")).toHaveValue("MU");
    await waitFor(() =>
      expect(fetchGuests).toHaveBeenLastCalledWith(3, "MU", expect.anything()),
    );

    fireEvent.click(screen.getByRole("button", { name: "Clear search term" }));
    expect(screen.getByTestId("guestlist-search-input")).toHaveValue("");
  });

  it("checks in a guest and closes", async () => {
    const { addToCart, onClose } = renderModal();

    fireEvent.click(
      await screen.findByRole("button", { name: "Add Alex Müller to cart" }),
    );

    expect(addToCart).toHaveBeenCalledWith(product, 1, guests[0]);
    expect(onClose).toHaveBeenCalled();
  });

  it("shows the arrival note in the guest's card before closing", async () => {
    const { addToCart, onClose } = renderModal();

    fireEvent.click(
      await screen.findByRole("button", { name: "Add Tim Mülheim to cart" }),
    );

    expect(addToCart).toHaveBeenCalledWith(product, 1, guests[1]);
    expect(onClose).not.toHaveBeenCalled();
    expect(screen.getByTestId("guestlist-arrival-note")).toHaveTextContent(
      "Hand over the press badge",
    );

    fireEvent.click(screen.getByRole("button", { name: "Got it" }));
    expect(onClose).toHaveBeenCalled();
  });

  it("sells without a list entry", async () => {
    const { addToCart, onClose } = renderModal();
    await screen.findByText("Alex Müller");

    fireEvent.click(
      screen.getByRole("button", {
        name: "Add product to cart without guestlist entry",
      }),
    );

    expect(addToCart).toHaveBeenCalledWith(product, 1, null);
    expect(onClose).toHaveBeenCalled();
  });

  it("closes with the close button", async () => {
    const { onClose } = renderModal();
    await screen.findByText("Alex Müller");

    fireEvent.click(screen.getByRole("button", { name: "Close guest list" }));

    expect(onClose).toHaveBeenCalled();
  });

  it("reports a failed search", async () => {
    fetchGuests.mockRejectedValue(new Error("offline"));
    renderModal();

    expect(
      await screen.findByText("Error fetching list entries: offline"),
    ).toBeInTheDocument();
  });
});
