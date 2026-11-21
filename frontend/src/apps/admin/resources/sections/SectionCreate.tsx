import React from "react";
import { Create, SimpleForm } from "react-admin";
import SectionFields from "./SectionForm";
import { Section } from "./types";

export const SectionCreate: React.FC = () => (
  <Create<Section> title="Create new section" redirect="list">
    <SimpleForm defaultValues={{ name: "", pos: 0 }}>
      <SectionFields />
    </SimpleForm>
  </Create>
);

export default SectionCreate;
