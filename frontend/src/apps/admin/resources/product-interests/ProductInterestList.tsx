import React from "react";
import {
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  DateField,
} from "react-admin";
import { ProductInterest } from "./types";
import { useConfig } from "../../../../core/config/providers/ConfigProvider";

export const ProductInterestList: React.FC = () => {
  const { locale } = useConfig();

  return (
    <List<ProductInterest> sort={{ field: "pos", order: "ASC" }}>
      <Datagrid rowClick="" bulkActionButtons={false}>
        <NumberField source="id" />
        <NumberField source="wuff" />
        <DateField
          showTime={true}
          locales={locale}
          options={{ weekday: "short", hour: "2-digit", minute: "2-digit" }}
          source="createdAt"
        />
        <NumberField source="product.id" />
        <TextField source="product.name" />
        <DeleteButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};
