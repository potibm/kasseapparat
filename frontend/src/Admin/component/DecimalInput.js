import React from "react";
import { TextInput } from "react-admin";
import Decimal from "decimal.js";

const DecimalInput = (props) => {
  const validate = (value) => {
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

  const parse = (value) => {
    if (!value) return null;

    const normalizedValue = value.replace(",", ".");
    return normalizedValue;
  };

  const format = (value) => {
    if (!value) return "";
    return value.toString().replace(",", ".");
  };

  return (
    <TextInput {...props} parse={parse} format={format} validate={validate} />
  );
};

export default DecimalInput;
