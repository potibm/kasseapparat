import React from "react";
import {
  usePermissions,
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  BooleanField,
  ReferenceField,
} from "react-admin";
import { useConfig } from "@core/config/hooks/useConfig";
import { ProductFilters } from "./ProductFilters";
import { Product } from "./types";

export const ProductList: React.FC = () => {
  const { currencyOptions, currencyLocale } = useConfig();
  const { permissions, isLoading: permissionsLoading } = usePermissions();

  if (permissionsLoading) return <>Loading...</>;

  return (
    <List<Product>
      sort={{ field: "pos", order: "ASC" }}
      filters={ProductFilters}
    >
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <NumberField source="vatRate" />
        <NumberField
          source="grossPrice"
          locales={currencyLocale}
          options={currencyOptions}
        />
        <ReferenceField
          source="sectionId"
          reference="sections"
          label="Section"
          link={false}
        />
        <ReferenceField
          source="productGroupId"
          reference="productGroups"
          label="Group"
          link={false}
          emptyText="–"
        />
        <TextField source="variant" sortable={false} emptyText="–" />
        <NumberField source="pos" />
        <BooleanField source="soldOut" sortable={false} />
        <BooleanField source="hidden" sortable={false} />
        {permissions === "admin" && <DeleteButton mutationMode="pessimistic" />}
      </Datagrid>
    </List>
  );
};
