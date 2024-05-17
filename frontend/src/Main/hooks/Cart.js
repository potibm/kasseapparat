// Cart.js
export const addToCart = (cart, product) => {
  const existingProductIndex = cart.findIndex(item => item.id === product.id)
  if (existingProductIndex !== -1) {
    const updatedCart = [...cart]
    updatedCart[existingProductIndex].quantity++
    updatedCart[existingProductIndex].totalPrice = updatedCart[existingProductIndex].quantity * updatedCart[existingProductIndex].price
    return updatedCart
  } else {
    const updatedProduct = { ...product, quantity: 1, totalPrice: product.price }
    return [...cart, updatedProduct]
  }
}

export const removeFromCart = (cart, product) => {
  const existingProductIndex = cart.findIndex(item => item.id === product.id)

  if (existingProductIndex !== -1) {
    return [...cart.slice(0, existingProductIndex), ...cart.slice(existingProductIndex + 1)]
  }
  return cart
}

export const removeAllFromCart = () => {
  return []
}

export const checkoutCart = () => {
  return []
}
