import React from "react";
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  BooleanField,
  ListProps,
  SimpleList,
} from "react-admin";
import { useIsSmallScreen } from "@admin/ui/hooks/useIsSmallScreen";
import ConditionalDeleteOnOwnershipButton from "@admin/ui/buttons/ConditionalDeleteOnOwnershipButton";
import CreateGuestlistEntryButton from "./components/CreateGuestlistEntryButton";
import { GuestlistFilters } from "./GuestlistFilters";
import GuestlistProgressField from "./components/GuestlistProgressField";

export const GuestlistList: React.FC<ListProps> = (props) => {
  const isSmall = useIsSmallScreen();

  return (
    <List
      {...props}
      filters={GuestlistFilters}
      sort={{ field: "id", order: "ASC" }}
    >
      {isSmall ? (
        <SimpleList
          linkType="edit"
          primaryText={(list) => list.name}
          secondaryText={(list) => list.product?.name ?? ""}
          tertiaryText={(list) =>
            list.entries
              ? `${list.arrivedEntries ?? 0} / ${list.entries} arrived`
              : "No entries"
          }
        />
      ) : (
        <Datagrid rowClick="edit" bulkActionButtons={false}>
          <NumberField source="id" />
          <TextField source="name" />
          <BooleanField source="typeCode" sortable={false} label="Code" />
          <TextField source="product.name" sortable={false} label="Product" />
          <GuestlistProgressField label="Arrived" />
          <CreateGuestlistEntryButton />
          <ConditionalDeleteOnOwnershipButton mutationMode="pessimistic" />
        </Datagrid>
      )}
    </List>
  );
};

export default GuestlistList;
