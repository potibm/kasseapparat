import React, { useState, useEffect } from "react";
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
  Button,
  useListContext,
} from "react-admin";
import { useConfig } from "../../provider/ConfigProvider";
import { Box, Tooltip } from "@mui/material";
import CreditCardIcon from "@mui/icons-material/CreditCard";
import LinkOffIcon from "@mui/icons-material/LinkOff";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";
import TouchAppIcon from "@mui/icons-material/TouchApp";
import ErrorIcon from "@mui/icons-material/Error";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import {
  getCurrentReaderId,
  setCurrentReaderId,
  clearCurrentReaderId,
} from "../../helper/ReaderCookie";
import PropTypes from "prop-types";

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

const renderReaderSelection = (record, selectedReaderId, onSelect) => {
  if (record.status !== "paired") return null;

  const isCurrent = selectedReaderId === record.id;

  return isCurrent ? (
    <Box sx={{ display: "flex", alignItems: "center", color: "success.main" }}>
      <CheckCircleOutlineIcon sx={{ mr: 1 }} />
      Selected
    </Box>
  ) : (
    <Tooltip title="Assign this reader to this device">
      <span>
        <Button
          onClick={() => onSelect(record.id)}
          size="small"
          startIcon={<TouchAppIcon />}
          variant="text"
          color="primary"
          label="Use this reader"
        />
      </span>
    </Tooltip>
  );
};

export const SumupReaderList = (props) => {
  const sumupEnabled = useConfig().sumupEnabled;
  const [selectedReaderId, setSelectedReaderId] = useState(undefined);

  useEffect(() => {
    setSelectedReaderId(getCurrentReaderId());
  }, []);

  const handleReaderSelect = (id) => {
    setCurrentReaderId(id);
    setSelectedReaderId(id);
  };

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
      {/* Innerhalb von List: Jetzt ist ListContext verfügbar */}
      <ReaderListContent
        selectedReaderId={selectedReaderId}
        onClear={() => {
          clearCurrentReaderId();
          setSelectedReaderId(undefined);
        }}
        onSelect={handleReaderSelect}
      />
    </List>
  );
};

const ReaderListContent = ({ selectedReaderId, onClear, onSelect }) => {
  const { data: readers = [], isLoading } = useListContext();

  const isSelectedReaderMissing =
    !isLoading &&
    selectedReaderId &&
    !readers.some((r) => r.id === selectedReaderId);

  return (
    <>
      {isSelectedReaderMissing && (
        <Box
          sx={{
            p: 2,
            mb: 2,
            backgroundColor: "#fff3cd",
            color: "#856404",
            borderRadius: 1,
            border: "1px solid #ffeeba",
          }}
        >
          <Box sx={{ mb: 1 }}>
            The previously selected reader ({selectedReaderId}) is no longer
            available or paired.
          </Box>
          <Button variant="outlined" size="small" onClick={onClear}>
            Clear selection
          </Button>
        </Box>
      )}

      <Datagrid
        rowClick=""
        bulkActionButtons={false}
        rowSx={(record) =>
          record.id === selectedReaderId ? { backgroundColor: "#e0f7fa" } : {}
        }
      >
        <TextField source="id" />
        <TextField source="name" />
        <FunctionField label="Status" render={renderStatus} />
        <TextField source="deviceIdentifier" />
        <TextField source="deviceModel" />
        <FunctionField
          label="Action"
          render={(record) =>
            renderReaderSelection(record, selectedReaderId, onSelect)
          }
        />
        <DeleteWithConfirmButton
          label="Unpair"
          confirmTitle="Unpair device"
          confirmContent="Are you sure you want to unpair this SumUp reader?"
          mutationMode="pessimistic"
          icon={<LinkOffIcon />}
        />
      </Datagrid>
    </>
  );
};

ReaderListContent.propTypes = {
  selectedReaderId: PropTypes.string,
  onClear: PropTypes.func.isRequired,
  onSelect: PropTypes.func.isRequired,
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
