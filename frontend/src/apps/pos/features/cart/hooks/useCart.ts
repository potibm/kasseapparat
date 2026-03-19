import { useState, useCallback, useRef, useEffect } from "react";
import { Cart } from "../services/Cart";
import { storePurchase } from "../../../utils/api";
import { PaymentMethodData } from "../types/cart.types";
import {
  Purchase as PurchaseType,
  Product as ProductType,
  Guest as GuestType,
} from "../../../utils/api.schemas";
import { createLogger } from "@core/logger/logger";

const cartLog = createLogger("Cart");
const purchaseLog = createLogger("Purchase");

interface EnrichedPurchase extends PurchaseType {
  onComplete: (success: boolean) => void;
}

export const useCart = (apiHost: string, getToken: () => Promise<string>) => {
  const [cart, setCart] = useState(new Cart());
  const [isPolling, setIsPolling] = useState(false);
  const [pendingPurchase, setPendingPurchase] =
    useState<EnrichedPurchase | null>(null);
  const [checkoutProcessing, setCheckoutProcessing] = useState<string | null>(
    null,
  );
  const isMounted = useRef(true);
  useEffect(() => {
    isMounted.current = true;
    return () => {
      isMounted.current = false;
    };
  }, []);

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
        return new Promise((resolve, reject) => {
          const enrichedPurchase = {
            ...createdPurchase,
            onComplete: (success: boolean) => {
              if (!isMounted.current) {
                reject(new Error("Component unmounted during polling"));
                return;
              }

              setIsPolling(false);
              setPendingPurchase(null);
              setCheckoutProcessing(null);

              if (success) {
                purchaseLog.info("Purchase completed successfully");
                clear();
                resolve(createdPurchase);
              } else {
                purchaseLog.error("Purchase failed or cancelled");
                reject(new Error("Payment failed or cancelled"));
              }
            },
          };

          setPendingPurchase(enrichedPurchase);
          setIsPolling(true);
        });
      }

      if (isMounted.current) {
        if (createdPurchase.status === "confirmed") {
          clear();
          setCheckoutProcessing(null);
          purchaseLog.info("Purchase confirmed immediately", {
            purchaseId: createdPurchase.id,
          });
        }
        return createdPurchase;
      }

      purchaseLog.error("Unknown purchase status", createdPurchase.status);
      throw new Error("Unknown purchase status: " + createdPurchase.status);
    } catch (error: unknown) {
      purchaseLog.error(
        "Error during checkout",
        error instanceof Error ? { message: error.message } : { error },
      );
      if (isMounted.current) {
        setCheckoutProcessing(null);
      }
      throw error;
    }
  };

  return {
    cart,
    add,
    remove,
    clear,
    checkout,
    checkoutProcessing,
    isPolling,
    pendingPurchase,
    setIsPolling,
  };
};

export default useCart;
