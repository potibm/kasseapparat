import React, { useState } from "react";
import { Badge, Card } from "flowbite-react";
import { HiShoppingCart, HiUserAdd, HiOutlineThumbUp } from "react-icons/hi";
import PropTypes from "prop-types";
import { useConfig } from "../../provider/ConfigProvider";
import GuestlistModal from "./GuestlistModal";
import MyButton from "./MyButton";
import ProductInterestModal from "./ProductInterestModal";

function Product({ product, addToCart, hasListItem, quantityByProductInCart }) {
  const [isGuestListModalOpen, setIsGuestListModalOpen] = useState(false);
  const [isPIModalOpen, setIsPIModalOpen] = useState(false);


  const handleAddToCart = () => {
    addToCart(product);
  };

  const handleShowGuestlist = () => {
    setIsGuestListModalOpen(true);
  };

  const handleHideGuestlist = () => {
    setIsGuestListModalOpen(false);
  };
  const currency = useConfig().currency;

  const compactCardTheme = {
    root: {
      children: "flex h-full flex-col justify-center gap-2 p-4",
    },
  };

  const handleCardClick = () => {
    if (product.soldOut) {
      console.log("ioen")
      setIsPIModalOpen(true);
    } else if (product.lists.length > 0) {
      handleShowGuestlist();
    } else {
      handleAddToCart();
    }
  };

  return (
    <>
      <Card
        theme={compactCardTheme}
        className="w-[22%] flex flex-col mb-5 mr-5 relative"
        href="#"
        onClick={handleCardClick}
      >
        {product.soldOut && (
          <Badge className="absolute top-2 right-2" color="gray">
            Sold Out
          </Badge>
        )}
        <div className="flex items-center justify-between mt-auto">
          <h5 className="text-1xl text-left text-balance font-bold tracking-tight text-gray-900 dark:text-white">
            {product.name}
          </h5>
          {!product.soldOut && product.totalStock > 0 && (
            <div className="text-sm">
              {product.unitsSold + quantityByProductInCart(product)} /{" "}
              {product.totalStock}
            </div>
          )}
        </div>
        <div className="flex-grow" style={{ flexGrow: 0.01 }}></div>{" "}
        <div className="flex items-center justify-between mt-auto">
          <p
            className={`text-2xl font-bold ${product.soldOut ? "text-gray-400" : "text-gray-900 dark:text-white"}`}
          >
            {currency.format(product.price)}
          </p>
          <div className="flex">
            {product.soldOut && (
              <MyButton aria-label="Register interest">
                <HiOutlineThumbUp className="h-5 w-5" />
              </MyButton>
            )}
            {!product.soldOut && product.lists.length > 0 && (
              <>
                <MyButton aria-label="Show guestlist">
                  <HiUserAdd className="h-5 w-5" />
                </MyButton>
              </>
            )}
            {!product.soldOut && product.lists.length === 0 && (
              <MyButton aria-label="Add to cart">
                <HiShoppingCart className="h-5 w-5" />
              </MyButton>
            )}
          </div>
        </div>
      </Card>
      {product.wrapAfter && <div className="w-full"></div>}
      {!product.soldOut && product.lists.length > 0 && (
        <GuestlistModal
                    isOpen={isGuestListModalOpen}
                    onClose={handleHideGuestlist}
                    product={product}
                    addToCart={addToCart}
                    hasListItem={hasListItem}
                  />
      )}
      {product.soldOut && (
        <ProductInterestModal show={isPIModalOpen} onClose={() => setIsPIModalOpen(false)} product={product} />
      )}
    </>
  );
}

Product.propTypes = {
  product: PropTypes.object.isRequired,
  addToCart: PropTypes.func.isRequired,
  hasListItem: PropTypes.func.isRequired,
  quantityByProductInCart: PropTypes.func.isRequired,
};

export default Product;
