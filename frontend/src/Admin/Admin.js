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
import {
  ListGroupCreate,
  ListGroupEdit,
  ListGroupIcon,
  ListGroupList,
} from "./compontent/ListGroup";
import PropTypes from "prop-types";

const AdminPanel = () => (
  <Admin
    layout={MyLayout}
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
      primaryText="List Entries"
      list={ListEntryList}
      edit={ListEntryEdit}
      create={ListEntryCreate}
      icon={ListEntryIcon}
    />
    <Resource
      name="listGroups"
      primaryText="List Groups"
      list={ListGroupList}
      edit={ListGroupEdit}
      create={ListGroupCreate}
      icon={ListGroupIcon}
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

const MyMenu = () => (
  <Menu>
    <Menu.ResourceItem name="products" />
    <Menu.ResourceItem name="purchases" />

    <div className="text-right mt-4 text-xs pr-1 font-bold tracking-wide border-b-2 border-b-pink-500 uppercase text-ellipsis">
      Guestlist
    </div>
    <Menu.ResourceItem name="lists" />
    <Menu.ResourceItem name="listGroups" primaryText="List Groups" />
    <Menu.ResourceItem name="listEntries" primaryText="List Entries" />

    <div className="text-right mt-4 text-xs pr-1 font-bold tracking-wide border-b-2 border-b-pink-500 uppercase text-ellipsis">
      Admin
    </div>
    <Menu.ResourceItem name="users" />
  </Menu>
);

const MyLayout = ({ children }) => <Layout menu={MyMenu}>{children}</Layout>;

MyLayout.propTypes = {
  children: PropTypes.object.isRequired,
};

export default AdminPanel;
