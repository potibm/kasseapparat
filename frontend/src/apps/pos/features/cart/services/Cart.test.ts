import { describe, it, expect, vi, beforeEach } from "vitest";
import { Cart } from "./Cart";
import Decimal from "decimal.js";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../utils/api.schemas";
import { PaymentMethodData } from "../types/cart.types";
import { createLogger } from "@core/logger/logger";

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

const createMockProduct = (id: number, price: number) =>
  ({
    id,
    netPrice: new Decimal(price),
    grossPrice: new Decimal(price * 1.19), // Beispiel 19% MwSt.
    vatAmount: new Decimal(price * 0.19),
  }) as any;

const createMockGuest = (id: number) =>
  ({
    id,
    name: `Guest ${id}`,
  }) as any;

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
      const product = createMockProduct(1, 100);

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
      const product = createMockProduct(1, 100);
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
      const product = createMockProduct(1, 100);

      // Wir prüfen verschiedene ungültige Inputs
      const invalidCounts = [0, -5, 1.5, Infinity, NaN];

      invalidCounts.forEach((count) => {
        const unchangedCart = cart.add(product, count);

        // Die Cart-Instanz sollte die gleiche bleiben
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
      const product = createMockProduct(1, 100);
      const guest = createMockGuest(99);

      const newCart = cart.add(product, 1, guest);

      expect(newCart.items[0].listItems).toHaveLength(1);
      expect(newCart.items[0].listItems[0].id).toBe(99);
      expect(newCart.items[0].listItems[0].attendedGuests).toBe(1);
    });

    it("should warn and abort when the guest is already assigned to the product", () => {
      const product = createMockProduct(1, 100);
      const guest = createMockGuest(99);

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
      const product1 = createMockProduct(1, 100);
      const product2 = createMockProduct(2, 200);
      const cart = new Cart().add(product1).add(product2);

      const newCart = cart.remove(1);

      expect(newCart.items).toHaveLength(1);
      expect(newCart.items[0].id).toBe(2);
    });
  });

  describe("Getters und Helfer (totalGross, totalNet, getQuantity, etc.)", () => {
    it("should calculate the correct totals for all items", () => {
      const cart = new Cart()
        .add(createMockProduct(1, 100), 2) // 200 Netto
        .add(createMockProduct(2, 50), 1); // 50 Netto

      expect(cart.totalNet).toEqual(new Decimal(250));
      expect(cart.totalQuantity).toBe(3);
      expect(cart.getQuantity(1)).toBe(2);
      expect(cart.getQuantity(999)).toBe(0); // Existiert nicht
    });

    it("should be able to check if a ListItem (guest) exists in the entire cart", () => {
      const cart = new Cart().add(
        createMockProduct(1, 100),
        1,
        createMockGuest(42),
      );

      expect(cart.hasListItem(42)).toBe(true);
      expect(cart.hasListItem(99)).toBe(false);
    });
  });

  describe("toApiPayload()", () => {
    it("should format the payload correctly for the API", () => {
      const cart = new Cart().add(createMockProduct(1, 100));
      const paymentData = {
        type: "CREDIT_CARD", // Sollte gefiltert werden
        token: "tok_123", // Sollte bleiben
      };

      const payload = cart.toApiPayload("cc", paymentData as any);

      expect(payload).toEqual(
        expect.objectContaining({
          paymentMethod: "cc",
          token: "tok_123",
          totalNetPrice: "100", // Decimal.toString()
        }),
      );

      // type should be removed from the payload
      expect((payload as any).type).toBeUndefined();

      // the items should have formatted fields (lists: null, guestlists: null)
      expect(payload.cart[0].lists).toBeNull();
      expect(payload.cart[0].guestlists).toBeNull();
    });
  });
});

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
