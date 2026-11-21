import { RaRecord } from "react-admin";

export interface Product extends RaRecord {
  id: number;
  name: string;
  vatRate: number;
  netPrice: number;
  grossPrice: number;
  pos: number;
  sectionId: number;
  soldOut: boolean;
  hidden: boolean;
  totalStock?: number;
  unitsSold?: number;
  soldOutRequestCount?: number;
  apiExport?: boolean;
}
