import React from "react";
import {
  TopToolbar,
  ExportButton,
  FilterButton,
  CreateButton,
} from "react-admin";
import { ImportDeineTicketsButton } from "./ImportDeineTicketsButton";

const GuestActions: React.FC = () => (
  <TopToolbar>
    <FilterButton />
    <CreateButton />
    <ExportButton />
    <ImportDeineTicketsButton />
  </TopToolbar>
);

export default GuestActions;
