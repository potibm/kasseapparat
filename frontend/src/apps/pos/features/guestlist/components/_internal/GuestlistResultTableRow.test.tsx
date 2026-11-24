import { describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import GuestlistResultTableRow from "./GuestlistResultTableRow";
import { Guest } from "@pos/api/schemas";

const guest = (overrides: Partial<Guest> = {}): Guest => ({
  id: 1,
  name: "Alex Müller",
  code: null,
  listName: "Crew",
  additionalGuests: 2,
  arrivalNote: null,
  attendedGuests: 0,
  ...overrides,
});

const renderRow = (
  entry: Guest,
  {
    inCart = false,
    showArrivalNote = false,
    query = "",
  }: { inCart?: boolean; showArrivalNote?: boolean; query?: string } = {},
) => {
  const onAddToCart = vi.fn();
  const onAcknowledgeNote = vi.fn();

  render(
    <GuestlistResultTableRow
      entry={entry}
      onAddToCart={onAddToCart}
      hasListItem={() => inCart}
      loadedSearchQuery={query}
      showArrivalNote={showArrivalNote}
      onAcknowledgeNote={onAcknowledgeNote}
    />,
  );

  return { onAddToCart, onAcknowledgeNote };
};

describe("GuestlistResultTableRow", () => {
  it("shows the guest with list and companions and checks them in", () => {
    const entry = guest();
    const { onAddToCart } = renderRow(entry, { query: "MUL" });

    const row = screen.getByTestId("guestlist-result-id-1");
    expect(row).toHaveTextContent("Alex Müller");
    expect(row).toHaveTextContent("Crew · up to 2 additional guests");
    expect(screen.getByText("Mül")).toHaveClass("underline");

    fireEvent.click(
      screen.getByRole("button", { name: "Add Alex Müller to cart" }),
    );
    expect(onAddToCart).toHaveBeenLastCalledWith(entry, 0);

    fireEvent.click(
      screen.getByRole("button", {
        name: "Add Alex Müller with 2 additional guest(s) to cart",
      }),
    );
    expect(onAddToCart).toHaveBeenLastCalledWith(entry, 2);
  });

  it("uses the singular for one companion and shows a code instead of a name", () => {
    renderRow(guest({ code: "ABCDEFGHI", additionalGuests: 1 }));

    expect(screen.getByTestId("guestlist-result-id-1")).toHaveTextContent(
      "Crew · up to 1 additional guest",
    );
    expect(
      screen.getByRole("button", { name: "Add ABCDEFGHI to cart" }),
    ).toBeInTheDocument();
  });

  it("marks a guest already in the cart", () => {
    renderRow(guest(), { inCart: true });

    expect(screen.getByText("In cart")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /to cart/ })).toBeNull();
  });

  it("flags an arrival note and shows it in the card after the check-in", () => {
    const entry = guest({ arrivalNote: "Hand over the press badge" });
    renderRow(entry);

    expect(screen.getByText("Arrival note")).toBeInTheDocument();
    expect(screen.queryByTestId("guestlist-arrival-note")).toBeNull();
  });

  it("asks to acknowledge the arrival note", () => {
    const { onAcknowledgeNote } = renderRow(
      guest({ arrivalNote: "Hand over the press badge" }),
      { inCart: true, showArrivalNote: true },
    );

    const note = screen.getByTestId("guestlist-arrival-note");
    expect(note).toHaveTextContent("Hand over the press badge");
    expect(screen.getByText("Added")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Got it" }));
    expect(onAcknowledgeNote).toHaveBeenCalled();
  });
});
