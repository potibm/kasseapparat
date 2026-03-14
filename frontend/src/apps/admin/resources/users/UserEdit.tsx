import React from "react";
import { Edit } from "react-admin";
import { UserEditFormFields } from "./components/UserEditFormFields";

export const UserEdit: React.FC = () => {
  return (
    <Edit mutationMode="pessimistic" title="Edit User">
      <UserEditFormFields />
    </Edit>
  );
};

export default UserEdit;
