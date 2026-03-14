import React from "react";
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  BooleanField,
  EmailField,
  ListProps,
} from "react-admin";
import { UserFilters } from "./components/UserFilters";
import { ConditionalDeleteButton } from "./components/ConditionalDeleteButton";

export const UserList: React.FC<ListProps> = (props) => {
  return (
    <List sort={{ field: "id", order: "ASC" }} filters={UserFilters} {...props}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="username" />
        <EmailField source="email" />
        <BooleanField source="admin" />
        <ConditionalDeleteButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};
