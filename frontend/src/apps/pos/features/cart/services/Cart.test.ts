import { describe, it, expect, vi, beforeEach } from "vitest";
import { Cart } from "./Cart";
import Decimal from "decimal.js";
import { PaymentMethodData } from "../types/cart.types";
import { createLogger } from "@core/logger/logger";
import {
  createMockProduct,
  createMockGuest,
} from "@pos/utils/api.schemas.mocks";

vi.mock("@core/logger/logger", () => {
  const mockLogger = {
    warn: vi.fn(),
    debug: vi.fn(),
    error: vi.fn(),
  };
  return {
    createLogger: vi.fn(() => mockLogger),
  };
});

describe("Cart", () => {
  const mockLogger = createLogger("Cart");

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Initializing", () => {
    it("should initialize with an empty cart", () => {
      const cart = new Cart();
      expect(cart.isEmpty).toBe(true);
      expect(cart.items).toEqual([]);
      expect(cart.totalQuantity).toBe(0);
    });

    it("should freeze the items array (immutable)", () => {
      const cart = new Cart();
      expect(Object.isFrozen(cart.items)).toBe(true);
    });
  });

  describe("add()", () => {
    it("should add a new product to the cart", () => {
      const cart = new Cart();
      const product = createMockProduct({ id: 1, netPrice: new Decimal(100) });

      const updatedCart = cart.add(product, 2);

      // the original cart should still be empty
      expect(cart.isEmpty).toBe(true);

      // the new cart should have the product with correct quantity and prices
      expect(updatedCart.items).toHaveLength(1);
      expect(updatedCart.items[0].id).toBe(1);
      expect(updatedCart.items[0].quantity).toBe(2);
      expect(updatedCart.items[0].totalNetPrice).toEqual(new Decimal(200));

      expect(mockLogger.debug).toHaveBeenCalledWith(
        "Adding new product to cart",
        expect.any(Object),
      );
    });

    it("should increase the quantity and prices when the product is already in the cart", () => {
      const product = createMockProduct({ id: 1, netPrice: new Decimal(100) });
      let cart = new Cart().add(product, 1);

      // Add again
      cart = cart.add(product, 2);

      expect(cart.items).toHaveLength(1); // still one entry
      expect(cart.items[0].quantity).toBe(3); // 1 + 2
      expect(cart.items[0].totalNetPrice).toEqual(new Decimal(300));

      expect(mockLogger.debug).toHaveBeenCalledWith(
        "Product already in cart, updating quantity",
        expect.any(Object),
      );
    });

    it("should block and warn when an invalid quantity is provided", () => {
      const cart = new Cart();
      const product = createMockProduct({ id: 1, netPrice: new Decimal(100) });

      const invalidCounts = [0, -5, 1.5, Infinity, Number.NaN];

      invalidCounts.forEach((count) => {
        const unchangedCart = cart.add(product, count);

        expect(unchangedCart).toBe(cart);
        expect(mockLogger.warn).toHaveBeenCalledWith(
          "Invalid quantity provided, skipping add to cart",
          expect.any(Object),
        );
      });
    });
  });

  describe("add() - list items (guests) ", () => {
    it("should add a guest to a new product", () => {
      const cart = new Cart();
      const product = createMockProduct({ id: 1, netPrice: new Decimal(100) });
      const guest = createMockGuest({ id: 99 });

      const newCart = cart.add(product, 1, guest);

      expect(newCart.items[0].listItems).toHaveLength(1);
      expect(newCart.items[0].listItems[0].id).toBe(99);
      expect(newCart.items[0].listItems[0].attendedGuests).toBe(1);
    });

    it("should warn and abort when the guest is already assigned to the product", () => {
      const product = createMockProduct({ id: 1, netPrice: new Decimal(100) });
      const guest = createMockGuest({ id: 99 });

      const cart = new Cart().add(product, 1, guest);
      // Try to assign the same guest again
      const unchangedCart = cart.add(product, 1, guest);

      expect(unchangedCart).toBe(cart); // Reference equality due to early return
      expect(mockLogger.warn).toHaveBeenCalledWith(
        "Guest was already in cart",
        expect.any(Object),
      );
    });
  });

  describe("remove()", () => {
    it("should remove a product by ID", () => {
      const product1 = createMockProduct({ id: 1, netPrice: new Decimal(100) });
      const product2 = createMockProduct({ id: 2, netPrice: new Decimal(200) });
      const cart = new Cart().add(product1).add(product2);

      const newCart = cart.remove(1);

      expect(newCart.items).toHaveLength(1);
      expect(newCart.items[0].id).toBe(2);
    });
  });

  describe("Getters und Helfer (totalGross, totalNet, getQuantity, etc.)", () => {
    it("should calculate the correct totals for all items", () => {
      const cart = new Cart()
        .add(createMockProduct({ id: 1, netPrice: new Decimal(100) }), 2) // 200 net
        .add(createMockProduct({ id: 2, netPrice: new Decimal(50) }), 1); // 50 net

      expect(cart.totalNet).toEqual(new Decimal(250));
      expect(cart.totalQuantity).toBe(3);
      expect(cart.getQuantity(1)).toBe(2);
      expect(cart.getQuantity(999)).toBe(0); // Does not exist
    });

    it("should be able to check if a ListItem (guest) exists in the entire cart", () => {
      const cart = new Cart().add(
        createMockProduct({ id: 1, netPrice: new Decimal(100) }),
        1,
        createMockGuest({ id: 42 }),
      );

      expect(cart.hasListItem(42)).toBe(true);
      expect(cart.hasListItem(99)).toBe(false);
    });
  });

  describe("toApiPayload()", () => {
    it("should format the payload correctly for the API", () => {
      const cart = new Cart().add(
        createMockProduct({ id: 1, netPrice: new Decimal(100) }),
        1,
      );
      const paymentData = {
        type: "CREDIT_CARD", // should be removed in the payload
        token: "tok_123", // should be included in the payload
      };

      const payload = cart.toApiPayload("cc", paymentData as PaymentMethodData);

      expect(payload).toEqual(
        expect.objectContaining({
          paymentMethod: "cc",
          token: "tok_123",
          totalNetPrice: "100", // Decimal.toString()
        }),
      );

      // type should be removed from the payload
      expect(payload).not.toHaveProperty("type");

      // the items should have formatted fields (lists: null, guestlists: null)
      expect(payload.cart[0].lists).toBeNull();
      expect(payload.cart[0].guestlists).toBeNull();
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
});
