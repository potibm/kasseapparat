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
  TabbedForm,
  FormTab,
  DateField,
  DateTimeInput,
  minValue,
  maxValue,
  email,
  number,
} from "react-admin";
import PersonIcon from "@mui/icons-material/Person";
import GuestActions from "./GuestAction";
import { GuestFilters } from "./GuestFilters";
import { useLocation } from "react-router";
import { useFormContext } from "react-hook-form";
import PropTypes from "prop-types";
import { Box } from "@mui/material";
import { useConfig } from "../../provider/ConfigProvider";

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

const ValidateAdditionalGuests = [
  required(),
  number(),
  minValue(0),
  maxValue(5),
];
const ValidateAttendedGuests = [required(), number(), minValue(0), maxValue(6)];

const AttendedGuestsBooleanField = (props) => {
  const record = useRecordContext();
  if (!record) return null;
  const hasAttendedGuests = record.attendedGuests > 0;

  const tempRecord = { ...record, hasAttendedGuests };

  return (
    <BooleanField {...props} source="hasAttendedGuests" record={tempRecord} />
  );
};

const ArrivedAtOrNullField = (props) => {
  const record = useRecordContext();
  if (!record) return null;
  const hasArrived = record.arrivedAt != null;

  if (!hasArrived) {
    return <Box>The person has not arrived, yet</Box>;
  }

  return <DateTimeInput source="arrivedAt" disabled={true} />;
};

export const GuestList = (props) => {
  const locale = useConfig().Locale;

  return (
    <List
      sort={{ field: "id", order: "ASC" }}
      filters={GuestFilters}
      actions={<GuestActions />}
    >
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <TextField source="guestlist.name" label="Guestlist" />
        <NumberField
          source="additionalGuests"
          min={0}
          max={5}
          sortable={false}
        />
        <AttendedGuestsBooleanField label="present" sortable={false} />
        <DateField
          source="arrivedAt"
          showTime={true}
          emptyText="-"
          locales={locale}
          options={{ weekday: "short", hour: "2-digit", minute: "2-digit" }}
        />
        <ConditionalDeleteButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};

export const GuestEdit = () => {
  return (
    <Edit>
      <TabbedForm
        toolbar={
          <Toolbar>
            <SaveButton />
          </Toolbar>
        }
      >
        <FormTab label="General">
          <NumberInput disabled source="id" />
          <ReferenceInput source="guestlistId" reference="guestlists">
            <SelectInput optionText="name" validate={required()} disabled />
          </ReferenceInput>
          <TextInput source="name" validate={required()} />
          <TextInput
            source="code"
            helperText="The entrance code on the ticket"
          />
        </FormTab>
        <FormTab label="Additional Guests">
          <NumberInput
            source="additionalGuests"
            min={0}
            max={5}
            defaultValue={0}
            validate={ValidateAdditionalGuests}
            helperText="Number of additional guests (read as +1)"
          />
          <NumberInput
            source="attendedGuests"
            min={0}
            max={5}
            defaultValue={0}
            validate={ValidateAttendedGuests}
            helperText="Number of visitors that are present"
          />
        </FormTab>
        <FormTab label="Arrival">
          <ArrivedAtOrNullField />

          <TextInput
            source="arrivalNote"
            label="Note"
            helperText="A text that will be displayed when selecting this person."
          />
          <TextInput
            source="notifyOnArrivalEmail"
            validate={email()}
            label="Notify Email"
            helperText="Email to notify on arrival"
          />
        </FormTab>
      </TabbedForm>
    </Edit>
  );
};

const GuestCreateToolbar = ({ guestlistId, ...props }) => {
  const redirect = useRedirect();
  const { reset } = useFormContext();
  const notify = useNotify();

  const onSuccess = (data) => {
    notify(`Entry saved!`);
    reset();
    redirect(`/admin/guests/create?guestlist_id=${data.guestlistId}`);
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
GuestCreateToolbar.propTypes = {
  guestlistId: PropTypes.number.isRequired,
};

export const GuestCreate = (props) => {
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const guestlistId = parseInt(params.get("guestlist_id"), 10);

  return (
    <Create {...props} title="Create new List Entry">
      <SimpleForm
        defaultValues={{ guestlistId }}
        toolbar={<GuestCreateToolbar guestlistId={guestlistId} />}
      >
        <NumberInput disabled source="id" />
        <ReferenceInput source="guestlistId" reference="guestlists">
          <AutocompleteInput optionText="name" validate={required()} />
        </ReferenceInput>
        <TextInput source="name" validate={required()} />
        <TextInput source="code" helperText="The entrance code on the ticket" />
        <NumberInput
          source="additionalGuests"
          min={0}
          max={5}
          defaultValue={0}
          validate={ValidateAdditionalGuests}
          helperText="Number of additional guests (read as +1)"
        />
        <TextInput
          source="arrivalNote"
          label="Note"
          helperText="A text that will be displayed when selecting this person."
        />
        <TextInput
          source="notifyOnArrivalEmail"
          validate={email()}
          label="Notify Email"
          helperText="Email to notify on arrival"
        />
      </SimpleForm>
    </Create>
  );
};

export const GuestIcon = () => <PersonIcon />;
