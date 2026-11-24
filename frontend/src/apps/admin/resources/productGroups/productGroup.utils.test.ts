import { describe, expect, it } from "vitest";
import { isInGroup, parseProductGroupId } from "./productGroup.utils";

describe("parseProductGroupId", () => {
  it("sends 0 for no group", () => {
    expect(parseProductGroupId("")).toBe(0);
    expect(parseProductGroupId(null)).toBe(0);
    expect(parseProductGroupId(undefined)).toBe(0);
    expect(parseProductGroupId("nope")).toBe(0);
  });

  it("keeps a group id", () => {
    expect(parseProductGroupId(7)).toBe(7);
    expect(parseProductGroupId("7")).toBe(7);
  });
});

describe("isInGroup", () => {
  it("is true only for a real group id", () => {
    expect(isInGroup(3)).toBe(true);
    expect(isInGroup(0)).toBe(false);
    expect(isInGroup(null)).toBe(false);
  });
});
