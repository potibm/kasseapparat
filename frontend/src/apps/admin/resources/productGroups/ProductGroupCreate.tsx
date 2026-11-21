import React from "react";
import { Create, SimpleForm } from "react-admin";
import ProductGroupFields from "./ProductGroupFields";
import { ProductGroup } from "./types";

export const ProductGroupCreate: React.FC = () => (
  <Create<ProductGroup> title="Create new product group" redirect="list">
    <SimpleForm defaultValues={{ name: "", pos: 0 }}>
      <ProductGroupFields preselectSection />
    </SimpleForm>
  </Create>
);

export default ProductGroupCreate;
