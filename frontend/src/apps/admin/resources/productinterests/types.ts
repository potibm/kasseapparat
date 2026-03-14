import { RaRecord } from "react-admin";
import { Product as ProductType } from "../products/types.ts";

export interface ProductInterest extends RaRecord {
  id: number;
  createdAt: string;
  product: ProductType;
}
