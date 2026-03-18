import { RaRecord, useRecordContext, FieldProps } from "react-admin";
import { Box } from "@mui/material";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";
import ErrorIcon from "@mui/icons-material/Error";

export const PurchaseStatusField = <T extends RaRecord>(
  props: FieldProps<T>,
) => {
  const { source } = props;
  const record = useRecordContext<T>(props);

  if (!record || !source || !record[source]) return null;

  const status = record[source];

  // We define the UI mapping here
  const statusConfig = {
    confirmed: {
      color: "success.main",
      icon: <CheckCircleIcon fontSize="small" sx={{ mr: 1 }} />,
      text: "CONFIRMED",
    },
    pending: {
      color: "warning.main",
      icon: <HourglassEmptyIcon fontSize="small" sx={{ mr: 1 }} />,
      text: "PENDING",
    },
    // Default / Error fallback
    fallback: {
      color: "error.main",
      icon: <ErrorIcon fontSize="small" sx={{ mr: 1 }} />,
      text: status.toUpperCase(),
    },
  };

  const config =
    statusConfig[status as keyof typeof statusConfig] || statusConfig.fallback;

  return (
    <Box
      component="span" // span is better for table cells
      sx={{ display: "flex", alignItems: "center", color: config.color }}
    >
      {config.icon}
      {config.text}
    </Box>
  );
};
