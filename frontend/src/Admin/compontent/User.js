import React from "react";
import {
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
  PasswordInput,
  BooleanField,
  BooleanInput,
} from "react-admin";
import PersonIcon from "@mui/icons-material/Person";

export const UserList = () => {
  return (
    <List sort={{ field: "id", order: "ASC" }}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="username" />
        <BooleanField source="admin" />
        <DeleteButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};

export const UserEdit = () => {
  const equalToPassword = (value, allValues) => {
    if (value !== allValues.password) {
      return "The two passwords must match";
    }
  };

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
        <TextInput source="username" />
        <PasswordInput source="password" />
        <PasswordInput source="confirm_password" validate={equalToPassword} />
        <BooleanInput source="admin" />
      </SimpleForm>
    </Edit>
  );
};

export const UserCreate = () => {
  const equalToPassword = (value, allValues) => {
    if (value !== allValues.password) {
      return "The two passwords must match";
    }
  };

  return (
    <Create title="Create new user">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="username" />
        <PasswordInput source="password" />
        <PasswordInput source="confirm_password" validate={equalToPassword} />
        <BooleanInput source="admin" />
      </SimpleForm>
    </Create>
  );
};

export const UserIcon = () => <PersonIcon />;
