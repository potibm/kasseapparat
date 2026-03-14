import React from "react";
import {
  SimpleForm,
  TextInput,
  Create,
  ReferenceInput,
  required,
  AutocompleteInput,
  email,
  CreateProps,
} from "react-admin";
import { useLocation } from "react-router";
import { GuestCreateToolbar } from "./components/GuestCreateToolbar";
import { AdditionalGuestsInput } from "./components/AdditionalGuestsInput";

export const GuestCreate: React.FC<CreateProps> = (props) => {
  const location = useLocation();
  const params = new URLSearchParams(location.search);

  const guestlistIdParam = params.get("guestlist_id");
  const guestlistId = guestlistIdParam
    ? Number.parseInt(guestlistIdParam, 10)
    : undefined;

  return (
    <Create {...props} title="Create new List Entry">
      <SimpleForm
        defaultValues={{
          guestlistId: guestlistId,
        }}
        toolbar={<GuestCreateToolbar />}
      >
        <ReferenceInput source="guestlistId" reference="guestlists">
          <AutocompleteInput optionText="name" validate={required()} />
        </ReferenceInput>
        <TextInput source="name" validate={required()} />
        <TextInput source="code" helperText="The entrance code on the ticket" />
        <AdditionalGuestsInput />
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
