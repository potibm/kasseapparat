import "../assets/styles/pos-style.css";

import React, { useState, useEffect } from "react";
import { Alert, Spinner } from "flowbite-react";
import Cart from "../features/cart/components/Cart";
import ProductList from "../features/product-list/components/ProductList";
import PurchaseHistory from "../features/purchase-history/components/PurchaseHistory";
import ErrorModal from "../components/ErrorModal";
import Menu from "../features/menu/compontents/Menu";
import PollingModal from "../features/purchase/components/PollingModal";
import {
  refundPurchaseById,
  fetchProducts,
  fetchPurchases,
  storePurchase,
  addProductInterest,
} from "../utils/api";
import {
  addToCart,
  removeFromCart,
  removeAllFromCart,
  checkoutCart,
  containsListItemID,
  getCartProductQuantity,
} from "../features/cart/components/services/cart.logic";
import { useAuth } from "../features/auth/providers/auth-provider";
import { useConfig } from "../../../core/config/providers/config-provider";
import Version from "../components/Version";
import PosLayout from "../layouts/PosLayout";
import { useProducts } from "../features/product-list/hooks/useProducts";

const Kasseapparat = () => {
  const [cart, setCart] = useState([]);
  //const [products, setProducts] = useState(null);
  const [purchaseHistory, setPurchaseHistory] = useState(null);
  const [errorMessage, setErrorMessage] = useState("");
  const { username, getToken, id: userId } = useAuth();
  const [pollingModalOpen, setPollingModalOpen] = useState(false);
  const [onPollingComplete, setOnPollingComplete] = useState(() => () => {});
  const [pendingPurchase, setPendingPurchase] = useState(null);
  const { apiHost, environmentMessage } = useConfig();

  const showError = (message) => {
    setErrorMessage(message);
  };

  const handleCloseError = () => {
    setErrorMessage("");
  };

  const { products, loading: productsLoading, refreshProducts } = useProducts(
    apiHost, 
    getToken, 
    showError
  );

  useEffect(() => {
    const getHistory = async () => {
      fetchPurchases(apiHost, await getToken(), userId)
        .then((history) => setPurchaseHistory(history))
        .catch((error) =>
          showError(
            "There was an error fetching the purchase history: " +
              error.message,
          ),
        );
    };
    getHistory();
  }, [apiHost, userId, getToken]);

  const handleAddToCart = (product, count = 1, listItem = null) => {
    setCart(addToCart(cart, product, count, listItem));
  };

  const handleRemoveFromCart = (product) => {
    setCart(removeFromCart(cart, product));
  };

  const hasListItem = (listItemID) => {
    return containsListItemID(cart, listItemID);
  };

  const handleRemoveAllFromCart = async () => {
    setCart(removeAllFromCart());
    const token = await getToken();
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
    console.log("Adding purchase to history: ", purchase);
    setPurchaseHistory([purchase, ...purchaseHistory]);
  };

  const handleRemoveFromPurchaseHistory = async (purchase) => {
    return refundPurchaseById(apiHost, await getToken(), purchase.id)
      .then(async () => {
        const token = await getToken();
        fetchPurchases(apiHost, token, userId)
          .then((history) => setPurchaseHistory(history))
          .catch((error) =>
            showError(
              "There was an error fetching the purchase history: " +
                error.message,
            ),
          );
        refreshProducts();
      })
      .catch((error) => {
        showError(
          "There was an error refunding the purchase: " + error.message,
        );
      });
  };

  const getQuantityByProductInCart = (product) => {
    return getCartProductQuantity(cart, product);
  };

  const handleCheckoutCart = async (paymentMethodCode, paymentMethodData) => {
    try {
      const createdPurchase = await storePurchase(
        apiHost,
        await getToken(),
        cart,
        paymentMethodCode,
        paymentMethodData,
      );

      // probably we need to wait for pending purchases to be processed
      console.log("Purchase created: ", createdPurchase);

      if (createdPurchase.status === "pending") {
        // open modal and wait

        let purchaseSucceeded;
        let pollingTimeoutTimer;

        try {
          purchaseSucceeded = await Promise.race([
            new Promise((resolve) => {
              setPollingModalOpen(true);
              setPendingPurchase(createdPurchase);
              setOnPollingComplete(() => resolve);
            }),
            new Promise((_, reject) => {
              const timeoutDuration = 3 * 60 * 1000; // 3 minutes timeout
              pollingTimeoutTimer = setTimeout(() => {
                reject(new Error("Polling timed out"));
              }, timeoutDuration);
            }),
          ]);

          if (purchaseSucceeded === false) {
            return;
          }
        } catch (error) {
          setPollingModalOpen(false);
          if (error.message === "Polling timed out") {
            showError("Payment processing timed out. Please try again.");
          } else {
            showError("An unexpected error occurred: " + error.message);
          }
        } finally {
          if (pollingTimeoutTimer) {
            clearTimeout(pollingTimeoutTimer);
          }
        }
      } else if (createdPurchase.status !== "confirmed") {
        throw new Error(
          "Purchase status is not confirmed: " + createdPurchase.status,
        );
      }
      console.log("Checkout successful, purchase: ", createdPurchase);

      setCart(checkoutCart());
      handleAddToPurchaseHistory(createdPurchase);
      refreshProducts();
    } catch (error) {
      showError("There was an error during checkout: " + error.message);
    }
  };

  const handleAddProductInterest = async (product) => {
    console.log("Adding product interest for product: ", product.id);
    return addProductInterest(apiHost, await getToken(), product.id)
      .then(() => {
        product.soldOutRequestCount++;
      })
      .catch((error) => {
        showError(
          "There was an error adding the product interest: " + error.message,
        );
      });
  };


  return (
    <PosLayout
      topAlert={
        environmentMessage && <Alert color="info">{environmentMessage}</Alert>
      }
      sidebar={
        <>
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
          <Menu username={username} />
          <p className="text-xs mt-10 dark:text-white">
            <Version />
          </p>
        </>
      }
      overlays={
        <>
          <ErrorModal message={errorMessage} onClose={handleCloseError} />
          {pollingModalOpen && (
            <PollingModal
              purchase={pendingPurchase}
              onComplete={onPollingComplete}
            />
          )}
        </>
      }
    >
      <ProductList
        products={products}
        addToCart={handleAddToCart}
        hasListItem={hasListItem}
        quantityByProductInCart={getQuantityByProductInCart}
        addProductInterest={handleAddProductInterest}
      />
    </PosLayout>
  );
};

export default Kasseapparat;
