// Admin.js
import React from 'react';
import { Admin, Resource} from 'react-admin';
import dataProvider from './dataProvider';
import authProvider from './authProvider';
import { ProductList, ProductEdit, ProductIcon, ProductCreate } from './compontent/Product';
import { PurchaseList, PurchaseShow } from './compontent/Purchase';

const AdminPanel = () => (
  <Admin dataProvider={dataProvider} authProvider={authProvider} basename='/admin'>
    <Resource name="products" list={ProductList} edit={ProductEdit} create={ProductCreate} icon={ProductIcon} />
    <Resource name="purchases" list={PurchaseList} show={PurchaseShow} />
  </Admin>
);

export default AdminPanel;