import React from "react";
import { TextInput, TextInputProps, Validator } from "react-admin";
import Decimal from "decimal.js";

const decimalValidator: Validator = (value: unknown) => {
  if (value === null || value === undefined || value === "") {
    return undefined;
  }

  try {
    const decimalValue = new Decimal(String(value));

    if (decimalValue.isNaN()) return "Invalid number";
    if (decimalValue.isNegative()) return "Negative number";

    return undefined;
  } catch {
    return "Invalid number";
  }
};

const DecimalInput: React.FC<TextInputProps> = ({ validate, ...props }) => {
  const compositeValidate = Array.isArray(validate)
    ? [decimalValidator, ...validate]
    : validate
      ? [decimalValidator, validate]
      : decimalValidator;

  const parse = (value: string | null): string | null => {
    if (!value) return null;

    const cleaned = value
      .trim()
      .replace(",", ".")
      .replaceAll(/[^\d.]/g, "");

    const parts = cleaned.split(".");
    return parts.length > 2
      ? `${parts[0]}.${parts.slice(1).join("")}`
      : cleaned;
  };

  const format = (value: unknown): string => {
    if (value === null || value === undefined) return "";
    return String(value).replace(".", ",");
  };

  return (
    <TextInput
      {...props}
      parse={parse}
      format={format}
      validate={compositeValidate}
    />
  );
};

export default DecimalInput;
