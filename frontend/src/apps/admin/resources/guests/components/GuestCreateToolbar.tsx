import React from "react";
import {
  SaveButton,
  Toolbar,
  useRedirect,
  useNotify,
  ToolbarProps,
  RaRecord,
} from "react-admin";
import { useFormContext } from "react-hook-form";
import { GuestRecord } from "../types";

export const GuestCreateToolbar: React.FC<ToolbarProps> = ({ ...props }) => {
  const redirect = useRedirect();
  const { reset } = useFormContext();
  const notify = useNotify();

  const handleSuccess = (data: RaRecord) => {
    notify(`Guest entry saved!`, { type: "info" });
    reset();

    const guest = data as GuestRecord;
    const finalId = guest.guestlistId;
    if (finalId) {
      redirect(`/admin/guests/create?guestlist_id=${finalId}`);
    }
  };

  return (
    <Toolbar {...props}>
      <SaveButton />
      <SaveButton
        type="button"
        label="Save and Add another"
        mutationOptions={{ onSuccess: handleSuccess }}
        sx={{ ml: 1 }}
      />
    </Toolbar>
  );
};

export default GuestCreateToolbar;
