import Decimal from "decimal.js";
import { CartItem, PaymentMethodData } from "../types/cart.types";
import { ApiCreatePayloadPurchase } from "../../../api/types";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../api/schemas";
import { createLogger } from "@core/logger/logger";
import { OutOfStockBehavior } from "@core/config/types/config.types";
import {
  availableStock,
  canAddMore,
  restrictsStock,
} from "@pos/features/stock/utils/stock";

const log = createLogger("Cart");

export class Cart {
  public readonly items: readonly CartItem[];

  constructor(items: CartItem[] = []) {
    this.items = Object.freeze([...items]);
  }

  public add(
    product: ProductType,
    count: number = 1,
    listItem: GuestType | null = null,
    outOfStockBehavior: OutOfStockBehavior = "ignore",
  ): Cart {
    if (count <= 0 || !Number.isFinite(count) || !Number.isInteger(count)) {
      log.warn("Invalid quantity provided, skipping add to cart", {
        productId: product.id,
        count,
      });
      return this;
    }

    const existingIndex = this.items.findIndex(
      (item) => item.id === product.id,
    );
    const newItems = [...this.items];
    const itemProductWasFoundInCart = existingIndex !== -1;
    const existingQuantity = itemProductWasFoundInCart
      ? this.items[existingIndex].quantity
      : 0;

    // Cap the quantity at the stock that is left when the behaviour enforces it. Under
    // "ignore" the sale stays allowed, so the quantity passes through untouched and
    // the cart only reports that it is over stock. The cap is the product's total
    // availability, so it bounds the final quantity rather than the increment.
    const enforcing = restrictsStock(outOfStockBehavior);
    const capacity = enforcing
      ? availableStock(product, 0)
      : Number.POSITIVE_INFINITY;
    const quantity = Math.min(existingQuantity + count, capacity);

    // The same predicate the product card uses, so the add button and the cart can
    // never disagree about whether another unit fits.
    if (!canAddMore(product, existingQuantity, outOfStockBehavior)) {
      log.warn("No stock left to add", {
        productId: product.id,
        existingQuantity,
        available: availableStock(product, existingQuantity),
      });

      return this;
    }

    const addedQuantity = quantity - existingQuantity;

    if (itemProductWasFoundInCart) {
      const existingItem = this.items[existingIndex];

      // Duplicate prevention for list items (guests)
      if (
        listItem &&
        existingItem.listItems.some((li) => li.id === listItem.id)
      ) {
        log.warn("Guest was already in cart", {
          productId: product.id,
          listItemId: listItem.id,
        });
        return this;
      }

      // Immutable update of the item
      const updatedItem: CartItem = {
        ...existingItem,
        quantity,
        listItems: listItem
          ? [
              ...existingItem.listItems,
              { ...listItem, attendedGuests: addedQuantity },
            ]
          : existingItem.listItems,
        totalNetPrice: existingItem.netPrice.mul(quantity),
        totalGrossPrice: existingItem.grossPrice.mul(quantity),
        totalVatAmount: existingItem.vatAmount.mul(quantity),
      };

      log.debug("Product already in cart, updating quantity", {
        productId: product.id,
        existingQuantity,
        requestedQuantity: existingQuantity + count,
        addedQuantity,
        clamped: addedQuantity !== count,
      });
      newItems[existingIndex] = updatedItem;
    } else {
      // Create new item
      const newItem: CartItem = {
        ...product,
        quantity,
        listItems: listItem
          ? [{ ...listItem, attendedGuests: addedQuantity }]
          : [],
        totalNetPrice: product.netPrice.mul(quantity),
        totalGrossPrice: product.grossPrice.mul(quantity),
        totalVatAmount: product.vatAmount.mul(quantity),
      };

      log.debug("Adding new product to cart", {
        productId: product.id,
        quantity,
        requestedQuantity: count,
        clamped: quantity !== count,
      });
      newItems.push(newItem);
    }

    return new Cart(newItems);
  }

  public remove(productId: number): Cart {
    return new Cart(this.items.filter((item) => item.id !== productId));
  }

  public get totalGross(): Decimal {
    return this.items.reduce(
      (sum, item) => sum.plus(item.totalGrossPrice),
      new Decimal(0),
    );
  }

  public get isEmpty(): boolean {
    return this.items.length === 0;
  }

  public getQuantity(productId: number): number {
    return this.items.find((i) => i.id === productId)?.quantity ?? 0;
  }

  /**
   * Whether a line exceeds the stock available, which the POS allows under the
   * "ignore" behaviour but should point out.
   */
  public isOverStock(item: CartItem): boolean {
    return availableStock(item, 0) < item.quantity;
  }

  public get totalQuantity(): number {
    return this.items.reduce((total, item) => total + item.quantity, 0);
  }

  public get totalNet(): Decimal {
    return this.items.reduce(
      (sum, item) => sum.plus(item.totalNetPrice),
      new Decimal(0),
    );
  }

  public hasListItem(listItemId: number): boolean {
    return this.items.some((product) =>
      product.listItems.some((listItem) => listItem.id === listItemId),
    );
  }

  public toApiPayload(
    paymentMethodCode: string,
    paymentMethodData: PaymentMethodData,
  ): ApiCreatePayloadPurchase {
    const { type: _type, ...cleanPaymentData } = paymentMethodData;

    return {
      paymentMethod: paymentMethodCode,
      cart: this.items.map((item) => ({
        ...item,
        lists: null,
        guestlists: null,
      })),
      totalGrossPrice: this.totalGross.toString(),
      totalNetPrice: this.totalNet.toString(),
      ...cleanPaymentData,
    };
  }
}
