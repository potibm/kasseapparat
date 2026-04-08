import React from "react";
import {
  ReferenceInput,
  AutocompleteInput,
  SelectInput,
  SelectArrayInput,
  NumberInput,
  FilterForm,
  FilterButton,
  ExportButton,
  TopToolbar,
  ToolbarProps,
} from "react-admin";
import { Box } from "@mui/material";
import { useConfig } from "@core/config/hooks/useConfig";
import { PurchaseExportButton } from "./PurchaseExportButton";

export const PurchaseListToolbar: React.FC<ToolbarProps> = () => {
  const { paymentMethods } = useConfig();

  const statusChoices = [
    { id: "pending", name: "Pending" },
    { id: "confirmed", name: "Confirmed" },
    { id: "cancelled", name: "Cancelled" },
    { id: "failed", name: "Failed" },
    { id: "refunded", name: "Refunded" },
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
    <SelectArrayInput key="status" source="status" choices={statusChoices} />,
    <NumberInput
      key="total-gross-price-gte"
      source="totalGrossPrice_gte"
      label="Min. total gross price"
      step={0.01}
    />,
    <NumberInput
      key="total-gross-price-lte"
      source="totalGrossPrice_lte"
      label="Max. total gross price"
      step={0.01}
    />,
  ];

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "space-between",
        alignItems: "flex-start",
        width: "100%",
      }}
    >
      {/* Left side: The active filters / search */}
      <FilterForm filters={purchaseFilters} />

      {/* Right side: Actions like Filter-Add, Export, etc. */}
      <TopToolbar sx={{ minHeight: "auto", p: 0 }}>
        <FilterButton filters={purchaseFilters} />
        {paymentMethods && (
          <PurchaseExportButton paymentMethods={paymentMethods} />
        )}
        <ExportButton />
      </TopToolbar>
    </Box>
  );
};
