import React from "react";
import { useConfig } from "@core/config/hooks/useConfig";
import {
  ConnectionStatus,
  useConnectionStatus,
} from "../hooks/useConnectionStatus";

const appearance: Record<
  ConnectionStatus,
  { label: string; dot: string; ring: string }
> = {
  online: { label: "Online", dot: "bg-green-600", ring: "ring-green-100" },
  slow: { label: "Slow", dot: "bg-amber-500", ring: "ring-amber-100" },
  offline: { label: "Offline", dot: "bg-red-600", ring: "ring-red-100" },
};

/**
 * Small dot in the POS header showing whether the backend is reachable. It
 * checks the public config endpoint, which also works behind the Vite dev proxy.
 */
export const ConnectionIndicator: React.FC = () => {
  const { apiBaseUrl } = useConfig();
  const status = useConnectionStatus(`${apiBaseUrl}/config`);
  const { label, dot, ring } = appearance[status];

  return (
    <div
      role="status"
      aria-live="polite"
      aria-label={`Connection: ${label}`}
      data-testid="connection-indicator"
      data-status={status}
      className="flex shrink-0 items-center gap-2 text-sm text-gray-700 dark:text-gray-200"
    >
      <span
        className={`h-2.5 w-2.5 rounded-full ring-4 dark:ring-transparent ${dot} ${ring}`}
        aria-hidden="true"
      />
      {label}
    </div>
  );
};

export default ConnectionIndicator;
