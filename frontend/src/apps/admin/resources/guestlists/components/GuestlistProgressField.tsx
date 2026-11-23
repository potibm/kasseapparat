import React from "react";
import { Box, LinearProgress, Typography } from "@mui/material";
import { useRecordContext } from "react-admin";

interface GuestlistStatsRecord {
  id: number;
  entries?: number;
  arrivedEntries?: number;
  arrivedGuests?: number;
}

/**
 * How many list entries have arrived, and – for information – how many people
 * that was in total. Additional guests are optional, so they only count once
 * they actually came.
 */
export const GuestlistProgressField: React.FC<{ label?: string }> = () => {
  const list = useRecordContext<GuestlistStatsRecord>();
  if (!list) return null;

  const entries = list.entries ?? 0;
  const arrived = list.arrivedEntries ?? 0;
  const people = list.arrivedGuests ?? 0;

  if (entries === 0) {
    return (
      <Typography variant="body2" color="text.secondary">
        No entries
      </Typography>
    );
  }

  const percent = Math.min(100, (arrived / entries) * 100);

  return (
    <Box sx={{ minWidth: 160 }} data-testid="guestlist-progress">
      <Typography variant="body2">
        {arrived} / {entries} arrived
        <Typography component="span" variant="body2" color="text.secondary">
          {" "}
          · {people} {people === 1 ? "person" : "people"}
        </Typography>
      </Typography>
      <LinearProgress
        variant="determinate"
        value={percent}
        aria-label={`${arrived} of ${entries} entries arrived`}
      />
    </Box>
  );
};

export default GuestlistProgressField;
