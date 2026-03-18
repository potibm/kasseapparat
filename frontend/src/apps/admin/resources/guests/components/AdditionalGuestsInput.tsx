import React from "react";
import {
  NumberInput,
  NumberInputProps,
  required,
  number,
  minValue,
  maxValue,
} from "react-admin";
import { MAX_ADDITIONAL_GUESTS } from "../constants";

const validateAdditionalGuests = [
  required(),
  number(),
  minValue(0),
  maxValue(MAX_ADDITIONAL_GUESTS),
];

export const AdditionalGuestsInput: React.FC<Partial<NumberInputProps>> = (
  props,
) => {
  return (
    <NumberInput
      source="additionalGuests"
      label="Additional Guests"
      defaultValue={0}
      min={0}
      max={MAX_ADDITIONAL_GUESTS}
      validate={validateAdditionalGuests}
      helperText="Number of additional guests (read as +1)"
      {...props}
    />
  );
};

export default AdditionalGuestsInput;
