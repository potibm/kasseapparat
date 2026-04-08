import { describe, it, expect } from "vitest";
import {
  getPurchaseErrorType,
  getErrorMessage,
  PurchaseErrorType,
} from "./PurchaseErrorHandler";

describe("PurchaseErrorHandler", () => {
  describe("getPurchaseErrorType()", () => {
    it("should return ReaderBusy when error message contains 'Reader Busy'", () => {
      const error = new Error("The terminal returned: Reader Busy right now.");
      const type = getPurchaseErrorType(error);
      expect(type).toBe(PurchaseErrorType.ReaderBusy);
    });

    it("should return Generic for errors without specific keywords", () => {
      const error = new Error("Network timeout");
      const type = getPurchaseErrorType(error);
      expect(type).toBe(PurchaseErrorType.Generic);
    });

    it("should return Generic for non-Error objects", () => {
      const type = getPurchaseErrorType({ some: "object" });
      expect(type).toBe(PurchaseErrorType.Generic);
    });
  });

  describe("getErrorMessage()", () => {
    it("should return correct message for ReaderBusy", () => {
      const message = getErrorMessage(PurchaseErrorType.ReaderBusy);
      expect(message).toBe(
        "The SumUp reader is currently busy. Please complete or cancel the ongoing transaction.",
      );
    });

    it("should return default message for Generic", () => {
      const message = getErrorMessage(PurchaseErrorType.Generic);
      expect(message).toBe("An error occurred while processing the purchase.");
    });
  });
});
