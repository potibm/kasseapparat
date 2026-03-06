import React from "react";
import { SearchInput } from "react-admin";
import { Chip } from "@mui/material";
import PropTypes from "prop-types";

const QuickFilter = ({ label }) => {
  return <Chip sx={{ marginBottom: 1 }} label={label} />;
};

QuickFilter.propTypes = {
  label: PropTypes.string,
};

export const UserFilters = [
  <SearchInput source="q" alwaysOn key="ID" />,
  <QuickFilter source="isAdmin" label="Admin" defaultValue={true} key="ID" />,
];
