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
  FunctionField,
} from "react-admin";
import InventoryIcon from "@mui/icons-material/Inventory";
import { useConfig } from "../../provider/ConfigProvider";
import { PurchaseListToolbar } from "./PurchaseListToolbar";

export const PurchaseList = () => {
  const {
    currencyOptions: currency,
    Locale: locale,
    paymentMethods,
  } = useConfig();

  return (
    <List sort={{ field: "createdAt", order: "DESC" }} actions={false}>
      <PurchaseListToolbar />
      <Datagrid rowClick="show" bulkActionButtons={false}>
        <FunctionField
          source="id"
          label="ID"
          render={(record) =>
            record.id ? (
              <span title={record.id}>
                {record.id.slice(0, 4)}…{record.id.slice(-4)}
              </span>
            ) : null
          }
        />
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
        <FunctionField
          source="paymentMethod"
          render={(record) => {
            if (!paymentMethods) {
              return record.paymentMethod;
            }
            if (Array.isArray(paymentMethods)) {
              const paymentMethod = paymentMethods.find(
                (pm) => pm.code === record.paymentMethod,
              );
              if (paymentMethod) {
                return paymentMethod.name;
              }
            }

            return record.paymentMethod;
          }}
        />
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
