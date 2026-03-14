import React from "react";
import {
  NumberInput,
  Edit,
  SimpleForm,
  TextInput,
  BooleanInput,
  SaveButton,
  Toolbar,
  required,
  ReferenceInput,
  SelectInput,
} from "react-admin";

export const GuestlistEdit: React.FC = () => {
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
        <BooleanInput source="typeCode" />
        <ReferenceInput source="productId" reference="products">
          <SelectInput optionText="name" validate={required()} />
        </ReferenceInput>
      </SimpleForm>
    </Edit>
  );
};

export default GuestlistEdit;
