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
  SaveButton,
  Toolbar,
  ReferenceInput,
  SelectInput,
  required,
} from "react-admin";
import { Chip } from "@mui/material";
import GroupIcon from "@mui/icons-material/Group";
import PropTypes from "prop-types";

const QuickFilter = ({ label }) => {
  return <Chip sx={{ marginBottom: 1 }} label={label} />;
};

QuickFilter.propTypes = {
  label: PropTypes.string,
};

const ListGroupFilters = [
  <ReferenceInput source="list" reference="lists" key="ID">
    <SelectInput label="List" source="list" optionText="name" />
  </ReferenceInput>,
];

export const ListGroupList = () => {
  const { permissions } = usePermissions();

  return (
    <List sort={{ field: "id", order: "ASC" }} filters={ListGroupFilters}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <TextField source="list.name" />
        {permissions === "admin" && <DeleteButton mutationMode="pessimistic" />}
      </Datagrid>
    </List>
  );
};

export const ListGroupEdit = () => {
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
        <TextInput source="name" validate={required()} />
        <ReferenceInput source="listId" reference="lists">
          <SelectInput optionText="name" disabled />
        </ReferenceInput>
      </SimpleForm>
    </Edit>
  );
};

export const ListGroupCreate = () => {
  return (
    <Create title="Create new List Group">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="name" validate={required()} />
        <ReferenceInput source="listId" reference="lists">
          <SelectInput optionText="name" validate={required()} />
        </ReferenceInput>
      </SimpleForm>
    </Create>
  );
};

export const ListGroupIcon = () => <GroupIcon />;
