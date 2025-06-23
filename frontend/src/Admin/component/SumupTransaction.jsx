import React from "react";
import {
  List,
  Datagrid,
  TextField,
  FunctionField,
  NumberField,
  DateField,
  Show,
  SimpleShowLayout,
  Filter,
  SelectInput,
  ArrayField,
} from "react-admin";
import { useConfig } from "../../provider/ConfigProvider";
import { Box } from "@mui/material";

import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";
import ErrorIcon from "@mui/icons-material/Error";

import ListIcon from "@mui/icons-material/List";

const renderStatus = (record) => {
  const status = record.status;
  switch (status) {
    case "SUCCESSFUL":
      return (
        <Box
          sx={{ display: "flex", alignItems: "center", color: "success.main" }}
        >
          <CheckCircleIcon sx={{ mr: 1 }} />
          SUCCESSFUL
        </Box>
      );
    case "PENDING":
      return (
        <Box
          sx={{ display: "flex", alignItems: "center", color: "warning.main" }}
        >
          <HourglassEmptyIcon sx={{ mr: 1 }} />
          PENDING
        </Box>
      );
    case "FAILED":
    case "CANCELLED":
    default:
      return (
        <Box
          sx={{ display: "flex", alignItems: "center", color: "error.main" }}
        >
          <ErrorIcon sx={{ mr: 1 }} />
          {status.charAt(0).toUpperCase() + status.slice(1)}
        </Box>
      );
  }
};

export const SumupTransactionList = (props) => {
  const {
    sumupEnabled,
    currencyOptions: currency,
    Locale: locale,
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
      filters={<TimeRangeFilter />}
      empty={false}
      {...props}
    >
      <Datagrid rowClick="show" bulkActionButtons={false}>
        <TextField source="id" sortable={false} />
        <TextField source="transactionCode" sortable={false} />
        <NumberField
          source="amount"
          locales={locale}
          options={currency}
          sortable={false}
        />
        <TextField source="currency" sortable={false} />
        <DateField
          source="createdAt"
          showTime={true}
          locales={locale}
          options={{ weekday: "short", hour: "2-digit", minute: "2-digit" }}
          sortable={false}
        />
        <FunctionField label="Status" render={renderStatus} />
      </Datagrid>
    </List>
  );
};

const getOldestTimeOptions = () => {
  const now = Date.now();

  const toISO = (msAgo) => new Date(now - msAgo).toISOString();

  return [
    { id: toISO(10 * 60 * 1000), name: "Last 10 minutes" },
    { id: toISO(60 * 60 * 1000), name: "Last 1 hour" },
    { id: toISO(24 * 60 * 60 * 1000), name: "Last 24 hours" },
    { id: toISO(48 * 60 * 60 * 1000), name: "Last 48 hours" },
    { id: toISO(7 * 24 * 60 * 60 * 1000), name: "Last 1 week" },
    { id: toISO(30 * 24 * 60 * 60 * 1000), name: "Last 1 month" },
  ];
};

const TimeRangeFilter = (props) => (
  <Filter {...props}>
    <SelectInput
      source="oldest_time"
      label="Time Range"
      choices={getOldestTimeOptions()}
      alwaysOn
    />
  </Filter>
);

export const SumupTransactionShow = (props) => {
  const currency = useConfig().currencyOptions;
  const locale = useConfig().Locale;

  return (
    <Show {...props}>
      <SimpleShowLayout>
        <TextField source="transactionId" />
        <TextField source="transactionCode" />
        <SimpleShowLayout direction="row">
          <NumberField source="amount" locales={locale} options={currency} />
          <TextField source="currency" />
        </SimpleShowLayout>
        <DateField
          source="createdAt"
          showDate={true}
          showTime={true}
          locales={locale}
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
        <FunctionField label="Status" render={renderStatus} />
        <ArrayField source="events">
          <Datagrid bulkActionButtons={false}>
            <TextField source="id" />
            <TextField source="type" />
            <TextField source="status" />
            <DateField
              source="timestamp"
              showDate={true}
              showTime={true}
              locales={locale}
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

export const SumupTransactionIcon = () => <ListIcon />;
