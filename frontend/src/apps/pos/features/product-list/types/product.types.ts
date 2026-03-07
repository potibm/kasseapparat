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
  guestlists: any[];
  unitsSold: number;
  soldOutRequestCount: number;
}
