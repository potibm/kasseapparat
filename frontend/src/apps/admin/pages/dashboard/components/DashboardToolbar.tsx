import * as React from "react";
import {
  Box,
  Button,
  FormControlLabel,
  Switch,
  Typography,
} from "@mui/material";
import RefreshIcon from "@mui/icons-material/Refresh";
import { useConfig } from "@core/config/hooks/useConfig";

interface DashboardToolbarProps {
  enabled: boolean;
  setEnabled: (enabled: boolean) => void;
  refreshNow: () => void;
  lastRefresh: Date;
}

/** Auto-refresh switch, manual refresh and the time of the last update. */
const DashboardToolbar: React.FC<DashboardToolbarProps> = ({
  enabled,
  setEnabled,
  refreshNow,
  lastRefresh,
}) => {
  const { dateLocale } = useConfig();

  return (
    <Box
      sx={{
        mt: 2,
        display: "flex",
        alignItems: "center",
        gap: 2,
        flexWrap: "wrap",
      }}
    >
      <FormControlLabel
        control={
          <Switch
            checked={enabled}
            onChange={(event) => setEnabled(event.target.checked)}
          />
        }
        label="Auto-refresh every minute"
      />
      <Button
        size="small"
        startIcon={<RefreshIcon />}
        onClick={refreshNow}
        aria-label="Refresh now"
      >
        Refresh
      </Button>
      <Typography variant="body2" color="text.secondary">
        Updated{" "}
        {lastRefresh.toLocaleTimeString(dateLocale, {
          hour: "2-digit",
          minute: "2-digit",
        })}
      </Typography>
    </Box>
  );
};

export default DashboardToolbar;
