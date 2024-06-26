// Cart.js
export const addToCart = (cart, product, count = 1, listItem = null) => {
  const existingProductIndex = cart.findIndex((item) => item.id === product.id);
  if (existingProductIndex !== -1) {
    // if list item (identified by id) exists in the existing product listItems: do nothing and return
    if (
      listItem &&
      cart[existingProductIndex].listItems.find(
        (item) => item.id === listItem.id,
      )
    ) {
      return cart;
    }

    const updatedCart = [...cart];
    updatedCart[existingProductIndex].quantity += count;
    updatedCart[existingProductIndex].totalPrice =
      updatedCart[existingProductIndex].quantity *
      updatedCart[existingProductIndex].price;

    // add list item to
    if (listItem) {
      updatedCart[existingProductIndex].listItems.push(listItem);
    }

    return updatedCart;
  } else {
    const updatedProduct = {
      ...product,
      quantity: count,
      totalPrice: product.price,
      listItems: [],
    };

    if (listItem) {
      updatedProduct.listItems.push(listItem);
    }

    return [...cart, updatedProduct];
  }
};

export const removeFromCart = (cart, product) => {
  const existingProductIndex = cart.findIndex((item) => item.id === product.id);

  if (existingProductIndex !== -1) {
    return [
      ...cart.slice(0, existingProductIndex),
      ...cart.slice(existingProductIndex + 1),
    ];
  }
  return cart;
};

export const removeAllFromCart = () => {
  return [];
};

export const checkoutCart = () => {
  return [];
};
