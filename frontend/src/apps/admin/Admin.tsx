import React, { useMemo } from "react";
import { Admin, Resource } from "react-admin";
import { createDataProvider } from "./providers/dataProvider";
import products from "./resources/products";
import productInterests from "./resources/product-interests";
import purchases from "./resources/purchases";
import guests from "./resources/guests";
import guestlists from "./resources/guestlists";
import sumupreaders from "./resources/sumup-readers";
import sumuptransactions from "./resources/sumup-transactions";
import Dashboard from "./pages/dashboard/Dashboard";
import { MyTheme, MyDarkTheme } from "./layouts/MyTheme";
import { MyLayout } from "./layouts/MyLayout";
import { createAuthProvider } from "@core/auth/authProvider";
import useConfig from "@core/config/hooks/useConfig";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3100";
const authProvider = createAuthProvider(API_HOST);

const AdminPanel: React.FC = () => {
  const { apiBaseUrl } = useConfig();

  const dataProvider = useMemo(() => {
    return createDataProvider(apiBaseUrl);
  }, [apiBaseUrl]);

  return (
    <Admin
      layout={MyLayout}
      loginPage={false}
      theme={MyTheme}
      darkTheme={MyDarkTheme}
      dashboard={Dashboard}
      dataProvider={dataProvider}
      authProvider={authProvider}
      title="Kasseapparat Admin"
      basename="/admin"
    >
      <Resource name="products" {...products} />
      <Resource name="productInterests" {...productInterests} />
      <Resource name="guestlists" {...guestlists} />
      <Resource name="guests" {...guests} />
      <Resource name="purchases" {...purchases} />
      <Resource
        name="sumupReaders"
        {...sumupreaders}
        options={{ label: "Readers" }}
      />
      <Resource
        name="sumupTransactions"
        {...sumuptransactions}
        options={{ label: "Transactions" }}
      />
    </Admin>
  );
};

export default AdminPanel;
