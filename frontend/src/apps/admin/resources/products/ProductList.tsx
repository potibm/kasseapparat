import React from "react";
import {
  usePermissions,
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  BooleanField,
} from "react-admin";
import { useConfig } from "@core/config/hooks/useConfig";
import { Product } from "./types";

export const ProductList: React.FC = () => {
  const { currencyOptions, currencyLocale } = useConfig();
  const { permissions, isLoading: permissionsLoading } = usePermissions();

  if (permissionsLoading) return <>Loading...</>;

  return (
    <List<Product> sort={{ field: "pos", order: "ASC" }}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <NumberField source="vatRate" />
        <NumberField
          source="grossPrice"
          locales={currencyLocale}
          options={currencyOptions}
        />
        <NumberField source="pos" />
        <BooleanField source="wrapAfter" sortable={false} />
        <BooleanField source="soldOut" sortable={false} />
        <BooleanField source="hidden" sortable={false} />
        {permissions === "admin" && <DeleteButton mutationMode="pessimistic" />}
      </Datagrid>
    </List>
  );
};
