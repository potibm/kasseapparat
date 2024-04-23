import React, { useState, useEffect } from 'react'
import './App.css'
import Cart from './components/Cart'
import ProductList from './components/ProductList'
import { fetchProducts, storePurchase } from './hooks/Api'
import { addToCart, removeFromCart, removeAllFromCart, checkoutCart } from './hooks/Cart'

const Currency = new Intl.NumberFormat('de-DE', {
  style: 'currency',
  currency: 'EUR',
  minimumFractionDigits: 0,
  maximumFractionDigits: 0
})

const API_HOST = process.env.REACT_APP_API_HOST ?? 'http://localhost:3001'

function App () {
  const [cart, setCart] = useState([])
  const [products, setProducts] = useState([])

  useEffect(() => {
    const getProducts = async () => {
      const products = await fetchProducts(API_HOST)
      if (products) {
        setProducts(products)
      }
    }
    getProducts()
  }, []) // Empty dependency array to run only once on mount

  const handleAddToCart = (product) => {
    setCart(addToCart(cart, product))
  }

  const handleRemoveFromCart = (product) => {
    setCart(removeFromCart(cart, product))
  }

  const handleRemoveAllFromCart = () => {
    setCart(removeAllFromCart())
    fetchProducts(API_HOST)
  }

  const handleCheckoutCart = () => {
    if (storePurchase(API_HOST, cart)) {
      setCart(checkoutCart())
      fetchProducts(API_HOST)
    }
  }

  return (
    <div className="App p-2">
      <div className="flex">
        <ProductList
          products={products}
          addToCart={handleAddToCart}
          currency={Currency} />
        <Cart cart={cart}
          currency={Currency}
          removeFromCart={handleRemoveFromCart}
          removeAllFromCart={handleRemoveAllFromCart}
          checkoutCart={handleCheckoutCart} />
      </div>
    </div>
  )
}

export default App
