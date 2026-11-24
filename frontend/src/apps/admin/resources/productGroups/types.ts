import { RaRecord } from "react-admin";

export interface ProductGroup extends RaRecord {
  id: number;
  sectionId: number;
  name: string;
  pos: number;
}
