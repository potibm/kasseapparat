import React, { useState } from "react";
import { Badge, Card } from "flowbite-react";
import { HiShoppingCart, HiUserAdd, HiOutlineThumbUp } from "react-icons/hi";
import { useConfig } from "@core/config/hooks/useConfig";
import {
  availableStock,
  canAddMore,
  isUnavailable,
} from "@pos/features/stock/utils/stock";
import GuestlistModal from "../../../guestlist/components/GuestlistModal";
import Button from "../../../../components/Button";
import ProductInterestModal from "./ProductInterestModal";
import {
  Product as ProductType,
  Guest as GuestType,
} from "../../../../api/schemas";

interface ProductProps {
  product: ProductType;
  addToCart: (
    product: ProductType,
    quantity: number,
    listItem: GuestType | null,
  ) => void;
  hasListItem: (guest: GuestType) => boolean;
  quantityByProductInCart: (product: ProductType) => number;
  addProductInterest: (product: ProductType) => Promise<void>;
}

const Product: React.FC<ProductProps> = ({
  product,
  addToCart,
  hasListItem,
  quantityByProductInCart,
  addProductInterest,
}) => {
  const [isGuestListModalOpen, setIsGuestListModalOpen] = useState(false);
  const [isPIModalOpen, setIsPIModalOpen] = useState(false);
  const { currency, outOfStockBehavior } = useConfig();

  // A product with limited stock that has run out cannot be added any more. Under the
  // "ignore" behaviour it stays sellable, so only a product an operator marked sold out
  // is treated as unavailable there.
  const soldOut = isUnavailable(product, outOfStockBehavior);
  const quantityInCart = quantityByProductInCart(product);
  const stockLeft = availableStock(product, quantityInCart);
  const hasGuestlist = product.guestlists && product.guestlists.length > 0;

  // Whether another unit still fits. Shared with the cart, so this never reads as
  // enabled while the cart would refuse the same quantity. Only consulted for a plain
  // product: a sold-out product still offers "Register interest", and a guestlist
  // product still opens its modal.
  const canAdd = canAddMore(product, quantityInCart, outOfStockBehavior);

  const getActionButton = () => {
    if (soldOut) {
      return (
        <Button aria-label={"Register interest in " + product.name}>
          <HiOutlineThumbUp className="h-5 w-5" />
        </Button>
      );
    } else if (hasGuestlist) {
      return (
        <Button aria-label={"Show guestlist for " + product.name}>
          <HiUserAdd className="h-5 w-5" />
        </Button>
      );
    } else if (!canAdd) {
      return (
        <Button aria-label={"No " + product.name + " left"} disabled>
          <HiShoppingCart className="h-5 w-5" />
        </Button>
      );
    } else {
      return (
        <Button
          aria-label={
            "Add " +
            product.name +
            " for " +
            currency.format(product.grossPrice.toNumber()) +
            " to cart"
          }
        >
          <HiShoppingCart className="h-5 w-5" />
        </Button>
      );
    }
  };

  const handleCardClick = (e: React.MouseEvent) => {
    e.preventDefault();
    if (soldOut) {
      setIsPIModalOpen(true);
    } else if (hasGuestlist) {
      setIsGuestListModalOpen(true);
    } else if (canAdd) {
      addToCart(product, 1, null);
    }
  };

  const compactCardTheme = {
    root: {
      children: "flex h-full flex-col justify-center gap-2 p-4",
    },
  };

  return (
    <>
      <Card
        theme={compactCardTheme}
        data-testid={"product-card-" + product.id}
        className="w-[22%] flex flex-col mb-5 mr-5 relative cursor-pointer"
        onClick={handleCardClick}
      >
        {soldOut && (
          <Badge className="absolute top-2 right-2" color="gray">
            Sold Out ({product.soldOutRequestCount})
          </Badge>
        )}

        <div className="flex items-center justify-between mt-auto">
          <h5
            className={`text-1xl text-left text-balance font-bold tracking-tight ${
              soldOut ? "text-gray-400" : "text-gray-900 dark:text-gray-200"
            }`}
          >
            {product.name}
          </h5>

          {!soldOut && product.totalStock > 0 && (
            <div className="text-sm dark:text-white">
              {stockLeft >= 0 && <span>{stockLeft} / </span>}
              {product.totalStock}
            </div>
          )}
        </div>

        <div className="flex items-center justify-between mt-auto">
          <p
            className={`text-2xl font-bold ${
              soldOut ? "text-gray-400" : "text-gray-900 dark:text-white"
            }`}
          >
            {currency.format(product.grossPrice.toNumber())}
          </p>

          <div className="flex">{getActionButton()}</div>
        </div>
      </Card>

      {product.wrapAfter && <div className="w-full"></div>}

      {!soldOut && hasGuestlist && (
        <GuestlistModal
          isOpen={isGuestListModalOpen}
          onClose={() => setIsGuestListModalOpen(false)}
          product={product}
          addToCart={addToCart}
          hasListItem={hasListItem}
        />
      )}

      {soldOut && (
        <ProductInterestModal
          show={isPIModalOpen}
          onClose={() => setIsPIModalOpen(false)}
          product={product}
          addProductInterest={addProductInterest}
        />
      )}
    </>
  );
};

export default Product;
