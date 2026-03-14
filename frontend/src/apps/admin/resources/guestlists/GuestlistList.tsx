import React from "react";
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  BooleanField,
  ListProps,
} from "react-admin";
import ConditionalDeleteOnOwnershipButton from "@admin/ui/buttons/ConditionalDeleteOnOwnershipButton";
import CreateGuestlistEntryButton from "./components/CreateGuestlistEntryButton";

export const GuestlistList: React.FC<ListProps> = (props) => {
  return (
    <List {...props} sort={{ field: "id", order: "ASC" }}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <BooleanField source="typeCode" sortable={false} label="Code" />
        <TextField source="product.name" sortable={false} label="Product" />
        <ConditionalDeleteOnOwnershipButton mutationMode="pessimistic" />
        <CreateGuestlistEntryButton />
      </Datagrid>
    </List>
  );
};

export default GuestlistList;
