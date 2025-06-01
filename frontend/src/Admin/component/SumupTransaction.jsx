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
  ShowActions,
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
        <h2>SumUp Readers</h2>
        <p>
          SumUp integration is not enabled. Please enable it in the
          configuration.
        </p>
      </div>
    );
  }

  return (
    <List title="SumUp Readers" {...props}>
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

export const SumupTransactionShow = (props) => {
  const currency = useConfig().currencyOptions;
  const locale = useConfig().Locale;

  return (
    <Show {...props}>
      <SimpleShowLayout>
        <TextField source="id" sortable={false} />
        <TextField source="transactionCode" sortable={false} />
        <SimpleShowLayout direction="row">
          <NumberField
            source="amount"
            locales={locale}
            options={currency}
            sortable={false}
          />
          <TextField source="currency" sortable={false} />
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
          sortable={false}
        />
        <FunctionField label="Status" render={renderStatus} />
      </SimpleShowLayout>
    </Show>
  );
};

export const SumupTransactionIcon = () => <ListIcon />;
