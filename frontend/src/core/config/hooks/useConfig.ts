import { use } from "react";
import { ConfigContext } from "../context/ConfigContext";

export const useConfig = () => {
  const context = use(ConfigContext);
  if (!context) {
    throw new Error("useConfig must be used within a ConfigProvider");
  }
  return context;
};

export default useConfig;
