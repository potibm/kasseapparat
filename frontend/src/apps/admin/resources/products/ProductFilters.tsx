import { ReferenceInput, SelectInput } from "react-admin";
import { sectionLabel } from "../sections/section.utils";
import { Section } from "../sections/types";

export const ProductFilters = [
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
