// Admin.js
import React from "react";
import {
  Admin,
  MenuItemLink,
  Layout,
  Menu,
  Resource,
  useSidebarState,
  useResourceDefinitions,
} from "react-admin";
import dataProvider from "./dataProvider";
import authProvider from "./authProvider";
import {
  ProductList,
  ProductEdit,
  ProductIcon,
  ProductCreate,
} from "./component/Product";
import { PurchaseList, PurchaseShow } from "./component/Purchase";
import { UserCreate, UserEdit, UserIcon, UserList } from "./component/User";
import {
  GuestlistCreate,
  GuestlistEdit,
  GuestlistIcon,
  GuestlistList,
} from "./component/Guestlist";
import {
  GuestCreate,
  GuestEdit,
  GuestIcon,
  GuestList,
} from "./component/Guest";
import {
  ProductInterestList,
  ProductInterestIcon,
} from "./component/ProductInterest";
import {
  SumupReaderCreate,
  SumupReaderList,
  SumupReaderIcon,
} from "./component/SumupReader";
import PropTypes from "prop-types";
import Dashboard from "./component/Dashboard";
import { useConfig } from "../provider/ConfigProvider";

const AdminPanel = () => (
  <Admin
    layout={MyLayout}
    dashboard={Dashboard}
    dataProvider={dataProvider}
    authProvider={authProvider}
    basename="/admin"
  >
    <Resource
      name="products"
      list={ProductList}
      edit={ProductEdit}
      create={ProductCreate}
      icon={ProductIcon}
    />
    <Resource
      name="productInterests"
      list={ProductInterestList}
      icon={ProductInterestIcon}
      options={{ label: "Product Interest" }}
    />
    <Resource
      name="guestlists"
      list={GuestlistList}
      edit={GuestlistEdit}
      create={GuestlistCreate}
      icon={GuestlistIcon}
    />
    <Resource
      name="guests"
      list={GuestList}
      edit={GuestEdit}
      create={GuestCreate}
      icon={GuestIcon}
      options={{ label: "Guests" }}
    />
    <Resource name="purchases" list={PurchaseList} show={PurchaseShow} />
    <Resource
      name="sumupReaders"
      list={SumupReaderList}
      create={SumupReaderCreate}
      icon={SumupReaderIcon}
      options={{ label: "Readers" }}
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

const MyMenuDivider = ({ name }) => (
  <div className="text-right mt-4 text-xs pr-1 font-bold tracking-wide border-b-2 border-b-pink-500 uppercase text-ellipsis">
    {name}
  </div>
);

MyMenuDivider.propTypes = {
  name: PropTypes.string.isRequired,
};

const MyMenuMessage = ({ message }) => (
  <div className="text-left mt-4 text-xs font-bold tracking-wide border-t-2 border-b-2 border-t-red-700 border-b-red-700 bg-info pl-5 pr-3 uppercase text-ellipsis">
    {message}
  </div>
);

MyMenuMessage.propTypes = {
  message: PropTypes.string.isRequired,
};

const MyMenu = () => {
  const environmentMessage = useConfig().environmentMessage;
  const sumupEnabled = useConfig().sumupEnabled;
  const resources = useResourceDefinitions();

  const [isSidebarOpen] = useSidebarState();

  return (
    <Menu>
      <Menu.DashboardItem />

      <MyMenuDivider name="POS" />
      <Menu.ResourceItem name="products" />
      <Menu.ResourceItem name="purchases" />
      <Menu.ResourceItem name="productInterests" />

      <MyMenuDivider name="Guestlist" />
      <Menu.ResourceItem name="guestlists" />
      <Menu.ResourceItem name="guests" />

      <MyMenuDivider name="Sumup" />
      {sumupEnabled ? (
        <Menu.ResourceItem name="sumupReaders" />
      ) : (
        <MenuItemLink
          to="#"
          primaryText={
            resources.sumupReaders?.options?.label || "Sumup Readers"
          }
          leftIcon={<SumupReaderIcon />}
          disabled
        />
      )}

      <MyMenuDivider name="Admin" />
      <Menu.ResourceItem name="users" />

      {environmentMessage && isSidebarOpen && (
        <MyMenuMessage message={environmentMessage} />
      )}
    </Menu>
  );
};

const MyLayout = ({ children }) => <Layout menu={MyMenu}>{children}</Layout>;

MyLayout.propTypes = {
  children: PropTypes.object.isRequired,
};

export default AdminPanel;
