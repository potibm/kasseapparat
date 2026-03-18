import React from "react";
import { useRecordContext, Button, RaRecord } from "react-admin";
import { Box, Tooltip } from "@mui/material";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import TouchAppIcon from "@mui/icons-material/TouchApp";

interface SumupReaderSelectionFieldProps {
  selectedReaderId: string | undefined;
  onSelect: (id: string) => void;
}

interface SumupReaderRecord extends RaRecord {
  status: string;
}

export const SumupReaderSelectionField: React.FC<
  SumupReaderSelectionFieldProps
> = ({ selectedReaderId, onSelect }) => {
  const record = useRecordContext<SumupReaderRecord>();

  if (record?.status !== "paired") return null;

  const isCurrent = selectedReaderId === record.id;

  if (isCurrent) {
    return (
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          color: "success.main",
          fontWeight: 500,
        }}
      >
        <CheckCircleOutlineIcon sx={{ mr: 1, fontSize: "1.2rem" }} />
        Selected
      </Box>
    );
  }

  return (
    <Tooltip title="Assign this reader to this device">
      <span>
        <Button
          onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
            e.stopPropagation();
            console.log("Selecting reader with id:", record.id);
            onSelect(String(record.id));
          }}
          size="small"
          startIcon={<TouchAppIcon />}
          variant="text"
          color="primary"
          label="Use this reader"
        />
      </span>
    </Tooltip>
  );
};

export default SumupReaderSelectionField;
