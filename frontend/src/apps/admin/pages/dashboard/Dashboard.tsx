import * as React from "react";
import { Title } from "react-admin";
import ProductStatsCard from "./components/ProductStatsCard";

const Dashboard: React.FC = () => {
  return (
    <>
      <Title title="Dashboard" />
      <ProductStatsCard />
    </>
  );
};

export default Dashboard;
