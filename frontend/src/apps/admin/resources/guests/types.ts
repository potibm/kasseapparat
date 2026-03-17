import { RaRecord } from "react-admin";

export interface GuestRecord extends RaRecord {
  id: number;
  guestlistId: number;
  name: string;
}
