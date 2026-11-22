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
  it("shows the arrived of the expected guests", () => {
    renderField({ id: 1, entries: 2, expectedGuests: 5, arrivedGuests: 2 });

    expect(screen.getByText("2 / 5 arrived")).toBeInTheDocument();
    expect(screen.getByText("2 entries")).toBeInTheDocument();
    expect(
      screen.getByRole("progressbar", { name: "2 of 5 guests arrived" }),
    ).toHaveAttribute("aria-valuenow", "40");
  });

  it("uses the singular for one entry", () => {
    renderField({ id: 1, entries: 1, expectedGuests: 1, arrivedGuests: 0 });

    expect(screen.getByText("1 entry")).toBeInTheDocument();
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
