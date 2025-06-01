import React from "react";
import {
  ReferenceInput,
  AutocompleteInput,
  SelectInput,
  NumberInput,
  FilterForm,
  FilterButton,
  ExportButton,
  Toolbar,
} from "react-admin";
import { useConfig } from "../../provider/ConfigProvider";
import { PurchaseExportButton } from "./PurchaseExportButton";

export const PurchaseListToolbar = (props) => {
  const { paymentMethods } = useConfig();

  const statusChoices = [
    { id: "pending", name: "Pending" },
    { id: "confirmed", name: "Confirmed" },
    { id: "cancelled", name: "Cancelled" },
    { id: "failed", name: "Failed" },
  ];

  const paymentMethodChoices =
    paymentMethods?.map((pm) => ({
      id: pm.code,
      name: pm.name,
    })) ?? [];

  const purchaseFilters = [
    <ReferenceInput key="created-by-id" source="createdById" reference="users">
      <AutocompleteInput optionText="username" />
    </ReferenceInput>,
    <SelectInput
      key="payment-method"
      source="paymentMethod"
      choices={paymentMethodChoices}
    />,
    <SelectInput key="status" source="status" choices={statusChoices} />,
    <NumberInput
      key="total-gross-price-gte"
      source="totalGrossPrice_gte"
      label="Min. total gross price"
      step={0.01}
      resettable={"true"}
    />,
    <NumberInput
      key="total-gross-price-lte"
      source="totalGrossPrice_lte"
      label="Max. total gross price"
      step={0.01}
      resettable={"true"}
    />,
  ];

  return (
    <Toolbar {...props}>
      <FilterForm filters={purchaseFilters} />
      <div className="MuiToolbar-root MuiToolbar-dense css-oynclu-MuiToolbar-root-RaTopToolbar-root">
        <FilterButton filters={purchaseFilters} />
        <PurchaseExportButton paymentMethods={paymentMethods} />
        <ExportButton {...props} />
      </div>
    </Toolbar>
  );
};
