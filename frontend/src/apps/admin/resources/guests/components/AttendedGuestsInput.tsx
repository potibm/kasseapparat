import React from "react";
import {
  NumberInput,
  NumberInputProps,
  minValue,
  maxValue,
  required,
  number,
} from "react-admin";
import { MAX_TOTAL_ATTENDEES } from "../constants";

const validateAttended = [
  required(),
  number(),
  minValue(0),
  maxValue(MAX_TOTAL_ATTENDEES),
];

export const AttendedGuestsInput: React.FC<Partial<NumberInputProps>> = (
  props,
) => (
  <NumberInput
    source="attendedGuests"
    label="Attended"
    defaultValue={0}
    min={0}
    max={MAX_TOTAL_ATTENDEES}
    validate={validateAttended}
    helperText={`Visitors present (max ${MAX_TOTAL_ATTENDEES})`}
    {...props}
  />
);
