import React, { useState } from "react";
import { Button, Card } from "flowbite-react";
import { HiShoppingCart, HiUserAdd } from "react-icons/hi";
import PropTypes from "prop-types";
import { useConfig } from "../../provider/ConfigProvider";
import GuestlistModal from "./GuestlistModal";

function Product({ product, addToCart, hasListItem }) {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleAddToCart = () => {
    addToCart(product);
  };

  const handleShowGuestlist = () => {
    setIsModalOpen(true); // Öffnen des Modals
  };

  const currency = useConfig().currency;

  return (
    <>
      <Card
        className={`w-1/4 mr-1.5 mb-1.5 float-left ${product.wrapAfter ? "float-none" : ""}`}
      >
        <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          {product.name}
        </h5>
        <div className="flex items-center justify-between">
          <span className="text-3xl font-bold text-gray-900 dark:text-white">
            {currency.format(product.price)}
          </span>
          {product.lists.length > 0 && (
            <Button onClick={handleShowGuestlist}>
              <HiUserAdd className="h-5 w-5" />
            </Button>
          )}
          <Button onClick={handleAddToCart}>
            <HiShoppingCart className="h-5 w-5" />
          </Button>
        </div>
      </Card>
      <GuestlistModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        product={product}
        addToCart={addToCart}
        hasListItem={hasListItem}
      />
    </>
  );
}

Product.propTypes = {
  product: PropTypes.object.isRequired,
  addToCart: PropTypes.func.isRequired,
  hasListItem: PropTypes.func.isRequired,
};

export default Product;
