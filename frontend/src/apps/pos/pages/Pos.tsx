import "../assets/styles/pos-style.css";

import React, { useState, useCallback } from "react";
import { Alert } from "flowbite-react";
// components & layouts
import Cart from "../features/cart/components/Cart";
import ProductList from "../features/product-list/components/ProductList";
import PurchaseHistory from "../features/purchase-history/components/PurchaseHistory";
import ErrorModal from "../components/ErrorModal";
import Menu from "../features/menu/components/Menu";
import PollingModal from "@pos/features/payment/components/PollingModal";
import Version from "../components/Version";
import PosLayout from "../layouts/PosLayout";
// hooks
import { useAuth } from "../features/auth/providers/AuthProvider";
import { useConfig } from "../../../core/config/providers/ConfigProvider";
import { useProducts } from "../features/product-list/hooks/useProducts";
import { useCart } from "../features/cart/hooks/useCart";
import { usePurchaseHistory } from "../features/purchase-history/hooks/usePurchaseHistory";
// types
import { PaymentMethodData } from "../features/cart/types/cart.types";
import {
  Product as ProductType,
  Purchase as PurchaseType,
  Guest as GuestType,
} from "../utils/api.schemas";
import { createLogger } from "@core/logger/logger";

const logPurchase = createLogger("Purchase");

const Kasseapparat: React.FC = () => {
  const { apiHost, environmentMessage } = useConfig();
  const { username, getSafeToken, id: userId } = useAuth();

  const [errorMessage, setErrorMessage] = useState<string>("");

  const showError = useCallback((message: string) => {
    setErrorMessage(message);
  }, []);

  const handleCloseError = () => {
    setErrorMessage("");
  };

  const {
    products,
    loading: _productsLoading,
    refreshProducts,
    addInterest,
  } = useProducts(apiHost, getSafeToken, showError);

  const {
    cart,
    add,
    remove,
    clear,
    checkout,
    checkoutProcessing,
    isPolling,
    pendingPurchase,
    finalizeCheckout,
  } = useCart(apiHost, getSafeToken);

  const {
    history,
    refreshHistory,
    refundPurchase,
    loading: historyLoading,
  } = usePurchaseHistory(apiHost, getSafeToken, userId, showError);

  const handlePurchaseSuccess = useCallback(async () => {
    await Promise.all([refreshHistory(), refreshProducts()]);
  }, [refreshHistory, refreshProducts]);

  const handleCheckout = useCallback(
    async (paymentMethodCode: string, paymentMethodData: PaymentMethodData) => {
      try {
        const purchase = await checkout(paymentMethodCode, paymentMethodData);

        // refresh directly when we know the purchase is successful, otherwise we wait for the polling to confirm it
        if (purchase.status === "confirmed") {
          try {
            await handlePurchaseSuccess();
          } catch (error: unknown) {
            logPurchase.error(
              "Refresh failed after immediate confirmation:",
              error,
            );
            showError(
              "Purchase was successful, but refreshing data failed. Please refresh the page.",
            );
          }
        }
        // on pending we wait for the PollingModal to confirm the purchase before refreshing
      } catch (error: unknown) {
        const errorMessage =
          error instanceof Error
            ? error.message
            : "An unknown error has occurred";

        showError(errorMessage);
      }
    },
    [checkout, handlePurchaseSuccess, showError],
  );

  const handleRefund = useCallback(
    async (purchaseId: string) => {
      try {
        await refundPurchase(purchaseId);
        await refreshProducts();
      } catch (error: unknown) {
        const errorMessage =
          error instanceof Error
            ? error.message
            : "An unknown error has occurred";

        showError(errorMessage);
      }
    },
    [refundPurchase, refreshProducts, showError],
  );

  const handlePurchaseModalComplete = useCallback(
    (success: boolean) => {
      finalizeCheckout(success);

      if (success) {
        handlePurchaseSuccess().catch((error: unknown) => {
          logPurchase.error("Refresh failed after polling:", error);
          showError(
            "Purchase was successful, but refreshing data failed. Please refresh the page.",
          );
        });
      }
    },
    [finalizeCheckout, handlePurchaseSuccess, showError],
  );

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
            loading={historyLoading}
            removeFromPurchaseHistory={(p: PurchaseType) => handleRefund(p.id)}
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
              onComplete={handlePurchaseModalComplete}
            />
          )}
        </>
      }
    >
      <ProductList
        products={products}
        addToCart={add}
        hasListItem={(g: GuestType) => cart.hasListItem(g.id)}
        quantityByProductInCart={(p: ProductType) => cart.getQuantity(p.id)}
        addProductInterest={(p: ProductType) => addInterest(p.id)}
      />
    </PosLayout>
  );
};

export default Kasseapparat;
