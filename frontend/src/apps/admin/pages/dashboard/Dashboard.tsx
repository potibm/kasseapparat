import * as React from "react";
import { Title } from "react-admin";
import ProductStatsCard from "./components/ProductStatsCard";
import PaymentMethodStatsCard from "./components/PaymentMethodStatsCard";

const Dashboard: React.FC = () => {
  return (
    <>
      <Title title="Dashboard" />
      <ProductStatsCard />
      <PaymentMethodStatsCard />
    </>
  );
};

export default Dashboard;
