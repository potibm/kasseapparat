import { createContext } from "react";
import { AppConfig } from "../types/config.types";

export const ConfigContext = createContext<AppConfig | null>(null);

export default ConfigContext;
