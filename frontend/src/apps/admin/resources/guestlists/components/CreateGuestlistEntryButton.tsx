import React from "react";
import {
  useRecordContext,
  Button,
  ButtonProps,
  RaRecord,
  Identifier,
} from "react-admin";
import { useNavigate } from "react-router";
import PersonAddIcon from "@mui/icons-material/PersonAdd";

interface GuestlistRecord extends RaRecord {
  id: Identifier;
}

const CreateGuestlistEntryButton: React.FC<ButtonProps> = (props) => {
  const record = useRecordContext<GuestlistRecord>(props);
  const navigate = useNavigate();

  const handleCreateEntry = (guestlistId: Identifier) => {
    navigate(`/admin/guests/create?guestlist_id=${guestlistId}`);
  };

  if (!record) return null;

  return (
    <Button
      {...props}
      label="Add Guest"
      onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        e.stopPropagation();
        handleCreateEntry(record?.id);
      }}
      startIcon={<PersonAddIcon />}
    />
  );
};

export default CreateGuestlistEntryButton;
