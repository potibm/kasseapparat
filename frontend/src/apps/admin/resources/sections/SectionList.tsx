import React from "react";
import {
  List,
  Datagrid,
  NumberField,
  FunctionField,
  ListProps,
} from "react-admin";
import SectionDeleteButton from "./SectionDeleteButton";
import { sectionLabel } from "./section.utils";
import { Section } from "./types";

export const SectionList: React.FC<ListProps> = (props) => (
  <List<Section> {...props} sort={{ field: "pos", order: "ASC" }}>
    <Datagrid rowClick="edit" bulkActionButtons={false}>
      <NumberField source="id" />
      <FunctionField<Section>
        source="name"
        label="Name"
        render={(record) => sectionLabel(record)}
      />
      <NumberField source="pos" />
      <SectionDeleteButton />
    </Datagrid>
  </List>
);

export default SectionList;
