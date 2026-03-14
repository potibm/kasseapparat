import { SearchInput } from "react-admin";
import { QuickFilter } from "@admin/ui/filters/QuickFilter";

export const UserFilters = [
  <SearchInput source="q" alwaysOn key="ID" />,
  <QuickFilter source="isAdmin" label="Admin" defaultValue={true} key="ID" />,
];

export default UserFilters;
