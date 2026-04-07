import React from "react";
import {
  NumberInput,
  SimpleForm,
  TextInput,
  Create,
  BooleanInput,
  required,
} from "react-admin";
import { useConfig } from "@core/config/hooks/useConfig";
import DecimalInput from "../../ui/inputs/DecimalInput";
import GrossPriceInput from "../../ui/inputs/GrossPriceInput";
import VatRateInput from "../../ui/inputs/VatRateInput";
import { Product } from "./types";

export const ProductCreate: React.FC = () => {
  const { currencyOptions } = useConfig();

  return (
    <Create<Product> title="Create new product">
      <SimpleForm>
        <TextInput source="name" validate={required()} />
        <h6>Pricing</h6>
        <DecimalInput source="netPrice" />
        <VatRateInput row={false} />
        <GrossPriceInput
          netSource="netPrice"
          vatSource="vatRate"
          source="grossPrice"
          options={currencyOptions}
        />
        <h6>Layout</h6>
        <NumberInput
          source="pos"
          helperText="The products will be shown in this order"
          validate={[required()]}
        />
        <BooleanInput
          source="wrapAfter"
          helperText="Create a line break after this product"
        />
      </SimpleForm>
    </Create>
  );
};
