import * as React from "react";
import { Title } from "react-admin";
import ProductStatsCard from "./components/ProductStatsCard";

const Dashboard: React.FC = () => {
  return (
    <>
      <Title title="Kasseapparat" />
      <ProductStatsCard />
    </>
  );
};

export default Dashboard;
