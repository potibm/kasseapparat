import { useCallback, useEffect, useState } from "react";

/** How often the dashboard reloads while auto-refresh is on. */
export const AUTO_REFRESH_MS = 60_000;

/**
 * Drives the dashboard refresh: a tick that goes up every minute while
 * auto-refresh is on and the page is visible, or when asked to refresh now.
 */
export const useAutoRefresh = (intervalMs: number = AUTO_REFRESH_MS) => {
  const [tick, setTick] = useState(0);
  const [enabled, setEnabled] = useState(true);
  const [lastRefresh, setLastRefresh] = useState(() => new Date());

  const refreshNow = useCallback(() => {
    setTick((current) => current + 1);
    setLastRefresh(new Date());
  }, []);

  useEffect(() => {
    if (!enabled) return;

    const handle = setInterval(() => {
      if (document.visibilityState === "visible") refreshNow();
    }, intervalMs);

    return () => clearInterval(handle);
  }, [enabled, intervalMs, refreshNow]);

  return { tick, enabled, setEnabled, refreshNow, lastRefresh };
};
