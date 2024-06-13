import React, { useState, useEffect } from "react";
import Cart from "./components/Cart";
import ProductList from "./components/ProductList";
import PurchaseHistory from "./components/PurchaseHistory";
import ErrorModal from "./components/ErrorModal";
import {
  deletePurchaseById,
  fetchProducts,
  fetchPurchases,
  storePurchase,
} from "./hooks/Api";
import {
  addToCart,
  removeFromCart,
  removeAllFromCart,
  checkoutCart,
} from "./hooks/Cart";
import { Link } from "react-router-dom";
import { Button } from "flowbite-react";
import { HiCog, HiOutlineUserCircle } from "react-icons/hi";
import { useAuth } from "../provider/AuthProvider";

const API_HOST = process.env.REACT_APP_API_HOST ?? "http://localhost:3001";

function Kasseapparat() {
  const [cart, setCart] = useState([]);
  const [products, setProducts] = useState([]);
  const [purchaseHistory, setPurchaseHistory] = useState([]);
  const [errorMessage, setErrorMessage] = useState("");
  const { username, token } = useAuth();

  useEffect(() => {
    const getProducts = async () => {
      fetchProducts(API_HOST)
        .then((products) => setProducts(products))
        .catch((error) =>
          showError(
            "There was an error fetching the products: " + error.message,
          ),
        );
    };
    const getHistory = async () => {
      const history = await fetchPurchases(API_HOST);
      setPurchaseHistory(history);
      fetchPurchases(API_HOST)
        .then((history) => setPurchaseHistory(history))
        .catch((error) =>
          showError(
            "There was an error fetching the purchase history: " +
              error.message,
          ),
        );
    };
    getProducts();
    getHistory();
  }, []); // Empty dependency array to run only once on mount

  const handleAddToCart = (product) => {
    setCart(addToCart(cart, product));
  };

  const handleRemoveFromCart = (product) => {
    setCart(removeFromCart(cart, product));
  };

  const handleRemoveAllFromCart = () => {
    setCart(removeAllFromCart());
    fetchProducts(API_HOST)
      .then((products) => setProducts(products))
      .catch((error) =>
        showError("There was an error fetching the products: " + error.message),
      );
  };

  const handleAddToPurchaseHistory = (purchase) => {
    setPurchaseHistory([purchase, ...purchaseHistory]);
  };

  const handleRemoveFromPurchaseHistory = (purchase) => {
    deletePurchaseById(API_HOST, token, purchase.id)
      .then((data) => {
        fetchPurchases(API_HOST)
          .then((history) => setPurchaseHistory(history))
          .catch((error) =>
            showError(
              "There was an error fetching the purchase history: " +
                error.message,
            ),
          );
      })
      .catch((error) => {
        showError("There was an error deleting the purchase: " + error.message);
      });
  };

  const handleCheckoutCart = async () => {
    storePurchase(API_HOST, token, cart)
      .then((createdPurchase) => {
        setCart(checkoutCart());
        handleAddToPurchaseHistory(createdPurchase.purchase);
        fetchProducts(API_HOST)
          .then((products) => setProducts(products))
          .catch((error) =>
            showError(
              "There was an error fetching the products: " + error.message,
            ),
          );
      })
      .catch((error) => {
        showError("There was an error storing the purchase: " + error.message);
      });
  };

  const showError = (message) => {
    setErrorMessage(message);
  };

  const handleCloseError = () => {
    setErrorMessage("");
  };

  return (
    <div className="App p-2">
      <div className="w-full overflow-hidden">
        <div className="border w-9/12">
          <ProductList products={products} addToCart={handleAddToCart} />
        </div>
        <div className="fixed inset-y-0 right-0 w-3/12 border bg-slate-200 p-2">
          <Cart
            cart={cart}
            removeFromCart={handleRemoveFromCart}
            removeAllFromCart={handleRemoveAllFromCart}
            checkoutCart={handleCheckoutCart}
          />

          <PurchaseHistory
            history={purchaseHistory}
            removeFromPurchaseHistory={handleRemoveFromPurchaseHistory}
          />

          <Button.Group className="mt-10">
            <Button>
              <HiOutlineUserCircle className="mr-2 h-5 w-5" /> {username}
            </Button>
            <Button as={Link} to="/logout">
              Logout
            </Button>
            <Button as={Link} target="blank" to="/admin">
              <HiCog className="mr-2 h-5 w-5" /> Admin
            </Button>
          </Button.Group>
        </div>
      </div>
      <ErrorModal message={errorMessage} onClose={handleCloseError} />
    </div>
  );
}

export default Kasseapparat;
