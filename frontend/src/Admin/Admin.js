// Admin.js
import React from "react";
import { Admin, Resource } from "react-admin";
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

const AdminPanel = () => (
  <Admin
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

export default AdminPanel;
