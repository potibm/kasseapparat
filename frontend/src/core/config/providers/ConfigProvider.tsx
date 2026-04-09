import React, { useState, useEffect, ReactNode } from "react";
import MinimalBaseCard from "../../../apps/pos/components/MinimalBaseCard";
import { Alert } from "flowbite-react";
import { HiInformationCircle } from "react-icons/hi";
import { AppConfig } from "../types/config.types";
import { ConfigSchema } from "../schemas/config.schemas";
import { transformConfig } from "../utils/config.transform";
import { createLogger } from "@core/logger/logger";
import { ConfigContext } from "../context/ConfigContext";

const log = createLogger("Config");
const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3001";

interface ConfigProviderProps {
  children: ReactNode;
  fetchUrl?: string;
}

export const ConfigProvider: React.FC<ConfigProviderProps> = ({
  children,
  fetchUrl = `${API_HOST}/api/v2/config`,
}) => {
  const [config, setConfig] = useState<AppConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadConfig = async () => {
      try {
        const response = await fetch(fetchUrl);
        if (!response.ok)
          throw new Error(`HTTP error! status: ${response.status}`);

        const json = await response.json();

        const parsedData = ConfigSchema.parse(json);

        const finalConfig = transformConfig(parsedData, API_HOST);

        setConfig(finalConfig);
        log.debug("Config loaded successfully", finalConfig);
      } catch (err) {
        log.error("Config loading failed", err);
        setError(err instanceof Error ? err.message : "Unknown error");
      } finally {
        setLoading(false);
      }
    };

    loadConfig();
  }, [fetchUrl]);

  if (loading) {
    return <div>⏳ Loading config...</div>;
  }

  if (error) {
    return (
      <MinimalBaseCard>
        Error loading config:
        <Alert
          color="failure"
          icon={HiInformationCircle}
          className="mt-3 font-mono"
        >
          {error}
        </Alert>
      </MinimalBaseCard>
    );
  }

  return <ConfigContext value={config}>{children}</ConfigContext>;
};

export default ConfigProvider;
