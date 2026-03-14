import { RaRecord } from "react-admin";

export interface UserRecord extends RaRecord {
  id: number;
  username: string;
  email: string;
  admin: boolean;
}
