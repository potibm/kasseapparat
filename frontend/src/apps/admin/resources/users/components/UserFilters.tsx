import { SearchInput } from "react-admin";
import { QuickFilter } from "@admin/ui/filters/QuickFilter";

const UserFilters = [
  <SearchInput source="q" alwaysOn key="search" />,
  <QuickFilter
    source="isAdmin"
    label="Admin"
    defaultValue={true}
    key="admin"
  />,
];

export default UserFilters;
