import Decimal from "decimal.js";

export interface Product {
  id: number;
  name: string;
  netPrice: Decimal;
  grossPrice: Decimal;
  vatRate: Decimal;
  vatAmount: Decimal;
  wrapAfter: boolean;
  hidden: boolean;
  soldOut: boolean;
  apiExport: boolean;
  pos: number;
  totalStock: number;
  guestlists?: Guestlist[] | null;
  unitsSold: number;
  soldOutRequestCount: number;
}

export interface Guestlist {
  id: number;
  name: string;
  typeCode: boolean;
  productId: number;
}
