import React from "react";
import { Button, Card } from "flowbite-react";
import { HiShoppingCart } from "react-icons/hi";
import PropTypes from "prop-types";
import { useConfig } from "../../provider/ConfigProvider";

function Product({ product, addToCart }) {
  const handleAddToCart = () => {
    addToCart(product);
  };

  const currency = useConfig().currency;

  return (
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
        <Button onClick={handleAddToCart}>
          <HiShoppingCart className="h-5 w-5" />
        </Button>
      </div>
    </Card>
  );
}

Product.propTypes = {
  product: PropTypes.object.isRequired,
  addToCart: PropTypes.func.isRequired,
};

export default Product;
