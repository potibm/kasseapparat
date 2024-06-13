import React from "react";
import Product from "./Product";
import PropTypes from "prop-types";

function ProductList({ products, addToCart }) {
  return (
    <div className="w-full">
      {products.map((product) => (
        <Product key={product.id} product={product} addToCart={addToCart} />
      ))}
    </div>
  );
}

ProductList.propTypes = {
  products: PropTypes.array.isRequired,
  addToCart: PropTypes.func.isRequired,
};

export default ProductList;
