import React from "react";
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  DateField,
  FunctionField,
} from "react-admin";
import { useConfig } from "../../../../core/config/providers/ConfigProvider";
import { PurchaseListToolbar } from "./components/PurchaseListToolbar";
import { PurchaseStatusField } from "./components/PurchaseStatusField";
import { PaymentMethodField } from "./components/PaymentMethodField";
import { ConditionalDeleteOnAdminButton } from "../../ui/buttons/ConditionalDeleteOnAdminButton";
import { ConditionalRefundButton } from "./components/ConditionalRefundButton";

export const PurchaseList: React.FC = () => {
  const { currencyOptions, locale } = useConfig();

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
          options={currencyOptions}
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
