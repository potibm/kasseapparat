import { useCallback, useEffect, useRef, useState } from "react";

export type ConnectionStatus = "online" | "slow" | "offline";

export interface ConnectionStatusOptions {
  /** How often the backend is checked while the POS is open. */
  intervalMs?: number;
  /** A check that takes longer than this counts as "slow". */
  slowThresholdMs?: number;
  /** A check that takes longer than this is aborted and counts as "offline". */
  timeoutMs?: number;
}

export const DEFAULT_INTERVAL_MS = 15_000;
export const DEFAULT_SLOW_THRESHOLD_MS = 2_000;
export const DEFAULT_TIMEOUT_MS = 5_000;

/**
 * Tells whether the backend is reachable, so a flaky Wi-Fi shows up before a
 * payment fails rather than during it.
 *
 * Checks `url` on mount, every `intervalMs`, when the tab becomes visible again
 * and when the browser reports it is back online. The browser's "offline" event
 * switches to "offline" right away.
 */
export const useConnectionStatus = (
  url: string,
  {
    intervalMs = DEFAULT_INTERVAL_MS,
    slowThresholdMs = DEFAULT_SLOW_THRESHOLD_MS,
    timeoutMs = DEFAULT_TIMEOUT_MS,
  }: ConnectionStatusOptions = {},
): ConnectionStatus => {
  const [status, setStatus] = useState<ConnectionStatus>("online");
  const inFlightRef = useRef(false);

  const check = useCallback(async () => {
    if (inFlightRef.current) return;

    if (globalThis.navigator?.onLine === false) {
      setStatus("offline");
      return;
    }

    inFlightRef.current = true;
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeoutMs);
    const startedAt = Date.now();

    try {
      const response = await fetch(url, {
        method: "GET",
        cache: "no-store",
        signal: controller.signal,
      });
      const duration = Date.now() - startedAt;

      if (!response.ok) {
        setStatus("offline");
      } else {
        setStatus(duration > slowThresholdMs ? "slow" : "online");
      }
    } catch {
      setStatus("offline");
    } finally {
      clearTimeout(timeoutId);
      inFlightRef.current = false;
    }
  }, [url, slowThresholdMs, timeoutMs]);

  useEffect(() => {
    void check();
    const intervalId = setInterval(() => void check(), intervalMs);

    const handleOnline = () => void check();
    const handleOffline = () => setStatus("offline");
    const handleVisibility = () => {
      if (document.visibilityState === "visible") void check();
    };

    globalThis.addEventListener("online", handleOnline);
    globalThis.addEventListener("offline", handleOffline);
    document.addEventListener("visibilitychange", handleVisibility);

    return () => {
      clearInterval(intervalId);
      globalThis.removeEventListener("online", handleOnline);
      globalThis.removeEventListener("offline", handleOffline);
      document.removeEventListener("visibilitychange", handleVisibility);
    };
  }, [check, intervalMs]);

  return status;
};

export default useConnectionStatus;
