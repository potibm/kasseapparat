import React from "react";
import { Spinner } from "flowbite-react";
import Product from "./_internal/Product";
import ProductGroupMatrix from "./_internal/ProductGroupMatrix";
import { groupProductsBySection } from "../utils/productSections";
import { buildSectionBlocks } from "../utils/productBlocks";
import {
  Product as ProductType,
  Guest as GuestType,
  Section,
  ProductGroup,
} from "../../../api/schemas";

interface ProductListProps {
  products: ProductType[] | null;
  sections?: Section[];
  productGroups?: ProductGroup[];
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
  productGroups = [],
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
          <div className="flex flex-col gap-3">
            {buildSectionBlocks(section.products, productGroups).map((block) =>
              block.kind === "groups" ? (
                <ProductGroupMatrix
                  key={block.key}
                  groups={block.groups}
                  addToCart={addToCart}
                  hasListItem={hasListItem}
                  quantityByProductInCart={quantityByProductInCart}
                  addProductInterest={addProductInterest}
                />
              ) : (
                <div
                  key={block.key}
                  className="grid auto-rows-[7.5rem] grid-cols-[repeat(auto-fill,minmax(11rem,1fr))] gap-3"
                >
                  {block.products.map((product) => (
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
              ),
            )}
          </div>
        </section>
      ))}
    </div>
  );
};

export default ProductList;
