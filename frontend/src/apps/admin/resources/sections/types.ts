import { RaRecord } from "react-admin";

export interface Section extends RaRecord {
  id: number;
  name: string;
  pos: number;
}
