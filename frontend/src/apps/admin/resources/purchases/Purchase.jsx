import React from "react";
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  DateField,
  Show,
  SimpleShowLayout,
  ArrayField,
  FunctionField,
  ReferenceField,
} from "react-admin";
import { useConfig } from "../../../../core/config/providers/ConfigProvider";
import { PurchaseListToolbar } from "./components/PurchaseListToolbar";
import { PurchaseStatusField } from "./components/PurchaseStatusField";
import { PaymentMethodField } from "./components/PaymentMethodField";
import { ConditionalDeleteOnAdminButton } from "../../ui/buttons/ConditionalDeleteOnAdminButton";
import { ConditionalRefundButton } from "./components/ConditionalRefundButton";

export const PurchaseList = () => {
  const { currencyOptions: currency, Locale: locale } = useConfig();

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
        <PaymentMethodField source="paymentMethod" />
        <PurchaseStatusField source="status" label="Status" />
        <ConditionalDeleteOnAdminButton mutationMode="pessimistic" />
        <ConditionalRefundButton />
      </Datagrid>
    </List>
  );
};

export const PurchaseShow = (props) => {
  const { currencyOptions: currency, locale } = useConfig();

  return (
    <Show {...props}>
      <SimpleShowLayout>
        <TextField source="id" />
        <DateField source="createdAt" showTime={true} />
        <PurchaseStatusField source="status" label="Status" />
        <PaymentMethodField source="paymentMethod" />
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
        <ReferenceField
          label="SumUp Transaction"
          source="sumupTransactionId"
          reference="sumupTransactions"
          link="show"
        >
          <TextField source="id" />
        </ReferenceField>
        <TextField source="sumupClientTransactionId" />
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
