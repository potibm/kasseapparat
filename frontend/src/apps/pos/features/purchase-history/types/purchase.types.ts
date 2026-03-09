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
  createdBy: User | null;
  sumupTransactionId: string | null;
  sumupClientTransactionId: string | null;
}

interface User {
  id: number;
  username: string;
  email: string;
  admin: boolean;
}

interface PurchaseItem {
  id: number;
  purchaseID: string;
  productID: number;
  product: Product;
  quantity: number;
  netPrice: Decimal;
  grossPrice: Decimal;
  vatRate: Decimal;
  vatAmount: Decimal;
  totalNetPrice: Decimal;
  totalGrossPrice: Decimal;
  totalVatAmount: Decimal;
}
