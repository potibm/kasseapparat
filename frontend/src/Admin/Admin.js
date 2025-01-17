// Admin.js
import React from "react";
import { Admin, Layout, Menu, Resource, useSidebarState } from "react-admin";
import dataProvider from "./dataProvider";
import authProvider from "./authProvider";
import {
  ProductList,
  ProductEdit,
  ProductIcon,
  ProductCreate,
} from "./compontent/Product";
import { PurchaseList, PurchaseShow } from "./compontent/Purchase";
import { UserCreate, UserEdit, UserIcon, UserList } from "./compontent/User";
import {
  GuestlistCreate,
  GuestlistEdit,
  GuestlistIcon,
  GuestlistList,
} from "./compontent/Guestlist";
import {
  ListEntryCreate,
  ListEntryEdit,
  ListEntryIcon,
  ListEntryList,
} from "./compontent/ListEntry";
import {
  ProductInterestList,
  ProductInterestIcon,
} from "./compontent/ProductInterest";
import PropTypes from "prop-types";
import Dashboard from "./compontent/Dashboard";
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
      name="listEntries"
      list={ListEntryList}
      edit={ListEntryEdit}
      create={ListEntryCreate}
      icon={ListEntryIcon}
      options={{ label: "List Entries" }}
    />
    <Resource name="purchases" list={PurchaseList} show={PurchaseShow} />
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
      <Menu.ResourceItem name="listEntries" primaryText="List Entries" />

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
