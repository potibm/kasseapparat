import React from "react";
import {
  BooleanField,
  useRecordContext,
  BooleanFieldProps,
  RaRecord,
} from "react-admin";

interface GuestlistRecord extends RaRecord {
  attendedGuests?: number;
}

interface AttendedGuestsBooleanFieldProps extends Omit<
  BooleanFieldProps,
  "source"
> {
  source?: string;
}

export const AttendedGuestsBooleanField: React.FC<
  AttendedGuestsBooleanFieldProps
> = (props) => {
  const record = useRecordContext<GuestlistRecord>();
  if (!record) return null;
  const hasAttendedGuests = (record.attendedGuests ?? 0) > 0;
  const tempRecord = { ...record, hasAttendedGuests };

  return (
    <BooleanField {...props} source="hasAttendedGuests" record={tempRecord} />
  );
};

export default AttendedGuestsBooleanField;
