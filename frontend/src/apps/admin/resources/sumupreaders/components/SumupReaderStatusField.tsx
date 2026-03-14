import { useRecordContext, FieldProps, RaRecord } from "react-admin";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";
import ErrorIcon from "@mui/icons-material/Error";
import { Typography, Box } from "@mui/material";

export const SumupReaderStatusField = <T extends RaRecord>(
  props: FieldProps<T>,
) => {
  const { source } = props;
  const record = useRecordContext<T>(props);

  if (!record || !source || record[source] === undefined) return null;

  const status = record[source];

  const getStatusConfig = () => {
    switch (status) {
      case "paired":
        return {
          color: "success.main",
          icon: <CheckCircleIcon sx={{ mr: 1, fontSize: "1.2rem" }} />,
          label: "Paired",
        };
      case "processing":
        return {
          color: "warning.main",
          icon: <HourglassEmptyIcon sx={{ mr: 1, fontSize: "1.2rem" }} />,
          label: "Processing",
        };
      case "expired":
      case "unknown":
      default:
        return {
          color: "error.main",
          icon: <ErrorIcon sx={{ mr: 1, fontSize: "1.2rem" }} />,
          label:
            typeof status === "string"
              ? status.charAt(0).toUpperCase() + status.slice(1)
              : "Unknown",
        };
    }
  };

  const config = getStatusConfig();

  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
        color: config.color,
      }}
    >
      {config.icon}
      <Typography variant="body2" component="span" sx={{ fontWeight: 500 }}>
        {config.label}
      </Typography>
    </Box>
  );
};

export default SumupReaderStatusField;
