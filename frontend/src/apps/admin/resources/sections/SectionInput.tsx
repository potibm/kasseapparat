import React, { useEffect } from "react";
import {
  ReferenceInput,
  SelectInput,
  required,
  useChoicesContext,
} from "react-admin";
import { useFormContext } from "react-hook-form";
import { sectionLabel } from "./section.utils";
import { Section } from "./types";

interface SectionSelectProps {
  preselectFirst?: boolean;
}

const SectionSelect: React.FC<SectionSelectProps> = ({ preselectFirst }) => {
  const { allChoices, source } = useChoicesContext<Section>();
  const { getValues, setValue } = useFormContext();

  useEffect(() => {
    if (!preselectFirst || !allChoices || allChoices.length === 0) return;
    if (getValues(source)) return;

    setValue(source, allChoices[0].id);
  }, [preselectFirst, allChoices, source, getValues, setValue]);

  return (
    <SelectInput
      label="Section"
      optionText={(record: Section) => sectionLabel(record)}
      validate={required()}
    />
  );
};

/**
 * Picks the section a product belongs to. On create the first section (lowest
 * position) is preselected, so a product never ends up without one.
 */
export const SectionInput: React.FC<SectionSelectProps> = ({
  preselectFirst = false,
}) => (
  <ReferenceInput
    source="sectionId"
    reference="sections"
    sort={{ field: "pos", order: "ASC" }}
    perPage={100}
  >
    <SectionSelect preselectFirst={preselectFirst} />
  </ReferenceInput>
);

export default SectionInput;
