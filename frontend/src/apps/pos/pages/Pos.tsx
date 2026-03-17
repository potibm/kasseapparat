import "../assets/styles/pos-style.css";

import { useState, useCallback } from "react";
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

const Kasseapparat = () => {
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
    setIsPolling,
  } = useCart(apiHost, getSafeToken);

  const { history, refreshHistory, refundPurchase } = usePurchaseHistory(
    apiHost,
    getSafeToken,
    userId,
    showError,
  );

  const handleCheckout = async (
    paymentMethodCode: string,
    paymentMethodData: PaymentMethodData,
  ) => {
    try {
      await checkout(paymentMethodCode, paymentMethodData);
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : "An unknown error has occured";

      showError(errorMessage);
    } finally {
      await Promise.all([
        refreshHistory(), // Historie neu vom Server laden
        refreshProducts(), // Lagerbestände/Produkte aktualisieren
      ]);
    }
  };

  const handleRefund = async (purchaseId: string) => {
    try {
      await refundPurchase(purchaseId);
      await refreshProducts();
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : "An unknown error has occured";

      showError(errorMessage);
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
              onComplete={pendingPurchase.onComplete}
              onConfirmed={() => setIsPolling(false)}
              onClose={() => setIsPolling(false)}
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
