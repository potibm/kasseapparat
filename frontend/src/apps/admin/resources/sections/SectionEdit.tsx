import React from "react";
import { Edit, SaveButton, SimpleForm, Toolbar } from "react-admin";
import SectionDeleteButton from "./SectionDeleteButton";
import SectionFields from "./SectionForm";
import { Section } from "./types";

export const SectionEdit: React.FC = () => (
  <Edit<Section> mutationMode="pessimistic" title="Edit section">
    <SimpleForm
      toolbar={
        <Toolbar sx={{ justifyContent: "space-between" }}>
          <SaveButton />
          <SectionDeleteButton />
        </Toolbar>
      }
    >
      <SectionFields />
    </SimpleForm>
  </Edit>
);

export default SectionEdit;
