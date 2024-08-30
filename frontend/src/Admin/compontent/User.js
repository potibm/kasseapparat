import React from "react";
import {
  usePermissions,
  useRecordContext,
  useGetIdentity,
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
  required,
  minLength,
  EmailField,
  email,
} from "react-admin";
import PersonIcon from "@mui/icons-material/Person";
import { UserFilters } from "./UserFilters";

const ConditionalDeleteButton = (props) => {
  const { permissions } = usePermissions();
  const { data: identity } = useGetIdentity();
  const record = useRecordContext(props);

  const isCurrentUser = record && record.id === identity.id;
  if (permissions === "admin" && !isCurrentUser) {
    return <DeleteButton {...props} />;
  }
  return null;
};

export const UserList = (props) => {
  return (
    <List sort={{ field: "id", order: "ASC" }} filters={UserFilters}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="username" />
        <EmailField source="email" />
        <BooleanField source="admin" />
        <ConditionalDeleteButton />
      </Datagrid>
    </List>
  );
};

export const UserEdit = () => {
  return (
    <Edit>
      <UserEditForm />
    </Edit>
  );
};

const UserEditForm = (props) => {
  const record = useRecordContext(props);
  const { data: permissions, isLoading: permissionsLoading } = usePermissions();
  const { data: identity, isLoading: identityLoading } = useGetIdentity();
  if (permissionsLoading || identityLoading) return <>Loading...</>;

  const equalToPassword = (value, allValues) => {
    if (value !== allValues.password) {
      return "The two passwords must match";
    }
  };

  const currentUserId = identity.id;
  const isCurrentUser = record && record.id === currentUserId;

  return (
    <SimpleForm
      toolbar={
        <Toolbar>
          <SaveButton />
        </Toolbar>
      }
    >
      <NumberInput disabled source="id" />
      <TextInput source="username" validate={required()} />
      <TextInput source="email" validate={[required(), email()]} />
      {(permissions === "admin" || isCurrentUser) && (
        <>
          <PasswordInput source="password" validate={[minLength(8)]} />
          <PasswordInput source="confirm_password" validate={equalToPassword} />
        </>
      )}
      <BooleanInput source="admin" disabled={permissions !== "admin"} />
    </SimpleForm>
  );
};

export const UserCreate = () => {
  const { isLoading, permissions } = usePermissions();
  if (isLoading) return <>Loading...</>;

  return (
    <Create title="Create new user">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="username" validate={required()} />
        <TextInput source="email" validate={[required(), email()]} />
        <BooleanInput source="admin" disabled={permissions !== "admin"} />
      </SimpleForm>
    </Create>
  );
};

export const UserIcon = () => <PersonIcon />;
