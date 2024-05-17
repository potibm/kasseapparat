export const fetchProducts = async (apiHost) => {
  try {
    const response = await fetch(apiHost + '/api/v1/products?_end=1000&_sort=pos&_order=asc')
    if (!response.ok) {
      throw new Error('Failed to fetch products')
    }
    const data = await response.json()
    return data
  } catch (error) {
    console.error(error)
    return null
  }
}

export const storePurchase = async (apiHost, cart, totalPrice) => {
  try {
    const response = await fetch(apiHost + '/api/v1/purchases', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ cart, totalPrice: cart.reduce((total, item) => total + item.totalPrice, 0) }) // : )
    })
    if (!response.ok) {
      throw new Error('Failed to store purchase')
    }
    const data = await response.json()
    return data.purchase
  } catch (error) {
    console.error(error)
    return null
  }
}

export const fetchPurchases = async (apiHost) => {
  try {
    const response = await fetch(apiHost + '/api/v1/purchases')
    if (!response.ok) {
      throw new Error('Failed to fetch purchases')
    }
    const data = await response.json()
    return data
  } catch (error) {
    console.error(error)
    return null
  }
}
