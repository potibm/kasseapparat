import React from "react";
import { NumberInput, TextInput, required } from "react-admin";
import SectionInput from "../sections/SectionInput";

/** The fields shared by the create and the edit form. */
export const ProductGroupFields: React.FC<{ preselectSection?: boolean }> = ({
  preselectSection = false,
}) => (
  <>
    <TextInput
      source="name"
      validate={required()}
      helperText="Shown at the start of the row in the POS, e.g. “Male”"
    />
    <SectionInput preselectFirst={preselectSection} />
    <NumberInput
      source="pos"
      helperText="Groups and single products of a section share this order"
      validate={[required()]}
    />
  </>
);

export default ProductGroupFields;
