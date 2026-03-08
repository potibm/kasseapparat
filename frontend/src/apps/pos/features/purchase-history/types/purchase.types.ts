import Decimal from "decimal.js";
import { Product } from "../../product-list/types/product.types";

export interface Purchase {
  id: string;
  status: "pending" | "confirmed" | "refunded" | "failed";
  paymentMethod: string;
  totalGrossPrice: Decimal;
  totalNetPrice: Decimal;
  totalVatAmount: Decimal;
  purchaseItems: PurchaseItem[];
  createdAt: string;
  createdById: number;
  createdBy: User;
  sumupTransactionId?: string;
  sumupClientTransactionId?: string;
}

interface User {
  id: number;
  username: string;
  email: string;
  admin: boolean;
}

interface PurchaseItem {
  id: string;
  purchaseID: string;
  productID: number;
  product: Product;
  quantity: number;
  netPrice: Decimal;
  grossPrice: Decimal;
  vatRate: number;
  vatAmount: Decimal;
  totalNetPrice: Decimal;
  totalGrossPrice: Decimal;
  totalVatAmount: Decimal;
}
