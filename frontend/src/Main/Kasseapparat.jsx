import React, { useState, useEffect } from "react";
import { Alert, Spinner } from "flowbite-react";
import Cart from "./components/Cart/Cart";
import ProductList from "./components/ProductList";
import PurchaseHistory from "./components/PurchaseHistory";
import ErrorModal from "./components/ErrorModal";
import MainMenu from "./components/MainMenu/MainMenu";
import PollingModal from "./components/Purchase/PollingModal";
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
import { useAuth } from "../Auth/provider/AuthProvider";
import { useConfig } from "../provider/ConfigProvider";
import Version from "../components/Version";
import Decimal from "decimal.js";

const Kasseapparat = () => {
  const [cart, setCart] = useState([]);
  const [products, setProducts] = useState(null);
  const [purchaseHistory, setPurchaseHistory] = useState(null);
  const [errorMessage, setErrorMessage] = useState("");
  const { username, token, id: userId } = useAuth();
  const [pollingModalOpen, setPollingModalOpen] = useState(false);
  const [pendingPurchase, setPendingPurchase] = useState(null);
  const apiHost = useConfig().apiHost;
  const envMessage = useConfig().environmentMessage;

  const convertProductsWithDecimals = (products) => {
    return products.map((product) => {
      return {
        ...product,
        netPrice: new Decimal(product.netPrice),
        grossPrice: new Decimal(product.grossPrice),
        vatAmount: new Decimal(product.vatAmount),
      };
    });
  };

  useEffect(() => {
    const getProducts = async () => {
      return fetchProducts(apiHost, token)
        .then((products) => setProducts(convertProductsWithDecimals(products)))
        .catch((error) =>
          showError(
            "There was an error fetching the products: " + error.message,
          ),
        );
    };
    const getHistory = async () => {
      fetchPurchases(apiHost, token, userId)
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
  }, [apiHost, token, userId]); // Empty dependency array to run only once on mount

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
      .then((products) => setProducts(convertProductsWithDecimals(products)))
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
      .then(() => {
        fetchPurchases(apiHost, token, userId)
          .then((history) => setPurchaseHistory(history))
          .catch((error) =>
            showError(
              "There was an error fetching the purchase history: " +
                error.message,
            ),
          );
        fetchProducts(apiHost, token)
          .then((products) =>
            setProducts(convertProductsWithDecimals(products)),
          )
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

  const handleCheckoutCart = async (paymentMethodCode, paymentMethodData) => {
    return storePurchase(
      apiHost,
      token,
      cart,
      paymentMethodCode,
      paymentMethodData,
    )
      .then((createdPurchase) => {
        // probably we need to wait for pending purchases to be processed
        console.log("Purchase created: ", createdPurchase);
        if (createdPurchase.status === "pending") {
          // open the polling modal
          setPendingPurchase(createdPurchase);
          setPollingModalOpen(true);
        } else if (createdPurchase.status !== "confirmed") {
          throw new Error(
            "Purchase status is not pending or confirmed: " +
              createdPurchase.status,
          );
        }

        setCart(checkoutCart());
        handleAddToPurchaseHistory(createdPurchase);
        fetchProducts(apiHost, token)
          .then((products) =>
            setProducts(convertProductsWithDecimals(products)),
          )
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
      {envMessage && (
        <Alert color="info" className="mb-5 rounded-none">
          {envMessage}
        </Alert>
      )}
      <div className="w-full overflow-hidden">
        {products === null && (
          <div className="w-9/12 text-gray-500 text-left p-5">
            Loading products...
            <Spinner className="ml-2" />
          </div>
        )}
        {products !== null && products.length === 0 && (
          <div className="w-9/12 text-gray-500 text-left p-5">
            No products, yet.
          </div>
        )}
        {products !== null && products.length > 0 && (
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

          <MainMenu username={username} />

          <p className="text-xs mt-10 dark:text-white">
            <Version />
          </p>
        </div>
      </div>
      <ErrorModal message={errorMessage} onClose={handleCloseError} />
      {pollingModalOpen && pendingPurchase && (
        <PollingModal
          purchase={pendingPurchase}
          show={pollingModalOpen}
          onClose={() => setPollingModalOpen(false)}
          onConfirmed={(updatedPurchase) => {
            setPollingModalOpen(false);
            handleAddToPurchaseHistory(updatedPurchase);
            setCart(checkoutCart());
          }}
        />
      )}
    </div>
  );
};

export default Kasseapparat;
