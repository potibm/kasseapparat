import React from "react";
import {
  usePermissions,
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  NumberInput,
  Edit,
  SimpleForm,
  TextInput,
  Create,
  BooleanInput,
  SaveButton,
  SearchInput,
  Toolbar,
  ReferenceInput,
  SelectInput,
} from "react-admin";
import { Chip } from "@mui/material";
import GroupsIcon from "@mui/icons-material/Groups";
import PropTypes from "prop-types";

const QuickFilter = ({ label }) => {
  return <Chip sx={{ marginBottom: 1 }} label={label} />;
};

QuickFilter.propTypes = {
  label: PropTypes.string,
};

const ListEntryFilters = [
  <SearchInput source="q" alwaysOn key="ID" />,
  <ReferenceInput source="list" reference="lists" key="ID">
    <SelectInput label="List" source="list" optionText="name" />
  </ReferenceInput>,
  <ReferenceInput source="listGroup" reference="listGroups" key="ID">
    <SelectInput label="List Group" source="listGroup" optionText="name" />
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

export const ListEntryList = () => {
  const { permissions } = usePermissions();

  return (
    <List sort={{ field: "id", order: "ASC" }} filters={ListEntryFilters}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <TextField source="list.name" />
        <TextField source="listGroup.name" />
        <NumberField source="additionalGuests" />
        {permissions === "admin" && <DeleteButton mutationMode="pessimistic" />}
      </Datagrid>
    </List>
  );
};

export const ListEntryEdit = () => {
  return (
    <Edit>
      <SimpleForm
        toolbar={
          <Toolbar>
            <SaveButton />
          </Toolbar>
        }
      >
        <NumberInput disabled source="id" />
        <TextInput source="name" />
        <BooleanInput source="typeCode" />
      </SimpleForm>
    </Edit>
  );
};

export const ListEntryCreate = () => {
  return (
    <Create title="Create new List Entry">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="name" />
      </SimpleForm>
    </Create>
  );
};

export const ListEntryIcon = () => <GroupsIcon />;
