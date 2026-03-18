import React from "react";
import {
  Datagrid,
  TextField,
  NumberField,
  DateField,
  Show,
  SimpleShowLayout,
  ArrayField,
  ShowProps,
} from "react-admin";
import { useConfig } from "../../../../core/config/providers/ConfigProvider";
import SumupTransactionStatusField from "./components/SumupTransactionStatusField";

export const SumupTransactionShow: React.FC<ShowProps> = (props) => {
  const currency = useConfig().currencyOptions;
  const { currencyLocale, dateLocale } = useConfig();

  return (
    <Show {...props}>
      <SimpleShowLayout>
        <TextField source="transactionId" />
        <TextField source="transactionCode" />
        <SimpleShowLayout direction="row">
          <NumberField
            source="amount"
            locales={currencyLocale}
            options={currency}
          />
          <TextField source="currency" />
        </SimpleShowLayout>
        <DateField
          source="createdAt"
          showDate={true}
          showTime={true}
          locales={dateLocale}
          options={{
            weekday: "long",
            year: "numeric",
            month: "long",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
            second: "2-digit",
          }}
        />
        <TextField source="cardType" />
        <SumupTransactionStatusField source="status" label="Status" />
        <ArrayField source="events">
          <Datagrid bulkActionButtons={false}>
            <TextField source="id" />
            <TextField source="type" />
            <TextField source="status" />
            <DateField
              source="timestamp"
              showDate={true}
              showTime={true}
              locales={dateLocale}
              options={{
                weekday: "long",
                year: "numeric",
                month: "long",
                day: "numeric",
                hour: "2-digit",
                minute: "2-digit",
                second: "2-digit",
              }}
            />
          </Datagrid>
        </ArrayField>
      </SimpleShowLayout>
    </Show>
  );
};

export default SumupTransactionShow;
