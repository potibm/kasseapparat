import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { RecordContextProvider } from "react-admin";
import GuestlistProgressField from "./GuestlistProgressField";

const renderField = (record: object | undefined) =>
  render(
    <RecordContextProvider value={record}>
      <GuestlistProgressField />
    </RecordContextProvider>,
  );

describe("GuestlistProgressField", () => {
  it("compares the arrived entries with all entries and adds the people", () => {
    renderField({ id: 1, entries: 20, arrivedEntries: 2, arrivedGuests: 3 });

    expect(screen.getByTestId("guestlist-progress")).toHaveTextContent(
      "2 / 20 arrived · 3 people",
    );
    expect(
      screen.getByRole("progressbar", { name: "2 of 20 entries arrived" }),
    ).toHaveAttribute("aria-valuenow", "10");
  });

  it("uses the singular for one person", () => {
    renderField({ id: 1, entries: 1, arrivedEntries: 1, arrivedGuests: 1 });

    expect(screen.getByTestId("guestlist-progress")).toHaveTextContent(
      "1 / 1 arrived · 1 person",
    );
  });

  it("says when a list is empty", () => {
    renderField({ id: 1 });

    expect(screen.getByText("No entries")).toBeInTheDocument();
  });

  it("renders nothing without a record", () => {
    const { container } = renderField(undefined);

    expect(container).toBeEmptyDOMElement();
  });
});
