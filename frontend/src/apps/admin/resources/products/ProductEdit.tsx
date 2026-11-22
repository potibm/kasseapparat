import React from "react";
import { useWatch } from "react-hook-form";
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
import { useConfig } from "@core/config/hooks/useConfig";
import { OutOfStockBehavior } from "@core/config/types/config.types";
import DecimalInput from "../../ui/inputs/DecimalInput";
import GrossPriceInput from "../../ui/inputs/GrossPriceInput";
import VatRateInput from "../../ui/inputs/VatRateInput";
import ProductLayoutInputs from "../productGroups/ProductLayoutInputs";
import { Product } from "./types";

const HiddenInput: React.FC<{ outOfStockBehavior: OutOfStockBehavior }> = ({
  outOfStockBehavior,
}) => {
  const totalStock = useWatch({ name: "totalStock" }) ?? 0;
  const limited = Number(totalStock) > 0;
  const managed = limited && outOfStockBehavior === "auto_hide";
  return (
    <BooleanInput
      source="hidden"
      disabled={managed}
      helperText={
        managed
          ? "Managed automatically while stock is limited (auto_hide)."
          : "Hide this product from the POS"
      }
    />
  );
};

const SoldOutInput: React.FC<{ outOfStockBehavior: OutOfStockBehavior }> = ({
  outOfStockBehavior,
}) => {
  const totalStock = useWatch({ name: "totalStock" }) ?? 0;
  const limited = Number(totalStock) > 0;
  const managed = limited && outOfStockBehavior === "auto_sold_out";
  return (
    <BooleanInput
      source="soldOut"
      disabled={managed}
      helperText={
        managed
          ? "Managed automatically while stock is limited (auto_sold_out)."
          : "Still show the product to collect information how big the interest is"
      }
    />
  );
};

export const ProductEdit: React.FC = () => {
  const { currencyOptions, outOfStockBehavior } = useConfig();

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
          <ProductLayoutInputs />
          <NumberInput
            source="pos"
            helperText="The products will be shown in this order"
          />
          <HiddenInput outOfStockBehavior={outOfStockBehavior} />
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
          <SoldOutInput outOfStockBehavior={outOfStockBehavior} />
          <NumberInput source="soldOutRequestCount" disabled={true} />
        </FormTab>
        <FormTab label="API">
          <BooleanInput source="apiExport" />
        </FormTab>
      </TabbedForm>
    </Edit>
  );
};
