import React, { useState, useEffect } from 'react'
import './App.css'
import Cart from './components/Cart'
import ProductList from './components/ProductList'
import PurchaseHistory from './components/PurchaseHistory'
import ErrorModal from './components/ErrorModal'
import { fetchProducts, fetchPurchases, storePurchase } from './hooks/Api'
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
  const [purchaseHistory, setPurchaseHistory] = useState([])
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    const getProducts = async () => {
      const products = await fetchProducts(API_HOST)
      if (products) {
        setProducts(products)
      }
    }
    const getHistory = async () => {
      const history = await fetchPurchases(API_HOST)
      if (history && history.data) {
        setPurchaseHistory(history.data)
      }
    }
    getProducts()
    getHistory()
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

  const handleAddToPurchaseHistory = (purchase) => {
    setPurchaseHistory([purchase, ...purchaseHistory])
  }

  const handleCheckoutCart = async () => {
    try {
      const createdPurchase = await storePurchase(API_HOST, cart)
      if (createdPurchase) {
        setCart(checkoutCart())
        handleAddToPurchaseHistory(createdPurchase)
        fetchProducts(API_HOST)
      }
    } catch (error) {
      console.error('Failed to checkout cart:', error)
      showError(error.message, error)
    }
  }

  const showError = (message) => {
    setErrorMessage(message);
  };

  const handleCloseError = () => {
      setErrorMessage('');
  };

  return (
    <div className="App p-2">
      <div className="w-full overflow-hidden">
        <div className='border w-9/12'>
          <ProductList
            products={products}
            addToCart={handleAddToCart}
            currency={Currency} />
        </div>
        <div className='fixed inset-y-0 right-0 w-3/12 border bg-slate-200 p-2'>
          <Cart cart={cart}
            currency={Currency}
            removeFromCart={handleRemoveFromCart}
            removeAllFromCart={handleRemoveAllFromCart}
            checkoutCart={handleCheckoutCart} />

          <PurchaseHistory 
            currency={Currency}
            history={purchaseHistory}
          />
        </div>
      </div>
      <ErrorModal message={errorMessage} onClose={handleCloseError} />
    </div>
  )
}

export default App