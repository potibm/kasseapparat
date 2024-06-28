// Admin.js
import React from "react";
import { Admin, Layout, Menu, Resource } from "react-admin";
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
import { ListCreate, ListEdit, ListIcon, ListList } from "./compontent/List";
import {
  ListEntryCreate,
  ListEntryEdit,
  ListEntryIcon,
  ListEntryList,
} from "./compontent/ListEntry";
import PropTypes from "prop-types";
import Dashboard from "./compontent/Dashboard";

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
      name="lists"
      list={ListList}
      edit={ListEdit}
      create={ListCreate}
      icon={ListIcon}
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

const MyMenu = () => (
  <Menu>
    <Menu.DashboardItem />

    <MyMenuDivider name="POS" />
    <Menu.ResourceItem name="products" />
    <Menu.ResourceItem name="purchases" />

    <MyMenuDivider name="Guestlist" />
    <Menu.ResourceItem name="lists" />
    <Menu.ResourceItem name="listEntries" primaryText="List Entries" />

    <MyMenuDivider name="Admin" />
    <Menu.ResourceItem name="users" />
  </Menu>
);

const MyLayout = ({ children }) => <Layout menu={MyMenu}>{children}</Layout>;

MyLayout.propTypes = {
  children: PropTypes.object.isRequired,
};

export default AdminPanel;
