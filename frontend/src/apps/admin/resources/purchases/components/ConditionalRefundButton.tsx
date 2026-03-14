import React from "react";
import {
  DeleteWithConfirmButtonProps,
  usePermissions,
  useRecordContext,
} from "react-admin";
import RefundWithConfirmButton from "./RefundWithConfirmButton";

export const ConditionalRefundButton: React.FC<DeleteWithConfirmButtonProps> = (
  props,
) => {
  const { permissions } = usePermissions();
  const record = useRecordContext(props);

  if (
    permissions === "admin" &&
    record.status === "confirmed" &&
    record.paymentMethod === "SUMUP"
  ) {
    return <RefundWithConfirmButton />;
  }
  return null;
};

export default ConditionalRefundButton;
