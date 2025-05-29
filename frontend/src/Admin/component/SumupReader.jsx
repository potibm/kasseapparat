import React from "react";
import {
  List,
  Datagrid,
  TextField,
  SimpleForm,
  TextInput,
  Create,
  SaveButton,
  Toolbar,
  DeleteWithConfirmButton,
  TopToolbar,
  CreateButton,
  FunctionField,
} from "react-admin";
import { useConfig } from "../../provider/ConfigProvider";
import { Box } from "@mui/material";
import CreditCardIcon from "@mui/icons-material/CreditCard";
import LinkOffIcon from "@mui/icons-material/LinkOff";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";
import ErrorIcon from "@mui/icons-material/Error";

const SumupListActions = () => (
  <TopToolbar>
    <CreateButton label="Pair" />
  </TopToolbar>
);

const renderStatus = (record) => {
  const status = record.status;
  switch (status) {
    case "paired":
      return (
        <Box
          sx={{ display: "flex", alignItems: "center", color: "success.main" }}
        >
          <CheckCircleIcon sx={{ mr: 1 }} />
          Paired
        </Box>
      );
    case "processing":
      return (
        <Box
          sx={{ display: "flex", alignItems: "center", color: "warning.main" }}
        >
          <HourglassEmptyIcon sx={{ mr: 1 }} />
          Processing
        </Box>
      );
    case "expired":
    case "unknown":
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

export const SumupReaderList = (props) => {
  const sumupEnabled = useConfig().sumupEnabled;
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
    <List title="SumUp Readers" {...props} actions={<SumupListActions />}>
      <Datagrid rowClick="" bulkActionButtons={false}>
        <TextField source="id" />
        <TextField source="name" />
        <FunctionField label="Status" render={renderStatus} />
        <TextField source="deviceIdentifier" />
        <TextField source="deviceModel" />
        <DeleteWithConfirmButton
          label="Unpair"
          confirmTitle="Unpair device"
          confirmContent="Are you sure you want to unpair this SumUp reader?"
          mutationMode="pessimistic"
          icon={<LinkOffIcon />}
        />
      </Datagrid>
    </List>
  );
};

const validatePairingCode = (value) => {
  if (!value) return "Pairing code is required";
  if (!/^[A-Za-z0-9]{9}$/.test(value)) {
    return "Must be 9 characters: only A–Z and 0–9";
  }
  return undefined;
};

export const SumupReaderCreate = () => {
  return (
    <Create title="Pair a SumUp Reader">
      <SimpleForm
        toolbar={
          <Toolbar>
            <SaveButton label="Pair" />
          </Toolbar>
        }
      >
        <TextInput
          source="pairingCode"
          label="Pairing Code"
          validate={validatePairingCode}
          helperText="Enter the 9-character code using only uppercase letters and digits."
          slotProps={{ htmlInput: { maxLength: 9 } }}
        />
        <TextInput source="name" />
      </SimpleForm>
    </Create>
  );
};

export const SumupReaderIcon = () => <CreditCardIcon />;
