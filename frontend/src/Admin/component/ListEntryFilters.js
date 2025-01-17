import React from "react";
import { SearchInput, ReferenceInput, AutocompleteInput } from "react-admin";
import { Chip } from "@mui/material";
import PropTypes from "prop-types";

const QuickFilter = ({ label }) => {
  return <Chip sx={{ marginBottom: 1 }} label={label} />;
};

QuickFilter.propTypes = {
  label: PropTypes.string,
};

export const ListEntryFilters = [
  <SearchInput source="q" alwaysOn key="ID" />,
  <ReferenceInput
    source="list"
    reference="lists"
    key="id"
    sort={{ field: "name", order: "ASC" }}
  >
    <AutocompleteInput optionText="name" />
  </ReferenceInput>,
  <QuickFilter
    source="isPresent"
    label="Present"
    defaultValue={true}
    key="ID"
  />,
  <QuickFilter
    source="isNotPresent"
    label="Not Present"
    defaultValue={true}
    key="ID"
  />,
];
