import React from "react";
import { useConfig } from "@core/config/hooks/useConfig";
import VariantButton from "./VariantButton";
import { ProductTapProps } from "./useProductTap";
import { GroupWithProducts } from "../../utils/productBlocks";
import { Product as DisplayProduct } from "@pos/api/schemas";

type MatrixProps = Omit<ProductTapProps, "product"> & {
  groups: GroupWithProducts<DisplayProduct>[];
};

/** The single price of a group, or null when its variants differ in price. */
const sharedPrice = (products: DisplayProduct[]) => {
  const [first, ...rest] = products;
  if (!first) return null;

  return rest.every((product) => product.grossPrice.equals(first.grossPrice))
    ? first.grossPrice
    : null;
};

/**
 * The size matrix: one row per product group, its name (and its price, when all
 * variants cost the same) at the start, then one button per variant (ADR 006).
 */
const ProductGroupMatrix: React.FC<MatrixProps> = ({ groups, ...tapProps }) => {
  const { currency } = useConfig();

  return (
    <div className="flex flex-col gap-2.5 rounded-xl bg-white p-3 shadow-sm dark:bg-gray-800">
      {groups.map(({ group, products }) => {
        const price = sharedPrice(products);

        return (
          <fieldset
            key={group.id}
            aria-label={group.name}
            data-testid={"product-group-" + group.id}
            className="flex items-center gap-2 border-0 p-0 m-0 min-w-0"
          >
            <div className="w-24 shrink-0">
              <div className="text-base leading-tight font-bold text-gray-900 dark:text-white">
                {group.name}
              </div>
              {price && (
                <div className="text-sm text-gray-600 dark:text-gray-300">
                  {currency.format(price.toNumber())}
                </div>
              )}
            </div>
            <div className="grid flex-1 grid-cols-[repeat(auto-fill,minmax(4.75rem,1fr))] gap-2">
              {products.map((product) => (
                <VariantButton
                  key={product.id}
                  {...tapProps}
                  product={product}
                  showPrice={price === null}
                />
              ))}
            </div>
          </fieldset>
        );
      })}
    </div>
  );
};

export default ProductGroupMatrix;
