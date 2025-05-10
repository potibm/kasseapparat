import React from "react";
import {
  ReferenceInput,
  AutocompleteInput,
  SelectInput,
  NumberInput,
  Filter,
} from "react-admin";
import { useConfig } from "../../provider/ConfigProvider";

export const PurchaseFilters = (props) => {
  const { paymentMethods } = useConfig();

  const paymentMethodChoices =
    paymentMethods?.map((pm) => ({
      id: pm.code,
      name: pm.name,
    })) ?? [];

  return (
    <Filter {...props}>
      <ReferenceInput source="createdById" reference="users">
        <AutocompleteInput optionText="username" />
      </ReferenceInput>
      <SelectInput source="paymentMethod" choices={paymentMethodChoices} />
      <NumberInput
        source="totalGrossPrice_gte"
        label="Min. total gross price"
        step={0.01}
        resettable
      />
      <NumberInput
        source="totalGrossPrice_lte"
        label="Max. total gross price"
        step={0.01}
        resettable
      />
    </Filter>
  );
};
