import { describe, expect, it } from "vitest";
import { foldText, splitByMatches } from "./searchHighlight";

describe("foldText", () => {
  it("drops case and diacritics", () => {
    expect(foldText("MÜLLER")).toBe("muller");
    expect(foldText("Zoë Çelik")).toBe("zoe celik");
    expect(foldText("Straße")).toBe("strasse");
    expect(foldText("Søren Łukasz")).toBe("soren lukasz");
  });
});

describe("splitByMatches", () => {
  it("highlights a match regardless of case and diacritics", () => {
    expect(splitByMatches("Alex Müller", "MUL")).toEqual([
      { text: "Alex ", match: false },
      { text: "Mül", match: true },
      { text: "ler", match: false },
    ]);
  });

  it("maps matches across letters that fold to two letters", () => {
    expect(splitByMatches("Hans Straße", "strasse")).toEqual([
      { text: "Hans ", match: false },
      { text: "Straße", match: true },
    ]);
  });

  it("highlights every occurrence", () => {
    expect(splitByMatches("Anna Ann", "ann")).toEqual([
      { text: "Ann", match: true },
      { text: "a ", match: false },
      { text: "Ann", match: true },
    ]);
  });

  it("returns the whole text without a query", () => {
    expect(splitByMatches("Alex", "  ")).toEqual([
      { text: "Alex", match: false },
    ]);
    expect(splitByMatches("", "a")).toEqual([]);
  });
});
