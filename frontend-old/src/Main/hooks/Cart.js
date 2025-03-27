// Cart.js
export const addToCart = (cart, product, count = 1, listItem = null) => {
  if (listItem) {
    listItem.attendedGuests = count;
  }

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
    updatedCart[existingProductIndex].totalNetPrice = updatedCart[
      existingProductIndex
    ].netPrice.mul(updatedCart[existingProductIndex].quantity);
    updatedCart[existingProductIndex].totalGrossPrice = updatedCart[
      existingProductIndex
    ].grossPrice.mul(updatedCart[existingProductIndex].quantity);
    updatedCart[existingProductIndex].totalVatAmount = updatedCart[
      existingProductIndex
    ].vatAmount.mul(updatedCart[existingProductIndex].quantity);

    // add list item to cart
    if (listItem) {
      updatedCart[existingProductIndex].listItems.push(listItem);
    }

    return updatedCart;
  } else {
    const updatedProduct = {
      ...product,
      quantity: count,
      totalNetPrice: product.netPrice,
      totalGrossPrice: product.grossPrice,
      totalVatAmount: product.vatAmount,
      listItems: [],
    };
    updatedProduct.unitsSold += count;

    if (listItem) {
      updatedProduct.listItems.push(listItem);
    }

    return [...cart, updatedProduct];
  }
};

export const containsListItemID = (cart, listItemID) => {
  // check if the cart (withing the list of products in the cart) contains a list item with the given ID
  return cart.some((product) =>
    product.listItems.some((listItem) => listItem.id === listItemID),
  );
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

export const getCartProductQuantity = (cart, product) => {
  const existingProductIndex = cart.findIndex((item) => item.id === product.id);
  if (existingProductIndex !== -1) {
    return cart[existingProductIndex].quantity;
  } else {
    return 0;
  }
};
