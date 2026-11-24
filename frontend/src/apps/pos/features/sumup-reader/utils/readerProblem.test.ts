import { describe, expect, it } from "vitest";
import { describeReaderProblem, findReaderProblem } from "./readerProblem";

const bar = { id: "rdr_1", name: "Bar 1", status: "paired" };

describe("findReaderProblem", () => {
  it("needs a selected reader", () => {
    const problem = findReaderProblem(undefined, [bar]);

    expect(problem).toEqual({ kind: "none-selected" });
    expect(describeReaderProblem(problem!)).toBe("No card reader selected");
  });

  it("reports a selected reader that is no longer paired", () => {
    const problem = findReaderProblem("rdr_gone", [bar]);

    expect(problem).toEqual({ kind: "missing", readerId: "rdr_gone" });
    expect(describeReaderProblem(problem!)).toBe(
      "The selected card reader is no longer paired",
    );
  });

  it("reports a reader that cannot take payments", () => {
    const expired = { ...bar, status: "expired" };
    const problem = findReaderProblem("rdr_1", [expired]);

    expect(problem).toEqual({ kind: "unusable", reader: expired });
    expect(describeReaderProblem(problem!)).toBe(
      "Card reader „Bar 1“ is expired",
    );
  });

  it("is happy with a paired or busy reader", () => {
    expect(findReaderProblem("rdr_1", [bar])).toBeNull();
    expect(
      findReaderProblem("rdr_1", [{ ...bar, status: "processing" }]),
    ).toBeNull();
  });

  it("does not guess while the readers are unknown", () => {
    expect(findReaderProblem("rdr_1", null)).toBeNull();
  });
});
