import React from "react";
import {
  NumberInput,
  Edit,
  TextInput,
  SaveButton,
  Toolbar,
  ReferenceInput,
  SelectInput,
  required,
  TabbedForm,
  FormTab,
  email,
} from "react-admin";
import { ArrivedAtOrNullField } from "./components/ArrivedAtOrNullField";
import { AdditionalGuestsInput } from "./components/AdditionalGuestsInput";
import { AttendedGuestsInput } from "./components/AttendedGuestsInput";

export const GuestEdit: React.FC = () => {
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
          <AdditionalGuestsInput />
          <AttendedGuestsInput />
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
