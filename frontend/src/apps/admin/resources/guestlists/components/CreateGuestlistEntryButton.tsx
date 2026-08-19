import React from "react";
import { useRecordContext, Button, RaRecord, Identifier } from "react-admin";
import { useNavigate } from "react-router";
import PersonAddIcon from "@mui/icons-material/PersonAdd";

interface GuestlistRecord extends RaRecord {
  id: Identifier;
}

const CreateGuestlistEntryButton: React.FC = () => {
  const record = useRecordContext<GuestlistRecord>();
  const navigate = useNavigate();

  const handleCreateEntry = (guestlistId: Identifier) => {
    navigate(`/admin/guests/create?guestlist_id=${guestlistId}`);
  };

  if (!record) return null;

  return (
    <Button
      label="Add Guest"
      onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        e.stopPropagation();
        handleCreateEntry(record?.id);
      }}
      sx={{
        "& .RaButton-label": {
          display: { xs: "none", sm: "inline" },
        },
      }}
    >
      <PersonAddIcon />
    </Button>
  );
};

export default CreateGuestlistEntryButton;
