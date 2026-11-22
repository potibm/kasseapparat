import { describe, expect, it } from "vitest";
import { findScannedGuest } from "./scan";
import { Guest } from "@pos/api/schemas";

const guest = (id: number, code: string | null): Guest => ({
  id,
  name: "Guest " + id,
  code,
  listName: "Deine Tickets",
  additionalGuests: 0,
  arrivalNote: null,
  attendedGuests: 0,
});

describe("findScannedGuest", () => {
  it("finds the guest whose code equals the scan", () => {
    const entries = [guest(1, "ABCDEFGHI"), guest(2, "ABCDEFGHIJ")];

    expect(findScannedGuest(entries, " abcdefghi ")).toBe(entries[0]);
  });

  it("ignores prefixes, names and empty scans", () => {
    const entries = [guest(1, "ABCDEFGHI"), guest(2, null)];

    expect(findScannedGuest(entries, "ABCDEF")).toBeNull();
    expect(findScannedGuest(entries, "Guest 2")).toBeNull();
    expect(findScannedGuest(entries, "")).toBeNull();
  });

  it("does not guess between two guests with the same code", () => {
    expect(
      findScannedGuest([guest(1, "SAME"), guest(2, "SAME")], "SAME"),
    ).toBeNull();
  });
});
