import React from "react";
import {
  NumberInput,
  SimpleForm,
  TextInput,
  Create,
  BooleanInput,
  required,
  ReferenceInput,
  SelectInput,
} from "react-admin";

export const GuestlistCreate: React.FC = () => {
  return (
    <Create title="Create new guestlist">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="name" validate={required()} />
        <BooleanInput source="typeCode" />
        <ReferenceInput source="productId" reference="products">
          <SelectInput optionText="name" validate={required()} />
        </ReferenceInput>
      </SimpleForm>
    </Create>
  );
};

export default GuestlistCreate;
