import { createContext, use } from "react";

/**
 * A counter that goes up whenever the dashboard should reload its figures. The
 * stats hooks include it in their effect dependencies.
 */
export const DashboardRefreshContext = createContext(0);

export const useDashboardRefreshTick = () => use(DashboardRefreshContext);
