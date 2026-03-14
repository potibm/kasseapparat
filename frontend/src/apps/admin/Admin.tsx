import React from "react";
import { Admin, Resource, Layout, LayoutProps } from "react-admin";
import dataProvider from "./providers/data-provider";
import authProvider from "./providers/auth-provider";
import products from "./resources/products";
import productInterests from "./resources/productinterests";
import purchases from "./resources/purchases";
import guests from "./resources/guests";
import { UserCreate, UserEdit, UserIcon, UserList } from "./components/User";
import {
  GuestlistCreate,
  GuestlistEdit,
  GuestlistIcon,
  GuestlistList,
} from "./components/Guestlist";
import {
  SumupReaderCreate,
  SumupReaderList,
  SumupReaderIcon,
} from "./components/SumupReader";
import {
  SumupTransactionIcon,
  SumupTransactionList,
  SumupTransactionShow,
} from "./components/SumupTransaction";
import Dashboard from "./components/Dashboard";
import { Menu } from "./layouts/Menu";

const MyLayout = (props: LayoutProps) => <Layout {...props} menu={Menu} />;

const AdminPanel: React.FC = () => (
  <Admin
    layout={MyLayout}
    dashboard={Dashboard}
    dataProvider={dataProvider}
    authProvider={authProvider}
    basename="/admin"
  >
    <Resource name="products" {...products} />
    <Resource name="productInterests" {...productInterests} />
    <Resource
      name="guestlists"
      list={GuestlistList}
      edit={GuestlistEdit}
      create={GuestlistCreate}
      icon={GuestlistIcon}
    />
    <Resource name="guests" {...guests} />
    <Resource name="purchases" {...purchases} />
    <Resource
      name="sumupReaders"
      list={SumupReaderList}
      create={SumupReaderCreate}
      icon={SumupReaderIcon}
      options={{ label: "Readers" }}
    />
    <Resource
      name="sumupTransactions"
      list={SumupTransactionList}
      show={SumupTransactionShow}
      icon={SumupTransactionIcon}
      options={{ label: "Transactions" }}
    />
    <Resource
      name="users"
      list={UserList}
      edit={UserEdit}
      create={UserCreate}
      icon={UserIcon}
    />
  </Admin>
);

export default AdminPanel;
