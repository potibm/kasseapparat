import React from "react";
import {
  usePermissions,
  BooleanInput,
  SimpleForm,
  TextInput,
  Create,
  required,
  email,
} from "react-admin";

export const UserCreate: React.FC = () => {
  const { isLoading, permissions } = usePermissions();
  if (isLoading) return <>Loading...</>;

  return (
    <Create title="Create new user">
      <SimpleForm>
        <TextInput source="username" validate={required()} />
        <TextInput source="email" validate={[required(), email()]} />
        <BooleanInput source="admin" disabled={permissions !== "admin"} />
      </SimpleForm>
    </Create>
  );
};

export default UserCreate;
