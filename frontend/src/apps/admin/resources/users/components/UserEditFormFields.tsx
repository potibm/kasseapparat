import React from "react";
import {
  usePermissions,
  useGetIdentity,
  useRecordContext,
  NumberInput,
  SimpleForm,
  TextInput,
  PasswordInput,
  BooleanInput,
  required,
  minLength,
  email,
} from "react-admin";
import { UserRecord } from "../types";
import { UserEditToolbar } from "./UserEditToolbar";

const validatePasswords = (value: string, allValues: UserFormValues) => {
  if (allValues.password && !value) {
    return "Please confirm your new password";
  }
  if (value && value !== allValues.password) {
    return "The two passwords must match";
  }
  return undefined;
};

interface UserFormValues extends UserRecord {
  confirm_password?: string;
  password?: string;
}

export const UserEditFormFields: React.FC = () => {
  const record = useRecordContext<UserRecord>();
  const { permissions, isLoading: permissionsLoading } = usePermissions();
  const { data: identity, isLoading: identityLoading } = useGetIdentity();

  if (permissionsLoading || identityLoading || !record) {
    return <>Loading permissions...</>;
  }

  const isAdmin = permissions === "admin";
  const isCurrentUser = record.id === identity?.id;

  return (
    <SimpleForm toolbar={<UserEditToolbar />}>
      <NumberInput disabled source="id" />
      <TextInput source="username" validate={required()} fullWidth />

      <TextInput
        source="email"
        validate={[required(), email()]}
        disabled={!isAdmin && !isCurrentUser}
        fullWidth
      />

      {(isAdmin || isCurrentUser) && (
        <>
          <PasswordInput
            source="password"
            label="New Password"
            helperText="Leave empty to keep current password"
            validate={[minLength(8)]}
            fullWidth
          />
          <PasswordInput
            source="confirm_password"
            label="Confirm New Password"
            validate={validatePasswords}
            fullWidth
          />
        </>
      )}

      <BooleanInput
        source="admin"
        label="Administrator Permissions"
        disabled={!isAdmin}
      />
    </SimpleForm>
  );
};

export default UserEditFormFields;
