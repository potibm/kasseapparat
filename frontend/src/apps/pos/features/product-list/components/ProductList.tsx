import React from "react";
import Product from "./_internal/Product";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../utils/api.schemas";

interface ProductListProps {
  products: ProductType[] | null;
  addToCart: (
    product: ProductType,
    quantity: number,
    listItem: GuestType | null,
  ) => void;
  hasListItem: (guest: GuestType) => boolean;
  quantityByProductInCart: (product: ProductType) => number;
  addProductInterest: (product: ProductType) => Promise<void>;
}

const ProductList: React.FC<ProductListProps> = ({
  products,
  addToCart,
  hasListItem,
  quantityByProductInCart,
  addProductInterest,
}) => {
  if (!products || products.length === 0) {
    return <p>No products available.</p>;
  }

  return (
    <div className="flex flex-wrap -m-1.5">
      {products.map((product) => (
        <Product
          key={product.id}
          product={product}
          addToCart={addToCart}
          hasListItem={hasListItem}
          quantityByProductInCart={quantityByProductInCart}
          addProductInterest={addProductInterest}
        />
      ))}
    </div>
  );
};

export default ProductList;
