import React, { useState } from "react";
import { Card } from "flowbite-react";
import { HiShoppingCart, HiUserAdd, HiQrcode } from "react-icons/hi";
import PropTypes from "prop-types";
import { useConfig } from "../../provider/ConfigProvider";
import GuestlistModal from "./GuestlistModal";
import MyButton from "./MyButton";
import ScannerModal from "./ScannerModal";

function Product({ product, addToCart, hasListItem }) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isScannerOpen, setIsScannerOpen] = useState(false);

  const handleAddToCart = () => {
    addToCart(product);
  };

  const handleShowGuestlist = () => {
    setIsModalOpen(true); 
  };

  const currency = useConfig().currency;

  return (
    <>
      <Card className="w-[24%] flex flex-col mb-2 mr-2">
        <h5 className="text-2xl text-left text-balance font-bold tracking-tight text-gray-900 dark:text-white">
          {product.name}
        </h5>
        <div className="flex-grow" style={{ flexGrow: 0.01 }}></div>{" "}
        <div className="flex items-center justify-between mt-auto">
          <p className="text-3xl font-bold text-gray-900 dark:text-white">
            {currency.format(product.price)}
          </p>
          <div className="flex">
            {product.lists.length > 0 && (
              <>
                <MyButton
                  onClick={() => setIsScannerOpen(true)}
                  aria-label="Open scanner"
                >
                  <HiQrcode />
                </MyButton>
                <ScannerModal 
                  isOpen={isScannerOpen} 
                  onClose={() => setIsScannerOpen(false)}
                  product={product}
                  addToCart={addToCart}
                  hasListItem={hasListItem}>
                </ScannerModal>

                <MyButton
                  onClick={handleShowGuestlist}
                  aria-label="Show guestlist"
                >
                  <HiUserAdd className="h-5 w-5" />
                </MyButton>
                <GuestlistModal
                  isOpen={isModalOpen}
                  onClose={() => setIsModalOpen(false)}
                  product={product}
                  addToCart={addToCart}
                  hasListItem={hasListItem}
                />
              </>
            )}
            {product.lists.length === 0 && (
              <MyButton onClick={handleAddToCart} aria-label="Add to cart">
                <HiShoppingCart className="h-5 w-5" />
              </MyButton>
            )}
          </div>
        </div>
      </Card>
      {product.wrapAfter && <div className="w-full"></div>}
    </>
  );
}

Product.propTypes = {
  product: PropTypes.object.isRequired,
  addToCart: PropTypes.func.isRequired,
  hasListItem: PropTypes.func.isRequired,
};

export default Product;
