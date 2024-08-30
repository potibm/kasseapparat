import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Cart from "./components/Cart";
import ProductList from "./components/ProductList";
import PurchaseHistory from "./components/PurchaseHistory";
import ErrorModal from "./components/ErrorModal";
import {
  deletePurchaseById,
  fetchProducts,
  fetchPurchases,
  storePurchase,
  addProductInterest,
} from "./hooks/Api";
import {
  addToCart,
  removeFromCart,
  removeAllFromCart,
  checkoutCart,
  containsListItemID,
  getCartProductQuantity,
} from "./hooks/Cart";
import { Button, Spinner } from "flowbite-react";
import { HiCog, HiOutlineUserCircle } from "react-icons/hi";
import { useAuth } from "../Auth/provider/AuthProvider";
import { useConfig } from "../provider/ConfigProvider";

function Kasseapparat() {
  const [cart, setCart] = useState([]);
  const [products, setProducts] = useState([]);
  const [purchaseHistory, setPurchaseHistory] = useState(null);
  const [errorMessage, setErrorMessage] = useState("");
  const { username, token } = useAuth();
  const version = useConfig().version;
  const apiHost = useConfig().apiHost;

  useEffect(() => {
    const getProducts = async () => {
      return fetchProducts(apiHost, token)
        .then((products) => setProducts(products))
        .catch((error) =>
          showError(
            "There was an error fetching the products: " + error.message,
          ),
        );
    };
    const getHistory = async () => {
      const history = await fetchPurchases(apiHost, token);
      setPurchaseHistory(history);
      fetchPurchases(apiHost, token)
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
  }, [apiHost, token]); // Empty dependency array to run only once on mount

  const handleAddToCart = (product, count = 1, listItem = null) => {
    setCart(addToCart(cart, product, count, listItem));
  };

  const handleRemoveFromCart = (product) => {
    setCart(removeFromCart(cart, product));
  };

  const hasListItem = (listItemID) => {
    return containsListItemID(cart, listItemID);
  };

  const handleRemoveAllFromCart = () => {
    setCart(removeAllFromCart());
    fetchProducts(apiHost, token)
      .then((products) => setProducts(products))
      .catch((error) =>
        showError("There was an error fetching the products: " + error.message),
      );
  };

  const handleAddToPurchaseHistory = (purchase) => {
    if (purchaseHistory === null) {
      return;
    }
    setPurchaseHistory([purchase, ...purchaseHistory]);
  };

  const handleRemoveFromPurchaseHistory = async (purchase) => {
    return deletePurchaseById(apiHost, token, purchase.id)
      .then((data) => {
        fetchPurchases(apiHost, token)
          .then((history) => setPurchaseHistory(history))
          .catch((error) =>
            showError(
              "There was an error fetching the purchase history: " +
                error.message,
            ),
          );
        fetchProducts(apiHost, token)
          .then((products) => setProducts(products))
          .catch((error) =>
            showError(
              "There was an error fetching the products: " + error.message,
            ),
          );
      })
      .catch((error) => {
        showError("There was an error deleting the purchase: " + error.message);
      });
  };

  const getQuantityByProductInCart = (product) => {
    return getCartProductQuantity(cart, product);
  };

  const handleCheckoutCart = async () => {
    return storePurchase(apiHost, token, cart)
      .then((createdPurchase) => {
        setCart(checkoutCart());
        handleAddToPurchaseHistory(createdPurchase.purchase);
        fetchProducts(apiHost, token)
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

  const handleAddProductInterest = (product) => {
    console.log("Adding product interest for product: ", product.id);
    return addProductInterest(apiHost, token, product.id)
      .then(() => {
        product.soldOutRequestCount++;
      })
      .catch((error) => {
        showError(
          "There was an error adding the product interest: " + error.message,
        );
      });
  };

  const showError = (message) => {
    setErrorMessage(message);
  };

  const handleCloseError = () => {
    setErrorMessage("");
  };

  return (
    <div className="App p-2 dark:bg-black">
      <div className="w-full overflow-hidden">
        {products.length === 0 && (
          <div className="w-9/12 text-gray-500 text-left p-5">
            Loading products...
            <Spinner className="ml-2" />
          </div>
        )}
        {products.length > 0 && (
          <div className="w-9/12">
            <ProductList
              products={products}
              addToCart={handleAddToCart}
              hasListItem={hasListItem}
              quantityByProductInCart={getQuantityByProductInCart}
              addProductInterest={handleAddProductInterest}
            />
          </div>
        )}
        <div className="fixed inset-y-0 right-0 w-3/12 bg-slate-200 dark:bg-gray-900 p-2">
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
              <HiOutlineUserCircle className="mr-2 h-5 w-5 max-lg:hidden" />{" "}
              {username}
            </Button>
            <Button
              as={Link}
              to="/logout"
              className="hover:bg-cyan-800 hover:text-white"
            >
              Logout
            </Button>
            <Button
              as={Link}
              target="blank"
              to="/admin"
              className="hover:bg-cyan-800 hover:text-white"
            >
              <HiCog className="mr-2 h-5 w-5 max-xl:hidden" /> Admin
            </Button>
          </Button.Group>

          <p className="text-xs mt-10 dark:text-white">Version {version}</p>
        </div>
      </div>
      <ErrorModal message={errorMessage} onClose={handleCloseError} />
    </div>
  );
}

export default Kasseapparat;
