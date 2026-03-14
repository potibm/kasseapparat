import React from "react";
import {
  DeleteWithConfirmButton,
  DeleteWithConfirmButtonProps,
  usePermissions,
} from "react-admin";

/**
 * A delete button that only renders if the user has admin permissions.
 * All props are passed through to the underlying DeleteWithConfirmButton.
 */
export const ConditionalDeleteOnAdminButton: React.FC<
  DeleteWithConfirmButtonProps
> = (props) => {
  const { permissions, isLoading } = usePermissions();

  if (isLoading) return null;

  if (permissions === "admin") {
    return <DeleteWithConfirmButton mutationMode="pessimistic" {...props} />;
  }

  return null;
};

export default ConditionalDeleteOnAdminButton;
