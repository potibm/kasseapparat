import React from "react";
import {
  SaveButton,
  Toolbar,
  useRedirect,
  useNotify,
  ToolbarProps,
} from "react-admin";
import { useFormContext } from "react-hook-form";

export const GuestCreateToolbar: React.FC<ToolbarProps> = ({ ...props }) => {
  const redirect = useRedirect();
  const { reset } = useFormContext();
  const notify = useNotify();

  const handleSuccess = (data: any) => {
    notify(`Guest entry saved!`, { type: "info" });
    reset();
    console.log("Created guest with data:", data);

    const finalId = data.guestlistId;
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
