import React from "react";
import {
  NumberInput,
  Edit,
  TextInput,
  BooleanInput,
  SaveButton,
  Toolbar,
  required,
  TabbedForm,
  FormTab,
} from "react-admin";
import { useConfig } from "../../../../core/config/providers/ConfigProvider";
import DecimalInput from "../../ui/inputs/DecimalInput";
import GrossPriceInput from "../../ui/inputs/GrossPriceInput";
import VatRateInput from "../../ui/inputs/VatRateInput";
import { Product } from "./types";

export const ProductEdit: React.FC = () => {
  const { currencyOptions } = useConfig();

  return (
    <Edit<Product> mutationMode="pessimistic" title="Edit product">
      <TabbedForm
        toolbar={
          <Toolbar>
            <SaveButton />
          </Toolbar>
        }
      >
        <FormTab label="General">
          <NumberInput disabled source="id" />
          <TextInput source="name" validate={required()} />
        </FormTab>
        <FormTab label="Pricing">
          <DecimalInput source="netPrice" />
          <VatRateInput row={false} />
          <GrossPriceInput
            netSource="netPrice"
            vatSource="vatRate"
            source="grossPrice"
            options={currencyOptions}
          />
        </FormTab>
        <FormTab label="Layout">
          <NumberInput
            source="pos"
            helperText="The products will be shown in this order"
          />
          <BooleanInput
            source="wrapAfter"
            helperText="Create a line break after this product"
          />
          <BooleanInput
            source="hidden"
            helperText="Hide this product from the POS"
          />
        </FormTab>
        <FormTab label="Stock">
          <h3>Stock</h3>
          <NumberInput
            source="totalStock"
            min={0}
            step={1}
            helperText="Number of available items. Shown for informational purposes, only."
          />
          <NumberInput
            source="unitsSold"
            min={0}
            step={1}
            disabled={true}
            helperText="Number of sold items. Shown for informational purposes, only."
          />
        </FormTab>
        <FormTab label="Sold Out">
          <BooleanInput
            source="soldOut"
            helperText="Still show the product to collect information how big the interest is"
          />
          <NumberInput source="soldOutRequestCount" disabled={true} />
        </FormTab>
        <FormTab label="API">
          <BooleanInput source="apiExport" />
        </FormTab>
      </TabbedForm>
    </Edit>
  );
};
