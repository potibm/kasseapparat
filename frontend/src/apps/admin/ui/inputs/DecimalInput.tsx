import React from "react";
import { TextInput, TextInputProps, Validator } from "react-admin";
import Decimal from "decimal.js";

const DecimalInput: React.FC<TextInputProps> = (props) => {
  const validate: Validator = (value: any) => {
    if (!value && props.validate) {
      return undefined;
    }
    if (!value) return "Required";

    try {
      const decimalValue = new Decimal(value);

      if (decimalValue.isNaN()) {
        return "Invalid number";
      }
      if (decimalValue.isNegative()) {
        return "Negative number";
      }

      return undefined;
    } catch {
      return "Invalid number";
    }
  };

  const parse = (value: string | null): string | null => {
    if (!value) return null;

    return value
      .trim()
      .replaceAll(",", ".") // Replace all commas
      .replaceAll(/\.(?=.*\.)/g, ""); // Keep only last decimal point
  };

  const format = (value: any) => {
    if (value === null || value === undefined) return "";
    return value.toString().replace(",", ".");
  };

  return (
    <TextInput {...props} parse={parse} format={format} validate={validate} />
  );
};

export default DecimalInput;
