import React from "react";
import { Edit, SaveButton, SimpleForm, Toolbar } from "react-admin";
import ProductGroupDeleteButton from "./ProductGroupDeleteButton";
import ProductGroupFields from "./ProductGroupFields";
import { ProductGroup } from "./types";

export const ProductGroupEdit: React.FC = () => (
  <Edit<ProductGroup> mutationMode="pessimistic" title="Edit product group">
    <SimpleForm
      toolbar={
        <Toolbar sx={{ justifyContent: "space-between" }}>
          <SaveButton />
          <ProductGroupDeleteButton />
        </Toolbar>
      }
    >
      <ProductGroupFields />
    </SimpleForm>
  </Edit>
);

export default ProductGroupEdit;
