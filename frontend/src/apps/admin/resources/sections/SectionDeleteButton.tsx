import React from "react";
import { DeleteWithConfirmButtonProps } from "react-admin";
import ConflictAwareDeleteButton from "@admin/ui/buttons/ConflictAwareDeleteButton";
import {
  SECTION_DELETE_FAILED_MESSAGE,
  SECTION_NOT_EMPTY_MESSAGE,
} from "./section.utils";

/**
 * Delete button for sections: only admins see it (ADR 006), and a section that still
 * contains products is reported in words instead of a bare "Conflict".
 */
export const SectionDeleteButton: React.FC<DeleteWithConfirmButtonProps> = (
  props,
) => (
  <ConflictAwareDeleteButton
    {...props}
    conflictMessage={SECTION_NOT_EMPTY_MESSAGE}
    fallbackMessage={SECTION_DELETE_FAILED_MESSAGE}
  />
);

export default SectionDeleteButton;
