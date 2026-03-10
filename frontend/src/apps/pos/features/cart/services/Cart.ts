import Decimal from "decimal.js";
import { CartItem, PaymentMethodData } from "../types/cart.types";
import { ApiCreatePayloadPurchase } from "../../../utils/api.types";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../utils/api.schemas";

export class Cart {
  public readonly items: readonly CartItem[];

  constructor(items: CartItem[] = []) {
    this.items = Object.freeze([...items]);
  }

  public add(
    product: ProductType,
    count: number = 1,
    listItem: GuestType | null = null,
  ): Cart {
    const existingIndex = this.items.findIndex(
      (item) => item.id === product.id,
    );
    let newItems = [...this.items];

    if (existingIndex !== -1) {
      const existingItem = this.items[existingIndex];

      // Dubletten-Check für ListItems
      if (
        listItem &&
        existingItem.listItems.find((li) => li.id === listItem.id)
      ) {
        return this;
      }

      // Immutable Update des Items
      const updatedItem: CartItem = {
        ...existingItem,
        quantity: existingItem.quantity + count,
        listItems: listItem
          ? [...existingItem.listItems, { ...listItem, attendedGuests: count }]
          : existingItem.listItems,
        totalNetPrice: existingItem.netPrice.mul(existingItem.quantity + count),
        totalGrossPrice: existingItem.grossPrice.mul(
          existingItem.quantity + count,
        ),
        totalVatAmount: existingItem.vatAmount.mul(
          existingItem.quantity + count,
        ),
      };
      newItems[existingIndex] = updatedItem;
    } else {
      // Neues Item erstellen
      const newItem: CartItem = {
        ...product,
        quantity: count,
        listItems: listItem ? [{ ...listItem, attendedGuests: count }] : [],
        totalNetPrice: product.netPrice.mul(count),
        totalGrossPrice: product.grossPrice.mul(count),
        totalVatAmount: product.vatAmount.mul(count),
      };
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
    return {
      paymentMethod: paymentMethodCode,
      cart: this.items.map((item) => ({
        ...item,
        lists: null,
        guestlists: null,
      })),
      totalGrossPrice: this.totalGross.toString(),
      totalNetPrice: this.totalNet.toString(),
      ...paymentMethodData,
    };
  }
}
