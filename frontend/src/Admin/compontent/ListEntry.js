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
  SearchInput,
  Toolbar,
  ReferenceInput,
  SelectInput,
  required,
  BooleanField,
  useRecordContext,
} from "react-admin";
import { Chip } from "@mui/material";
import PersonIcon from "@mui/icons-material/Person";
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

const AttendedGuestsBooleanField = (props) => {
  const record = useRecordContext();
  if (!record) return null;
  const hasAttendedGuests = record.attendedGuests > 0;

  // Erstellen eines temporären Records, der dem BooleanField übergeben wird
  const tempRecord = { ...record, hasAttendedGuests };

  return (
    <BooleanField {...props} source="hasAttendedGuests" record={tempRecord} />
  );
};

export const ListEntryList = () => {
  const { permissions } = usePermissions();

  return (
    <List sort={{ field: "id", order: "ASC" }} filters={ListEntryFilters}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <TextField source="list.name" />
        <TextField source="listGroup.name" />
        <NumberField source="additionalGuests" sortable={false} />
        <AttendedGuestsBooleanField label="present" sortable={false} />
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
        <ReferenceInput source="listId" reference="lists">
          <SelectInput optionText="name" validate={required()} disabled />
        </ReferenceInput>
        <ReferenceInput source="listGroupId" reference="listGroups">
          <SelectInput optionText="name" defaultValue={null} disabled />
        </ReferenceInput>
        <TextInput source="name" validate={required()} />
        <TextInput source="code" helperText="The entrance code on the ticket" />
        <NumberInput
          source="additionalGuests"
          min={0}
          helperText="Number of additional guests (read as +1)"
        />
        <NumberInput
          source="attendedGuests"
          min={0}
          helperText="Number of visitors that are present"
        />
      </SimpleForm>
    </Edit>
  );
};

export const ListEntryCreate = () => {
  return (
    <Create title="Create new List Entry">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <ReferenceInput source="listId" reference="lists">
          <SelectInput optionText="name" validate={required()} />
        </ReferenceInput>
        <ReferenceInput source="listGroupId" reference="listGroups">
          <SelectInput optionText="name" defaultValue={null} />
        </ReferenceInput>
        <TextInput source="name" validate={required()} />
        <TextInput source="code" helperText="The entrance code on the ticket" />
        <NumberInput
          source="additionalGuests"
          min={0}
          defaultValue={0}
          helperText="Number of additional guests (read as +1)"
        />
      </SimpleForm>
    </Create>
  );
};

export const ListEntryIcon = () => <PersonIcon />;
