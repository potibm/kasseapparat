import React from "react";
import {
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  DateField,
  Show,
  SimpleShowLayout,
  ArrayField,
  ReferenceInput,
  AutocompleteInput,
  SelectArrayInput,
  NumberInput,
  Filter,
} from "react-admin";
import InventoryIcon from "@mui/icons-material/Inventory";
import { useConfig } from "../../provider/ConfigProvider";
import { PurchaseFilters } from "./PurchaseFilters";

export const PurchaseList = () => {
  const { currencyOptions: currency, Locale: locale } = useConfig();

  return (
    <List
      filters={<PurchaseFilters />}
      sort={{ field: "createdAt", order: "DESC" }}
    >
      <Datagrid rowClick="show" bulkActionButtons={false}>
        <NumberField source="id" />
        <DateField
          source="createdAt"
          showTime={true}
          locales={locale}
          options={{ weekday: "short", hour: "2-digit", minute: "2-digit" }}
        />
        <NumberField
          source="totalGrossPrice"
          locales={locale}
          options={currency}
        />
        <TextField source="createdBy.username" />
        <TextField source="paymentMethod" />
        <DeleteButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};

export const PurchaseShow = (props) => {
  const currency = useConfig().currencyOptions;
  const locale = useConfig().Locale;

  return (
    <Show {...props}>
      <SimpleShowLayout>
        <NumberField source="id" />
        <DateField source="createdAt" showTime={true} />
        <TextField source="paymentMethod" />
        <NumberField
          source="totalNetPrice"
          locales={locale}
          options={currency}
        />
        <NumberField
          source="totalVatAmount"
          locales={locale}
          options={currency}
        />
        <NumberField
          source="totalGrossPrice"
          locales={locale}
          options={currency}
        />
        <ArrayField source="purchaseItems">
          <Datagrid bulkActionButtons={false}>
            <NumberField source="quantity" />
            <TextField source="product.name" />
            <NumberField
              source="totalNetPrice"
              locales={locale}
              options={currency}
            />
            <NumberField
              source="totalGrossPrice"
              locales={locale}
              options={currency}
            />
          </Datagrid>
        </ArrayField>
      </SimpleShowLayout>
    </Show>
  );
};

export const ProductIcon = () => <InventoryIcon />;
