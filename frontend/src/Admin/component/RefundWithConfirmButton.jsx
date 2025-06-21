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

const RefundWithConfirmButton = () => {
  const [open, setOpen] = useState(false);
  const record = useRecordContext();
  const dataProvider = useDataProvider();
  const notify = useNotify();
  const refresh = useRefresh();

  const handleClick = () => setOpen(true);
  const handleClose = () => setOpen(false);

  const handleConfirm = async () => {
    try {
      await dataProvider.refund("purchases", { id: record.id });
      notify("Refund successful", { type: "success" });
      refresh();
    } catch (error) {
      console.error("Refund error:", error);
      notify("Refund failed", { type: "error" });
    } finally {
      setOpen(false);
    }
  };

  return (
    <>
      <Button
        label="Refund"
        onClick={(event) => {
          event.stopPropagation(); // prevent rowClick from triggering
          handleClick();
        }}
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
