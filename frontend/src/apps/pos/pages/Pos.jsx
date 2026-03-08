import "../assets/styles/pos-style.css";

import React, { useState, useCallback } from "react";
import { Alert } from "flowbite-react";
import Cart from "../features/cart/components/Cart";
import ProductList from "../features/product-list/components/ProductList";
import PurchaseHistory from "../features/purchase-history/components/PurchaseHistory";
import ErrorModal from "../components/ErrorModal";
import Menu from "../features/menu/compontents/Menu";
import PollingModal from "../features/purchase/components/PollingModal";
import { useAuth } from "../features/auth/providers/auth-provider";
import { useConfig } from "../../../core/config/providers/config-provider";
import Version from "../components/Version";
import PosLayout from "../layouts/PosLayout";
import { useProducts } from "../features/product-list/hooks/useProducts";
import { useCart } from "../features/cart/hooks/useCart";
import { usePurchaseHistory } from "../features/purchase-history/hooks/usePurchaseHistory";

const Kasseapparat = () => {
  const { apiHost, environmentMessage } = useConfig();
  const { username, getToken, id: userId } = useAuth();
  const [errorMessage, setErrorMessage] = useState("");

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
    addInterest,
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

  const { history, refreshHistory, refundPurchase } = usePurchaseHistory(
    apiHost,
    getToken,
    userId,
    showError,
  );

  const handleCheckout = async (paymentMethodCode, paymentMethodData) => {
    try {
      await checkout(paymentMethodCode, paymentMethodData);
    } catch (error) {
      showError(error.message);
    } finally {
      await Promise.all([
        refreshHistory(), // Historie neu vom Server laden
        refreshProducts(), // Lagerbestände/Produkte aktualisieren
      ]);
    }
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
            checkoutCart={handleCheckout}
          />
          <PurchaseHistory
            history={history}
            removeFromPurchaseHistory={(p) => refundPurchase(p.id)}
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
        addProductInterest={(p) => addInterest(p.id)}
      />
    </PosLayout>
  );
};

export default Kasseapparat;
