import React from "react";
import { Chip, ChipProps } from "@mui/material";

// We extend ChipProps to stay compatible with MUI
interface QuickFilterProps extends ChipProps {
  label: string;
  source?: string; // React-admin uses these props
  defaultValue?: any; // to handle the filter logic
}

export const QuickFilter: React.FC<QuickFilterProps> = ({
  label,
  ...props
}) => {
  return <Chip sx={{ marginBottom: 1 }} label={label} {...props} />;
};
