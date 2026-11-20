import { describe, expect, it } from "vitest";
import { HttpError } from "react-admin";
import {
  SECTION_NOT_EMPTY_MESSAGE,
  UNNAMED_SECTION_LABEL,
  sectionDeleteErrorMessage,
  sectionLabel,
} from "./section.utils";

describe("sectionLabel", () => {
  it("uses the name of a named section", () => {
    expect(sectionLabel({ id: 1, name: "Tickets" })).toBe("Tickets");
  });

  it("labels an unnamed section as a row break", () => {
    expect(sectionLabel({ id: 3, name: "" })).toBe(
      `${UNNAMED_SECTION_LABEL} #3`,
    );
    expect(sectionLabel({ id: 4, name: "   " })).toBe(
      `${UNNAMED_SECTION_LABEL} #4`,
    );
  });

  it("returns an empty string without a section", () => {
    expect(sectionLabel(undefined)).toBe("");
    expect(sectionLabel(null)).toBe("");
  });
});

describe("sectionDeleteErrorMessage", () => {
  it("explains a conflict as a section that still has products", () => {
    const error = new HttpError("Conflict", 409, { error: "Conflict" });

    expect(sectionDeleteErrorMessage(error)).toBe(SECTION_NOT_EMPTY_MESSAGE);
  });

  it("prefers the details the backend sent", () => {
    const error = new HttpError("Forbidden", 403, {
      error: "Forbidden",
      details: "Only admins can delete sections",
    });

    expect(sectionDeleteErrorMessage(error)).toBe(
      "Only admins can delete sections",
    );
  });

  it("falls back to the error message", () => {
    expect(sectionDeleteErrorMessage(new HttpError("Not Found", 404))).toBe(
      "Not Found",
    );
    expect(sectionDeleteErrorMessage(new Error("Network down"))).toBe(
      "Network down",
    );
  });

  it("has a generic message for anything else", () => {
    expect(sectionDeleteErrorMessage("boom")).toBe(
      "The section could not be deleted.",
    );
  });
});
