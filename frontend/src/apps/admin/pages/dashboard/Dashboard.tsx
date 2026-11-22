import * as React from "react";
import { Title } from "react-admin";
import ProductStatsCard from "./components/ProductStatsCard";
import PaymentMethodStatsCard from "./components/PaymentMethodStatsCard";
import HourlyRevenueChart from "./components/HourlyRevenueChart";
import HourlyQuantityChart from "./components/HourlyQuantityChart";
import KpiCards from "./components/KpiCards";
import DashboardToolbar from "./components/DashboardToolbar";
import { DashboardRefreshContext } from "./context/DashboardRefreshContext";
import { useAutoRefresh } from "./hooks/useAutoRefresh";

const Dashboard: React.FC = () => {
  const { tick, enabled, setEnabled, refreshNow, lastRefresh } =
    useAutoRefresh();

  return (
    <DashboardRefreshContext.Provider value={tick}>
      <Title title="Dashboard" />
      <DashboardToolbar
        enabled={enabled}
        setEnabled={setEnabled}
        refreshNow={refreshNow}
        lastRefresh={lastRefresh}
      />
      <KpiCards />
      <ProductStatsCard />
      <PaymentMethodStatsCard />
      <HourlyRevenueChart />
      <HourlyQuantityChart />
    </DashboardRefreshContext.Provider>
  );
};

export default Dashboard;
