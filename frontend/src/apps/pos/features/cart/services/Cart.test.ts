import { describe, it, expect } from "vitest";
import { Cart } from "./Cart";
import Decimal from "decimal.js";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../utils/api.schemas";
import { PaymentMethodData } from "../types/cart.types";

describe("Cart.toApiPayload", () => {
  it("should remove irrelevant fields and keep only necessary payment data", () => {
    const mockProduct = {
      id: 1,
      netPrice: new Decimal("32.00"),
      grossPrice: new Decimal("40.00"),
      vatAmount: new Decimal("8.00"),
      guestlists: [
        {
          id: 2,
          name: "Guest List 1",
          typeCode: true,
          productId: 1,
        },
      ],
    };

    const guest: GuestType = {
      id: 1,
      name: "John Doe",
      code: null,
      listName: "Guest List 1",
      additionalGuests: 0,
      arrivalNote: null,
      attendedGuests: 0,
    };

    const cart = new Cart().add(mockProduct as ProductType, 1, guest);

    const paymentData = {
      type: "empty",
    };

    const payload = cart.toApiPayload(
      "sumup_terminal",
      paymentData as PaymentMethodData,
    );

    expect(payload.totalGrossPrice).toBe("40");
    expect(payload.totalNetPrice).toBe("32");
    expect(payload.cart).toHaveLength(1);
    expect(payload.cart[0]).toMatchObject({
      id: 1,
      netPrice: new Decimal("32.00"),
      grossPrice: new Decimal("40.00"),
      vatAmount: new Decimal("8.00"),
      quantity: 1,
      listItems: [
        {
          id: 1,
          name: "John Doe",
          code: null,
          listName: "Guest List 1",
          additionalGuests: 0,
          arrivalNote: null,
          attendedGuests: 1,
        },
      ],
    });
  });

  it("should remove the 'type' discriminator and keep only relevant payment data", () => {
    const mockProduct = {
      id: 1,
      netPrice: new Decimal("10.00"),
      grossPrice: new Decimal("11.90"),
      vatAmount: new Decimal("1.90"),
    };

    const cart = new Cart().add(mockProduct as ProductType, 1);

    const paymentData = {
      type: "sumup", // this field should be removed in the payload
      sumupReaderId: "ABC-123",
      someOtherField: "value",
    };

    const payload = cart.toApiPayload(
      "sumup_terminal",
      paymentData as PaymentMethodData,
    );

    expect(payload).not.toHaveProperty("type");

    expect(payload).toMatchObject({
      paymentMethod: "sumup_terminal",
      sumupReaderId: "ABC-123",
      someOtherField: "value",
    });

    expect(payload.totalGrossPrice).toBe("11.9");
    expect(payload.totalNetPrice).toBe("10");
  });

  it("should handle empty payment data gracefully", () => {
    const cart = new Cart();
    const payload = cart.toApiPayload("cash", {
      type: "empty",
    } as PaymentMethodData);

    expect(payload).toEqual(
      expect.objectContaining({
        paymentMethod: "cash",
        totalGrossPrice: "0",
      }),
    );
    expect(payload).not.toHaveProperty("type");
  });
});
