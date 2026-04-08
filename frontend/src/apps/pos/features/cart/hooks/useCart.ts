import { useState, useCallback } from "react";
import { Cart } from "../services/Cart";
import { storePurchase } from "../../../utils/api";
import { PaymentMethodData } from "../types/cart.types";
import {
  Purchase as PurchaseType,
  Product as ProductType,
  Guest as GuestType,
} from "../../../utils/api.schemas";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";
import { useConfig } from "@core/config/hooks/useConfig";

const cartLog = createLogger("Cart");
const purchaseLog = createLogger("Purchase");

export const useCart = (apiHost: string, getToken: () => Promise<string>) => {
  const [cart, setCart] = useState(new Cart());
  const [isPolling, setIsPolling] = useState(false);
  const [pendingPurchase, setPendingPurchase] = useState<PurchaseType | null>(
    null,
  );
  const [checkoutProcessing, setCheckoutProcessing] = useState<string | null>(
    null,
  );
  const { showToast } = useToast();
  const { currency } = useConfig();

  const add = useCallback(
    (product: ProductType, count: number, listItem: GuestType | null) => {
      cartLog.debug("Adding product to cart", { productId: product.id, count });
      setCart((prevCart) => prevCart.add(product, count, listItem));
    },
    [],
  );

  const remove = useCallback((product: ProductType) => {
    cartLog.debug("Removing product from cart", { productId: product.id });
    setCart((prevCart) => prevCart.remove(product.id));
  }, []);

  const clear = useCallback(() => {
    cartLog.debug("Clearing cart");
    setCart(new Cart());
  }, []);

  const finalizeCheckout = useCallback(
    (success: boolean) => {
      setIsPolling(false);
      setPendingPurchase(null);
      setCheckoutProcessing(null);

      if (success) {
        purchaseLog.info("Purchase completed successfully via polling");
        clear();
      } else {
        purchaseLog.error("Purchase failed or cancelled during polling");
      }
    },
    [clear],
  );

  const checkout = async (
    paymentMethodCode: string,
    paymentMethodData: PaymentMethodData,
  ) => {
    setCheckoutProcessing(paymentMethodCode);
    purchaseLog.debug("Initiating purchase", { paymentMethodCode });

    try {
      const token = await getToken();
      const payload = cart.toApiPayload(paymentMethodCode, paymentMethodData);
      const createdPurchase = await storePurchase(apiHost, token, payload);

      if (createdPurchase.status === "pending") {
        setPendingPurchase(createdPurchase);
        setIsPolling(true);
        return createdPurchase;
      }

      if (createdPurchase.status === "confirmed") {
        clear();
        setCheckoutProcessing(null);
        purchaseLog.info("Purchase confirmed immediately", {
          purchaseId: createdPurchase.id,
        });

        showToast({
          severity: "success",
          message: `Purchase at ${currency.format(createdPurchase.totalGrossPrice.toNumber())} confirmed!`,
        });
        return createdPurchase;
      }

      purchaseLog.error("Unknown purchase status", createdPurchase.status);
      throw new Error("Unknown purchase status: " + createdPurchase.status);
    } catch (error: unknown) {
      purchaseLog.error(
        "Error during checkout",
        error instanceof Error ? { message: error.message } : { error },
      );
      showToast({
        severity: "error",
        message: "An error occurred while processing the purchase.",
        autoClose: false,
      });

      setCheckoutProcessing(null);
      throw error;
    }
  };

  return {
    cart,
    add,
    remove,
    clear,
    checkout,
    finalizeCheckout,
    checkoutProcessing,
    isPolling,
    pendingPurchase,
  };
};

export default useCart;
