import Decimal from "decimal.js";
import { Product } from "../../product-list/types/product.types";
import { Guest } from "../../guestlist/types/guest.types";

export interface CartItem extends Product {
  quantity: number;
  listItems: Guest[];
  totalNetPrice: Decimal;
  totalGrossPrice: Decimal;
  totalVatAmount: Decimal;
}

interface EmptyPaymentMethodData {}

interface SumUpPaymentMethodData {
  sumupReaderId: string;
}

export type PaymentMethodData = EmptyPaymentMethodData | SumUpPaymentMethodData;
