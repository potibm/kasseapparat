import Decimal from "decimal.js";
import { Product } from "../../product-list/types/product.types";

export interface CartItem extends Product {
  quantity: number;
  listItems: ListItem[];
  totalNetPrice: Decimal;
  totalGrossPrice: Decimal;
  totalVatAmount: Decimal;
}

export interface ListItem {
  id: number;
  name: string;
  attendedGuests?: number;
}
