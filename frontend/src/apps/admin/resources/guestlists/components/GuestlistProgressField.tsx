import React from "react";
import { Box, LinearProgress, Typography } from "@mui/material";
import { useRecordContext } from "react-admin";

interface GuestlistStatsRecord {
  id: number;
  entries?: number;
  expectedGuests?: number;
  arrivedGuests?: number;
}

/** How many of the expected guests (entries plus companions) have arrived. */
export const GuestlistProgressField: React.FC<{ label?: string }> = () => {
  const list = useRecordContext<GuestlistStatsRecord>();
  if (!list) return null;

  const entries = list.entries ?? 0;
  const expected = list.expectedGuests ?? 0;
  const arrived = list.arrivedGuests ?? 0;

  if (entries === 0) {
    return (
      <Typography variant="body2" color="text.secondary">
        No entries
      </Typography>
    );
  }

  const percent = expected > 0 ? Math.min(100, (arrived / expected) * 100) : 0;

  return (
    <Box sx={{ minWidth: 140 }} data-testid="guestlist-progress">
      <Typography variant="body2">
        {arrived} / {expected} arrived
      </Typography>
      <LinearProgress
        variant="determinate"
        value={percent}
        aria-label={`${arrived} of ${expected} guests arrived`}
      />
      <Typography variant="caption" color="text.secondary">
        {entries} {entries === 1 ? "entry" : "entries"}
      </Typography>
    </Box>
  );
};

export default GuestlistProgressField;
