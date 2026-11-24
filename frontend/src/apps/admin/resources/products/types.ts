import { RaRecord } from "react-admin";

export interface Product extends RaRecord {
  id: number;
  name: string;
  vatRate: number;
  netPrice: number;
  grossPrice: number;
  pos: number;
  sectionId: number;
  productGroupId?: number | null;
  variant?: string;
  soldOut: boolean;
  hidden: boolean;
  totalStock?: number;
  unitsSold?: number;
  soldOutRequestCount?: number;
  apiExport?: boolean;
}
