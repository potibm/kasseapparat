import { useState, useCallback, useRef, useEffect } from "react";
import { Cart } from "../services/Cart";
import { storePurchase } from "../../../utils/api";
import { PaymentMethodData } from "../types/cart.types";
import {
  Purchase as PurchaseType,
  Product as ProductType,
  Guest as GuestType,
} from "../../../utils/api.schemas";

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
      setCart((prevCart) => prevCart.add(product, count, listItem));
    },
    [],
  );

  const remove = useCallback((product: ProductType) => {
    setCart((prevCart) => prevCart.remove(product.id));
  }, []);

  const clear = useCallback(() => {
    setCart(new Cart());
  }, []);

  const checkout = async (
    paymentMethodCode: string,
    paymentMethodData: PaymentMethodData,
  ) => {
    setCheckoutProcessing(paymentMethodCode);

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
                clear();
                resolve(createdPurchase);
              } else {
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
        }
        return createdPurchase;
      }

      throw new Error("Unknown purchase status: " + createdPurchase.status);
    } catch (error: unknown) {
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
