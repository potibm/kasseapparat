import React, { useState } from "react";
import { List, ListProps } from "react-admin";
import { useConfig } from "@core/config/providers/ConfigProvider";
import {
  getCurrentReaderId,
  setCurrentReaderId,
  clearCurrentReaderId,
} from "@core/localstorage/helper/reader";
import SumupReaderListContent from "./components/SumupReaderListContent";
import SumupReaderListActions from "./components/SumupReaderListActions";
import { Box, Typography } from "@mui/material";

export const SumupReaderList: React.FC<ListProps> = (props) => {
  const { sumupEnabled } = useConfig();
  const [selectedReaderId, setSelectedReaderId] = useState<string | undefined>(
    () => {
      return getCurrentReaderId() || undefined;
    },
  );

  const handleReaderSelect = (id: string) => {
    setCurrentReaderId(id);
    setSelectedReaderId(id);
  };

  const handleClear = () => {
    clearCurrentReaderId();
    setSelectedReaderId(undefined);
  };

  if (!sumupEnabled) {
    return (
      <Box p={2}>
        <Typography variant="h5">SumUp Readers</Typography>
        <Typography>
          SumUp integration is not enabled. Please enable it in the
          configuration.
        </Typography>
      </Box>
    );
  }

  return (
    <List title="SumUp Readers" {...props} actions={<SumupReaderListActions />}>
      <SumupReaderListContent
        selectedReaderId={selectedReaderId}
        onClear={handleClear}
        onSelect={handleReaderSelect}
      />
    </List>
  );
};
