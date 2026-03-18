import React from "react";
import { TextInput, TextInputProps, required } from "react-admin";

const validatePairingCode = (value: string) => {
  if (!value) return undefined;
  if (!/^[A-Z0-9]{9}$/.test(value.toUpperCase())) {
    return "Must be 9 characters: only letters and digits";
  }
  return undefined;
};

export const SumupReaderPairingCodeInput: React.FC<Partial<TextInputProps>> = (
  props,
) => {
  return (
    <TextInput
      source="pairingCode"
      label="Pairing Code"
      validate={[required(), validatePairingCode]}
      helperText="Enter the 9-character code using only uppercase letters and digits."
      parse={(value: string) => value?.toUpperCase()}
      slotProps={{
        htmlInput: {
          maxLength: 9,
          style: { textTransform: "uppercase" },
        },
      }}
      {...props}
    />
  );
};
