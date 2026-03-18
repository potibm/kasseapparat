import React from "react";
import {
  usePermissions,
  DeleteButton,
  useRecordContext,
  useGetIdentity,
  DeleteWithConfirmButtonProps,
  RaRecord,
} from "react-admin";

interface OwnableRecord extends RaRecord {
  createdById?: string | number;
}

export const ConditionalDeleteOnOwnershipButton: React.FC<
  DeleteWithConfirmButtonProps
> = (props) => {
  const record = useRecordContext<OwnableRecord>(props);
  const { permissions, isLoading: permissionsLoading } = usePermissions();
  const { data: identity, isLoading: identityLoading } = useGetIdentity();

  if (permissionsLoading || identityLoading) return null;

  const currentUserId = identity?.id;
  const isOwner =
    record?.createdById &&
    currentUserId &&
    record.createdById === currentUserId;
  const isAdmin = permissions === "admin";

  if (isAdmin || isOwner) {
    return <DeleteButton {...props} />;
  }
  return null;
};

export default ConditionalDeleteOnOwnershipButton;
