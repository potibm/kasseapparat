import React from "react";
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  DateField,
  ListProps,
} from "react-admin";
import { useConfig } from "../../../../core/config/providers/ConfigProvider";
import { SumupTransactionTimeRangeFilter } from "./components/SumupTransactionTimeRangeFilter";
import { SumupTransactionStatusField } from "./components/SumupTransactionStatusField";

export const SumupTransactionList: React.FC<ListProps> = (props) => {
  const {
    sumupEnabled,
    currencyOptions: currency,
    currencyLocale,
    dateLocale,
  } = useConfig();

  if (!sumupEnabled) {
    return (
      <div>
        <h2>SumUp Transactions</h2>
        <p>
          SumUp integration is not enabled. Please enable it in the
          configuration.
        </p>
      </div>
    );
  }

  return (
    <List
      title="SumUp Transactions"
      filters={<SumupTransactionTimeRangeFilter />}
      empty={false}
      {...props}
    >
      <Datagrid rowClick="show" bulkActionButtons={false}>
        <TextField source="id" sortable={false} />
        <TextField source="transactionCode" sortable={false} />
        <NumberField
          source="amount"
          locales={currencyLocale}
          options={currency}
          sortable={false}
        />
        <TextField source="currency" sortable={false} />
        <DateField
          source="createdAt"
          showTime={true}
          locales={dateLocale}
          options={{ weekday: "short", hour: "2-digit", minute: "2-digit" }}
          sortable={false}
        />
        <SumupTransactionStatusField source="status" />
      </Datagrid>
    </List>
  );
};

export default SumupTransactionList;
