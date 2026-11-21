import React from "react";
import { Spinner } from "flowbite-react";
import Product from "./_internal/Product";
import { DisplayProduct } from "../utils/vanishedProducts";
import { groupProductsBySection } from "../utils/productSections";
import {
  Product as ProductType,
  Guest as GuestType,
  Section,
} from "../../../api/schemas";

interface ProductListProps {
  products: DisplayProduct[] | null;
  sections?: Section[];
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
  sections = [],
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
    <div className="flex flex-col gap-4 p-1" data-testid="product-grid">
      {groupProductsBySection(products, sections).map((section) => (
        <section
          key={section.key}
          aria-label={section.name || undefined}
          data-testid={section.key}
        >
          {section.name && (
            <h2 className="mb-1.5 px-1 text-xs font-semibold tracking-wider text-slate-500 uppercase dark:text-gray-400">
              {section.name}
            </h2>
          )}
          <div className="grid auto-rows-[7.5rem] grid-cols-[repeat(auto-fill,minmax(11rem,1fr))] gap-3">
            {section.products.map((product) => (
              <Product
                key={product.id}
                product={product}
                addToCart={addToCart}
                hasListItem={hasListItem}
                quantityByProductInCart={quantityByProductInCart}
                addProductInterest={addProductInterest}
                vanished={product.vanished}
              />
            ))}
          </div>
        </section>
      ))}
    </div>
  );
};

export default ProductList;
