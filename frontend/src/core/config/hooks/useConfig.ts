import { use } from "react";
import { ConfigContext } from "../context/ConfigContext";

export const useConfig = () => {
  const context = use(ConfigContext);
  if (!context) {
    throw new Error(
      "useAppConfig must be used within a ConfigContext.Provider",
    );
  }
  return context;
};

export default useConfig;
