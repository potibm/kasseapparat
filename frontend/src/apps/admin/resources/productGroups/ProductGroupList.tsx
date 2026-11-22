import React from "react";
import {
  Datagrid,
  List,
  NumberField,
  ReferenceField,
  ReferenceInput,
  SelectInput,
  TextField,
} from "react-admin";
import { sectionLabel } from "../sections/section.utils";
import { Section } from "../sections/types";
import ProductGroupDeleteButton from "./ProductGroupDeleteButton";
import { ProductGroup } from "./types";

const productGroupFilters = [
  <ReferenceInput
    source="sectionId"
    reference="sections"
    sort={{ field: "pos", order: "ASC" }}
    perPage={100}
    alwaysOn
    key="sectionId"
  >
    <SelectInput
      label="Section"
      optionText={(record: Section) => sectionLabel(record)}
    />
  </ReferenceInput>,
];

export const ProductGroupList: React.FC = () => (
  <List<ProductGroup>
    sort={{ field: "pos", order: "ASC" }}
    filters={productGroupFilters}
  >
    <Datagrid rowClick="edit" bulkActionButtons={false}>
      <NumberField source="id" />
      <TextField source="name" />
      <ReferenceField
        source="sectionId"
        reference="sections"
        label="Section"
        link={false}
      />
      <NumberField source="pos" />
      <ProductGroupDeleteButton />
    </Datagrid>
  </List>
);

export default ProductGroupList;
