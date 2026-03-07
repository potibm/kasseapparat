import "../assets/styles/pos-style.css";

import React, { useState, useEffect, useCallback } from "react";
import { Alert, Spinner } from "flowbite-react";
import Cart from "../features/cart/components/Cart";
import ProductList from "../features/product-list/components/ProductList";
import PurchaseHistory from "../features/purchase-history/components/PurchaseHistory";
import ErrorModal from "../components/ErrorModal";
import Menu from "../features/menu/compontents/Menu";
import PollingModal from "../features/purchase/components/PollingModal";
import {
  refundPurchaseById,
  fetchPurchases,
  addProductInterest,
} from "../utils/api";
import { useAuth } from "../features/auth/providers/auth-provider";
import { useConfig } from "../../../core/config/providers/config-provider";
import Version from "../components/Version";
import PosLayout from "../layouts/PosLayout";
import { useProducts } from "../features/product-list/hooks/useProducts";
import useCart from "../features/cart/hooks/useCart";

const Kasseapparat = () => {
  const { apiHost, environmentMessage } = useConfig();
  const { username, getToken, id: userId } = useAuth();
  const [errorMessage, setErrorMessage] = useState("");

  const [purchaseHistory, setPurchaseHistory] = useState(null);

  const showError = useCallback((message) => {
    setErrorMessage(message);
  }, []);

  const handleCloseError = () => {
    setErrorMessage("");
  };

  const {
    products,
    loading: productsLoading,
    refreshProducts,
  } = useProducts(apiHost, getToken, showError);
  const {
    cart,
    add,
    remove,
    clear,
    checkout,
    checkoutProcessing,
    isPolling,
    pendingPurchase,
    setIsPolling,
  } = useCart(apiHost, getToken);

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
            checkoutProcessing={checkoutProcessing}
            removeFromCart={remove}
            removeAllFromCart={clear}
            checkoutCart={checkout}
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
          {isPolling && pendingPurchase && (
            <PollingModal
              purchase={pendingPurchase}
              onComplete={pendingPurchase.onComplete}
              onConfirmed={() => setIsPolling(false)}
            />
          )}
        </>
      }
    >
      <ProductList
        products={products}
        addToCart={add}
        hasListItem={(id) => cart.hasListItem(id)}
        quantityByProductInCart={(p) => cart.getQuantity(p.id)}
        addProductInterest={handleAddProductInterest}
      />
    </PosLayout>
  );
};

export default Kasseapparat;
