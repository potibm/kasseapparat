import React from "react";
import { Admin, Resource, Layout, LayoutProps } from "react-admin";
import dataProvider from "./providers/data-provider";
import authProvider from "./providers/auth-provider";
import products from "./resources/products";
import productInterests from "./resources/product-interests";
import purchases from "./resources/purchases";
import guests from "./resources/guests";
import guestlists from "./resources/guestlists";
import sumupreaders from "./resources/sumup-readers";
import sumuptransactions from "./resources/sumup-transactions";
import users from "./resources/users";
import Dashboard from "./pages/dashboard/Dashboard";
import { Menu } from "./layouts/Menu";

const MyLayout = (props: LayoutProps) => <Layout {...props} menu={Menu} />;

const AdminPanel: React.FC = () => (
  <Admin
    layout={MyLayout}
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
    <Resource name="users" {...users} />
  </Admin>
);

export default AdminPanel;
