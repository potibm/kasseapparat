import React from "react";
import {
  Button,
  Confirm,
  useNotify,
  useRefresh,
  useRecordContext,
  useDataProvider,
} from "react-admin";
import { useState } from "react";
import CurrencyExchangeIcon from "@mui/icons-material/CurrencyExchange";

const RefundWithConfirmButton: React.FC = () => {
  const [open, setOpen] = useState(false);
  const record = useRecordContext();
  const dataProvider = useDataProvider();
  const notify = useNotify();
  const refresh = useRefresh();

  if (!record || record.status === "refunded") return null;

  const handleClick = (event: React.MouseEvent) => {
    event.stopPropagation(); // Prevent Datagrid rowClick
    setOpen(true);
  };

  const handleClose = () => setOpen(false);

  const handleConfirm = async () => {
    try {
      await dataProvider.refund("purchases", { id: record.id });
      notify("Refund successful", { type: "success" });
      refresh();
    } catch (error: any) {
      console.error("Refund error:", error);
      notify(error.message || "Refund failed", { type: "error" });
    } finally {
      setOpen(false);
    }
  };

  return (
    <>
      <Button
        label="Refund"
        onClick={handleClick}
        startIcon={<CurrencyExchangeIcon />}
      />
      <Confirm
        isOpen={open}
        title="Confirm refund"
        content="Are you sure you want to refund this transaction?"
        onConfirm={handleConfirm}
        onClose={handleClose}
      />
    </>
  );
};

export default RefundWithConfirmButton;
