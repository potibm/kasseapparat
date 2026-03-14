import Decimal from "decimal.js";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../utils/api.schemas";

export interface CartItem extends ProductType {
  quantity: number;
  listItems: GuestType[];
  totalNetPrice: Decimal;
  totalGrossPrice: Decimal;
  totalVatAmount: Decimal;
}

interface EmptyPaymentMethodData {}

interface SumUpPaymentMethodData {
  sumupReaderId: string;
}

export type PaymentMethodData = EmptyPaymentMethodData | SumUpPaymentMethodData;
