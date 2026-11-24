import React from "react";
import { NumberInput, TextInput, required } from "react-admin";

/** The fields shared by the create and the edit form. */
export const SectionFields: React.FC = () => (
  <>
    <TextInput
      source="name"
      helperText="Shown as a heading in the POS. Leave empty for a plain row break."
    />
    <NumberInput
      source="pos"
      helperText="The sections will be shown in this order"
      validate={[required()]}
    />
  </>
);

export default SectionFields;
