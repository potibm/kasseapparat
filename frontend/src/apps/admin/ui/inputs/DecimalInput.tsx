import React from "react";
import { TextInput, TextInputProps, Validator } from "react-admin";
import {
  parseDecimal,
  formatDecimal,
  decimalValidator,
} from "../../utils/decimal-utils";

const DecimalInput: React.FC<TextInputProps> = ({ validate, ...props }) => {
  let additionalValidators: Validator[] = [];

  if (Array.isArray(validate)) {
    additionalValidators = validate;
  } else if (validate) {
    additionalValidators = [validate];
  }

  const compositeValidate = [decimalValidator, ...additionalValidators];

  return (
    <TextInput
      {...props}
      parse={parseDecimal}
      format={formatDecimal}
      validate={compositeValidate}
    />
  );
};

export default DecimalInput;
