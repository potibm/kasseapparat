import React from "react";
import { DeleteWithConfirmButtonProps, useNotify } from "react-admin";
import ConditionalDeleteOnAdminButton from "@admin/ui/buttons/ConditionalDeleteOnAdminButton";
import { sectionDeleteErrorMessage } from "./section.utils";

/**
 * Delete button for sections: only admins see it (ADR 006), and a section that still
 * contains products is reported in words instead of a bare "Conflict".
 */
export const SectionDeleteButton: React.FC<DeleteWithConfirmButtonProps> = (
  props,
) => {
  const notify = useNotify();

  return (
    <ConditionalDeleteOnAdminButton
      {...props}
      mutationOptions={{
        onError: (error: unknown) => {
          notify(sectionDeleteErrorMessage(error), { type: "error" });
        },
      }}
    />
  );
};

export default SectionDeleteButton;
