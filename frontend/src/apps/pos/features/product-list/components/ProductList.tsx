import React from "react";
import { Spinner } from "flowbite-react";
import Product from "./_internal/Product";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../api/schemas";

interface ProductListProps {
  products: ProductType[] | null;
  loading?: boolean;
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
  loading = false,
  addToCart,
  hasListItem,
  quantityByProductInCart,
  addProductInterest,
}) => {
  if (products === null) {
    return loading ? (
      <div className="flex items-center gap-3 p-4 dark:text-white">
        <Spinner size="md" /> Loading products…
      </div>
    ) : (
      <p className="p-4 dark:text-white">Products could not be loaded.</p>
    );
  }

  if (products.length === 0) {
    return <p className="p-4 dark:text-white">No products available.</p>;
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
