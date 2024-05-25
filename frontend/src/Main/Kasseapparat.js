import React, { useState, useEffect } from 'react'
import Cart from './components/Cart'
import ProductList from './components/ProductList'
import PurchaseHistory from './components/PurchaseHistory'
import ErrorModal from './components/ErrorModal'
import { deletePurchaseById, fetchProducts, fetchPurchases, storePurchase } from './hooks/Api'
import { addToCart, removeFromCart, removeAllFromCart, checkoutCart } from './hooks/Cart'
import { Link } from 'react-router-dom'
import { Button } from 'flowbite-react'
import { HiCog,HiOutlineUserCircle } from "react-icons/hi";
import { useAuth } from "../provider/authProvider";

// @TODO retrieve those settings from the backend
const Currency = new Intl.NumberFormat('de-DE', {
  style: 'currency',
  currency: 'EUR',
  minimumFractionDigits: 0,
  maximumFractionDigits: 0
})

const API_HOST = process.env.REACT_APP_API_HOST ?? 'http://localhost:3001'

function Kasseapparat () {
  const [cart, setCart] = useState([])
  const [products, setProducts] = useState([])
  const [purchaseHistory, setPurchaseHistory] = useState([])
  const [errorMessage, setErrorMessage] = useState('');
  const { username } = useAuth();

  useEffect(() => {
    const getProducts = async () => {
      fetchProducts(API_HOST)
      .then(products => setProducts(products))
      .catch(error => showError("There was an error fetching the products: " + error.message));

    }
    const getHistory = async () => {
      const history = await fetchPurchases(API_HOST)
      setPurchaseHistory(history)
      fetchPurchases(API_HOST)
        .then(history => setPurchaseHistory(history))
        .catch(error => showError("There was an error fetching the purchase history: " + error.message));
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
    .then(products => setProducts(products))
    .catch(error => showError("There was an error fetching the products: " + error.message));

  }

  const handleAddToPurchaseHistory = (purchase) => {
    setPurchaseHistory([purchase, ...purchaseHistory])
  }

  const handleRemoveFromPurchaseHistory = (purchase) => {
    deletePurchaseById(API_HOST, purchase.id)
    .then(data => {
      fetchPurchases(API_HOST)
      .then(history => setPurchaseHistory(history))
      .catch(error => showError("There was an error fetching the purchase history: " + error.message));
    })
    .catch(error => {
        showError("There was an error deleting the purchase: " + error.message);
    });
  }

  const handleCheckoutCart = async () => {
    storePurchase(API_HOST, cart)
    .then(createdPurchase => {
      setCart(checkoutCart())
        handleAddToPurchaseHistory(createdPurchase.purchase)
        fetchProducts(API_HOST)
          .then(products => setProducts(products))
          .catch(error => showError("There was an error fetching the products: " + error.message));
      }
    )
    .catch(error => { 
      showError("There was an error storing the purchase: " + error.message);
    });
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
            removeFromPurchaseHistory={handleRemoveFromPurchaseHistory}
          />

          <div className="mt-10">
            <Link to="/admin" target="_blank"><Button><HiCog  className="mr-2 h-5 w-5"/> Admin</Button></Link>
            <Link to="/logout"><Button><HiOutlineUserCircle  className="mr-2 h-5 w-5"/> Logout {username}</Button></Link>
          </div>
        </div>
      </div>
      <ErrorModal message={errorMessage} onClose={handleCloseError} />
    </div>
  )
}

export default Kasseapparat