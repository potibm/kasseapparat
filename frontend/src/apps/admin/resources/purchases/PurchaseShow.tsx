import React from "react";
import {
  Datagrid,
  TextField,
  NumberField,
  DateField,
  Show,
  SimpleShowLayout,
  ArrayField,
  ReferenceField,
} from "react-admin";
import { useConfig } from "../../../../core/config/providers/ConfigProvider";
import { PurchaseStatusField } from "./components/PurchaseStatusField";
import { PaymentMethodField } from "./components/PaymentMethodField";

export const PurchaseShow: React.FC = (props) => {
  const { currencyOptions: currency, currencyLocale } = useConfig();

  return (
    <Show {...props}>
      <SimpleShowLayout>
        <TextField source="id" />
        <DateField source="createdAt" showTime={true} />
        <PurchaseStatusField source="status" label="Status" />
        <PaymentMethodField source="paymentMethod" />
        <NumberField
          source="totalNetPrice"
          locales={currencyLocale}
          options={currency}
        />
        <NumberField
          source="totalVatAmount"
          locales={currencyLocale}
          options={currency}
        />
        <NumberField
          source="totalGrossPrice"
          locales={currencyLocale}
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
              locales={currencyLocale}
              options={currency}
            />
            <NumberField
              source="totalGrossPrice"
              locales={currencyLocale}
              options={currency}
            />
          </Datagrid>
        </ArrayField>
      </SimpleShowLayout>
    </Show>
  );
};
