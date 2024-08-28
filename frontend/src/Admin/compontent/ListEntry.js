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
  BooleanField,
  useRecordContext,
  useGetIdentity,
  AutocompleteInput,
  useRedirect,
  useNotify,
} from "react-admin";
import PersonIcon from "@mui/icons-material/Person";
import ListEntryActions from "./ListEntryAction";
import { ListEntryFilters } from "./ListEntryFilters";
import { useLocation } from "react-router-dom";
import { useFormContext } from "react-hook-form";
import PropTypes from "prop-types";

const ConditionalDeleteButton = (props) => {
  const record = useRecordContext(props);
  const { permissions, isLoading: permissionsLoading } = usePermissions();
  const { data: identity, isLoading: identityLoading } = useGetIdentity();
  if (permissionsLoading || identityLoading) return <>Loading...</>;

  const currentUserId = identity.id;
  const createdByCurrentUser = record && record.createdById === currentUserId;

  if (permissions === "admin" || createdByCurrentUser) {
    return <DeleteButton {...props} />;
  }
  return null;
};

const AttendedGuestsBooleanField = (props) => {
  const record = useRecordContext();
  if (!record) return null;
  const hasAttendedGuests = record.attendedGuests > 0;

  const tempRecord = { ...record, hasAttendedGuests };

  return (
    <BooleanField {...props} source="hasAttendedGuests" record={tempRecord} />
  );
};

export const ListEntryList = (props) => {
  return (
    <List
      sort={{ field: "id", order: "ASC" }}
      filters={ListEntryFilters}
      actions={<ListEntryActions />}
    >
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <TextField source="list.name" />
        <NumberField source="additionalGuests" sortable={false} />
        <AttendedGuestsBooleanField label="present" sortable={false} />
        <ConditionalDeleteButton mutationMode="pessimistic" />
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

const ListEntryCreateToolbar = ({ guestlistId, ...props }) => {
  const redirect = useRedirect();
  const { reset } = useFormContext();
  const notify = useNotify();

  const onSuccess = (data) => {
    notify(`Entry saved!`);
    reset();
    redirect(`/admin/listEntries/create?list_id=${data.listId}`);
  };

  return (
    <Toolbar {...props}>
      <SaveButton />
      <SaveButton
        type="button"
        label="Save and Add another"
        mutationOptions={{ onSuccess }}
        style={{ marginLeft: "10px" }}
      />
    </Toolbar>
  );
};
ListEntryCreateToolbar.propTypes = {
  guestlistId: PropTypes.number.isRequired,
};

export const ListEntryCreate = (props) => {
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const guestlistId = parseInt(params.get("list_id"), 10);

  return (
    <Create {...props} title="Create new List Entry">
      <SimpleForm
        defaultValues={{ listId: guestlistId }}
        toolbar={<ListEntryCreateToolbar guestlistId={guestlistId} />}
      >
        <NumberInput disabled source="id" />
        <ReferenceInput source="listId" reference="lists">
          <AutocompleteInput optionText="name" validate={required()} />
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
