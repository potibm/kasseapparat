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
  ReferenceField,
} from "react-admin";
import { useConfig } from "../../provider/ConfigProvider";
import { Box } from "@mui/material";
import { PurchaseListToolbar } from "./PurchaseListToolbar";
import ListIcon from "@mui/icons-material/List";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";
import ErrorIcon from "@mui/icons-material/Error";

const renderPaymentMethod = (record, paymentMethods) => {
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
};

const renderStatus = (record) => {
  const status = record.status;
  switch (status) {
    case "confirmed":
      return (
        <Box
          sx={{ display: "flex", alignItems: "center", color: "success.main" }}
        >
          <CheckCircleIcon sx={{ mr: 1 }} />
          CONFIRMED
        </Box>
      );
    case "pending":
      return (
        <Box
          sx={{ display: "flex", alignItems: "center", color: "warning.main" }}
        >
          <HourglassEmptyIcon sx={{ mr: 1 }} />
          PENDING
        </Box>
      );
    case "failed":
    case "cancelled":
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
          render={(record) => renderPaymentMethod(record, paymentMethods)}
        />
        <FunctionField label="Status" render={renderStatus} />
        <DeleteButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};

export const PurchaseShow = (props) => {
  const { currencyOptions: currency, locale, paymentMethods } = useConfig();

  return (
    <Show {...props}>
      <SimpleShowLayout>
        <TextField source="id" />
        <DateField source="createdAt" showTime={true} />
        <FunctionField label="Status" render={renderStatus} />
        <FunctionField
          source="paymentMethod"
          render={(record) => renderPaymentMethod(record, paymentMethods)}
        />
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

export const PurchaseIcon = () => <ListIcon />;
