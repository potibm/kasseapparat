import React from "react";
import {
  RadioButtonGroupInput,
  RadioButtonGroupInputProps,
  required,
} from "react-admin";
import { useConfig } from "@core/config/hooks/useConfig";

interface VatRateType {
  name: string;
  rate: number;
}

export const VatRateInput: React.FC<RadioButtonGroupInputProps> = (props) => {
  const { vatRates } = useConfig();

  const vatChoices = vatRates || [];

  const vatOptionRenderer = (choice: VatRateType) => {
    return `${choice.name} (${choice.rate}%)`;
  };

  return (
    <RadioButtonGroupInput
      source="vatRate"
      choices={vatChoices}
      validate={[required()]}
      optionValue="rate"
      optionText={vatOptionRenderer}
      {...props}
    />
  );
};

export default VatRateInput;
