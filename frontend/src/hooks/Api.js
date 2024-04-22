export const fetchProducts = async (apiHost) => {
  try {
    const response = await fetch(apiHost + '/api/v1/products')
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
