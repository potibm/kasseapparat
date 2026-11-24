import React from "react";
import { DeleteWithConfirmButtonProps } from "react-admin";
import ConflictAwareDeleteButton from "@admin/ui/buttons/ConflictAwareDeleteButton";
import {
  PRODUCT_GROUP_DELETE_FAILED_MESSAGE,
  PRODUCT_GROUP_NOT_EMPTY_MESSAGE,
} from "./productGroup.utils";

/** Delete button for product groups: admins only, empty groups only (ADR 006). */
export const ProductGroupDeleteButton: React.FC<
  DeleteWithConfirmButtonProps
> = (props) => (
  <ConflictAwareDeleteButton
    {...props}
    conflictMessage={PRODUCT_GROUP_NOT_EMPTY_MESSAGE}
    fallbackMessage={PRODUCT_GROUP_DELETE_FAILED_MESSAGE}
  />
);

export default ProductGroupDeleteButton;
