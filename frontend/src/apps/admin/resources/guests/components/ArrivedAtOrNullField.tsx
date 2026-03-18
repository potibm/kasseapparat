import React from "react";
import {
  useRecordContext,
  DateField,
  DateFieldProps,
  RaRecord,
} from "react-admin";
import { Typography } from "@mui/material";

interface GuestRecord extends RaRecord {
  arrivedAt?: string | Date | null;
}

interface ArrivedAtOrNullFieldProps extends Omit<DateFieldProps, "source"> {
  source?: string;
}

export const ArrivedAtOrNullField: React.FC<ArrivedAtOrNullFieldProps> = (
  props,
) => {
  const record = useRecordContext<GuestRecord>();
  if (!record) return null;
  const hasArrived = record.arrivedAt != null;

  if (!hasArrived) {
    return (
      <Typography
        component="span"
        variant="body2"
        sx={{ color: "text.secondary", fontStyle: "italic" }}
      >
        The person has not arrived, yet.
      </Typography>
    );
  }

  return <DateField {...props} source="arrivedAt" showTime={true} />;
};

export default ArrivedAtOrNullField;
