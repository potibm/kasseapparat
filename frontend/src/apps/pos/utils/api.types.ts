import { Guestlist } from "../features/product-list/types/product.types";

interface User {
  id: number;
  name: string;
}

export interface ApiGetResponseProduct {
  id: number;
  name: string;
  netPrice: string;
  grossPrice: string;
  vatRate: string;
  vatAmount: string;
  wrapAfter: boolean;
  hidden: boolean;
  soldOut: boolean;
  apiExport: boolean;
  pos: number;
  totalStock: number;
  guestlists?: Guestlist[];
  unitsSold: number;
  soldOutRequestCount: number;
}

export interface ApiCreateResponsePurchase {
  id: string;
  createdAt: string;
  createdById?: number;
  createdBy?: User;
  paymentMethod: string;
  totalNetPrice: string;
  sumupTransactionId?: string;
  sumupClientTransactionId?: string;
  totalGrossPrice: string;
  totalVatAmount: string;
  purchaseItems: ApiGetResponsePurchase[];
  status: string;
}

export interface ApiGetResponsePurchase {
  id: number;
  purchaseID: string;
  productID: number;
  product: ApiGetResponseProduct;
  quantity: number;
  netPrice: string;
  grossPrice: string;
  vatRate: string;
  vatAmount: string;
  totalNetPrice: string;
  totalGrossPrice: string;
  totalVatAmount: string;
}

export interface ApiCreateResponseProductInterest {
  id: number;
  productID: number;
}
