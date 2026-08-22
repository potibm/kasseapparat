import * as React from "react";
import { Title } from "react-admin";
import ProductStatsCard from "./components/ProductStatsCard";
import PaymentMethodStatsCard from "./components/PaymentMethodStatsCard";
import HourlyRevenueChart from "./components/HourlyRevenueChart";
import HourlyQuantityChart from "./components/HourlyQuantityChart";

const Dashboard: React.FC = () => {
  return (
    <>
      <Title title="Dashboard" />
      <ProductStatsCard />
      <PaymentMethodStatsCard />
      <HourlyRevenueChart />
      <HourlyQuantityChart />
    </>
  );
};

export default Dashboard;
