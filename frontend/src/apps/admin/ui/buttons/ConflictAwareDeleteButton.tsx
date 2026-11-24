import React from "react";
import { DeleteWithConfirmButtonProps, useNotify } from "react-admin";
import ConditionalDeleteOnAdminButton from "./ConditionalDeleteOnAdminButton";
import { deleteErrorMessage } from "./deleteErrorMessage";

export type ConflictAwareDeleteButtonProps = DeleteWithConfirmButtonProps & {
  /** Shown when the backend refuses the delete with 409. */
  conflictMessage: string;
  /** Shown when the error carries no message at all. */
  fallbackMessage: string;
};

/**
 * Delete button for things that can only be deleted when empty: only admins see it,
 * and a refused delete is reported in words instead of a bare "Conflict".
 */
export const ConflictAwareDeleteButton: React.FC<
  ConflictAwareDeleteButtonProps
> = ({ conflictMessage, fallbackMessage, ...props }) => {
  const notify = useNotify();

  return (
    <ConditionalDeleteOnAdminButton
      {...props}
      mutationOptions={{
        onError: (error: unknown) => {
          notify(deleteErrorMessage(error, conflictMessage, fallbackMessage), {
            type: "error",
          });
        },
      }}
    />
  );
};

export default ConflictAwareDeleteButton;
