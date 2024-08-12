import React from "react";
import Product from "./Product";
import PropTypes from "prop-types";

function ProductList({
  products,
  addToCart,
  hasListItem,
  quantityByProductInCart,
}) {
  return (
    <div className="flex flex-wrap -m-1.5">
      {products.map((product) => (
        <Product
          key={product.id}
          product={product}
          addToCart={addToCart}
          hasListItem={hasListItem}
          quantityByProductInCart={quantityByProductInCart}
        />
      ))}
    </div>
  );
}

ProductList.propTypes = {
  products: PropTypes.array.isRequired,
  addToCart: PropTypes.func.isRequired,
  hasListItem: PropTypes.func.isRequired,
  quantityByProductInCart: PropTypes.func.isRequired,
};

export default ProductList;
