import { RaRecord } from "react-admin";

export interface Product extends RaRecord {
  id: number;
  name: string;
  vatRate: number;
  netPrice: number;
  grossPrice: number;
  pos: number;
  wrapAfter: boolean;
  soldOut: boolean;
  hidden: boolean;
  totalStock?: number;
  unitsSold?: number;
  soldOutRequestCount?: number;
  apiExport?: boolean;
}
