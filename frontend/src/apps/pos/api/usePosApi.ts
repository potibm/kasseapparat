import { use, useMemo } from "react";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { createPosApiClient } from "./factory";

export const usePosApi = () => {
  const config = use(ConfigContext);
  if (!config) throw new Error("ConfigContext is missing");

  return useMemo(
    () => createPosApiClient(config.apiBaseUrl),
    [config.apiBaseUrl],
  );
};
