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
  const { dateLocale } = useConfig();

  return (
    <List<ProductInterest> sort={{ field: "createdAt", order: "DESC" }}>
      <Datagrid rowClick="" bulkActionButtons={false}>
        <NumberField source="id" />
        <DateField
          showTime={true}
          locales={dateLocale}
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
