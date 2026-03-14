import React from "react";
import {
  Datagrid,
  TextField,
  DeleteWithConfirmButton,
  Button,
  useListContext,
} from "react-admin";
import { Box } from "@mui/material";
import LinkOffIcon from "@mui/icons-material/LinkOff";
import { useTheme } from "@mui/material/styles";
import { SumupReaderStatusField } from "./SumupReaderStatusField";
import { SumupReaderSelectionField } from "./SumupReaderSelectionField";

interface SumupReaderListContentProps {
  selectedReaderId: string | undefined;
  onClear: () => void;
  onSelect: (id: string) => void;
}

export const SumupReaderListContent: React.FC<SumupReaderListContentProps> = ({
  selectedReaderId,
  onClear,
  onSelect,
}) => {
  const { data: readers = [], isLoading } = useListContext();
  const theme = useTheme();

  const isSelectedReaderMissing =
    !isLoading &&
    selectedReaderId &&
    !readers.some((r) => String(r.id) === selectedReaderId);

  return (
    <>
      {isSelectedReaderMissing && (
        <Box
          sx={{
            p: 2,
            mb: 2,
            backgroundColor: theme.palette.warning.light,
            color: theme.palette.warning.contrastText,
            borderRadius: 1,
            border: `1px solid ${theme.palette.warning.main}`,
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <Box>
            The previously selected reader (<strong>{selectedReaderId}</strong>)
            is no longer available.
          </Box>
          <Button
            variant="contained"
            color="warning"
            size="small"
            onClick={onClear}
          >
            Clear selection
          </Button>
        </Box>
      )}

      <Datagrid
        rowClick={false}
        bulkActionButtons={false}
        rowSx={(record) =>
          String(record.id) === selectedReaderId
            ? { backgroundColor: theme.palette.action.selected }
            : {}
        }
      >
        <TextField source="id" sortable={false} />
        <TextField source="name" sortable={false} />
        <SumupReaderStatusField
          source="status"
          label="Status"
          sortable={false}
        />
        <TextField
          source="deviceIdentifier"
          label="Identifier"
          sortable={false}
        />
        <TextField source="deviceModel" label="Model" sortable={false} />

        <SumupReaderSelectionField
          selectedReaderId={selectedReaderId}
          onSelect={onSelect}
          label="Action"
        />

        <DeleteWithConfirmButton
          label="Unpair"
          confirmTitle="Unpair device"
          confirmContent="Are you sure you want to unpair this SumUp reader?"
          mutationMode="pessimistic"
          icon={<LinkOffIcon />}
          sx={{ color: "error.main" }}
        />
      </Datagrid>
    </>
  );
};

export default SumupReaderListContent;
