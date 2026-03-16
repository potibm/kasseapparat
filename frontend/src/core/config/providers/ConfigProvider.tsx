import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import MinimalBaseCard from "../../../apps/pos/components/MinimalBaseCard";
import { Alert, Spinner } from "flowbite-react";
import { HiInformationCircle } from "react-icons/hi";
import { AppConfig } from "../types/config.types";
import { ConfigSchema, RawConfig } from "../schemas/config.schemas";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3001";

export const transformConfig = (
  rawData: RawConfig,
  apiHost: string,
): AppConfig => {
  const sumupEnabled = rawData.paymentMethods.some((m) => m.code === "SUMUP");
  const websocketHost = apiHost.replace(/^http/, "ws");

  const currencyOptions: Intl.NumberFormatOptions = {
    style: "currency",
    currency: rawData.currencyCode,
    minimumFractionDigits: rawData.fractionDigitsMin,
    maximumFractionDigits: rawData.fractionDigitsMax,
  };

  return {
    ...rawData,
    apiHost,
    websocketHost,
    sumupEnabled,
    currencyOptions,
    currency: new Intl.NumberFormat(rawData.currencyLocale, currencyOptions),
  };
};

export const ConfigContext = createContext<AppConfig | null>(null);

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
      } catch (err) {
        console.error("Config loading failed:", err);
        setError(err instanceof Error ? err.message : "Unknown error");
      } finally {
        setLoading(false);
      }
    };

    loadConfig();
  }, [fetchUrl]);

  if (loading) {
    return (
      <div>⏳ Loading config...</div>
    );
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

  return (
    <ConfigContext.Provider value={config}>{children}</ConfigContext.Provider>
  );
};

export const useConfig = () => {
  const context = useContext(ConfigContext);
  if (!context) {
    throw new Error("useConfig must be used within a ConfigProvider");
  }
  return context;
};

export default ConfigProvider;
