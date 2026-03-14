import React from "react";
import {
  usePermissions,
  DeleteButton,
  useRecordContext,
  useGetIdentity,
  DeleteWithConfirmButtonProps,
} from "react-admin";
import { UserRecord } from "../types";

export const ConditionalDeleteButton: React.FC<DeleteWithConfirmButtonProps> = (
  props,
) => {
  const record = useRecordContext<UserRecord>(props);
  const { permissions, isLoading: permissionsLoading } = usePermissions();
  const { data: identity, isLoading: identityLoading } = useGetIdentity();

  if (permissionsLoading || identityLoading) return null;

  const currentUserId = identity?.id;
  const isRecordOwner = record?.id === currentUserId;
  const isAdmin = permissions === "admin";

  if (isAdmin && !isRecordOwner) {
    return <DeleteButton {...props} />;
  }
  return null;
};

export default ConditionalDeleteButton;
