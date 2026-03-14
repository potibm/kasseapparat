import React from "react";
import {
  SimpleForm,
  TextInput,
  Create,
  SaveButton,
  Toolbar,
} from "react-admin";
import { SumupReaderPairingCodeInput } from "./components/SumupReaderPairingCodeInput";

export const SumupReaderCreate: React.FC = () => {
  return (
    <Create title="Pair a SumUp Reader">
      <SimpleForm
        toolbar={
          <Toolbar>
            <SaveButton label="Pair" />
          </Toolbar>
        }
      >
        <SumupReaderPairingCodeInput />
        <TextInput source="name" />
      </SimpleForm>
    </Create>
  );
};
