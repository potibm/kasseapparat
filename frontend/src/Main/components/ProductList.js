import React from "react";
import Product from "./Product";
import PropTypes from "prop-types";

function ProductList({ products, addToCart, currency }) {
  return (
    <div className="w-full">
      {products.map((product) => (
        <Product
          key={product.id}
          product={product}
          addToCart={addToCart}
          currency={currency}
        />
      ))}
    </div>
  );
}

ProductList.propTypes = {
  products: PropTypes.array.isRequired,
  addToCart: PropTypes.func.isRequired,
  currency: PropTypes.object.isRequired,
};

export default ProductList;
