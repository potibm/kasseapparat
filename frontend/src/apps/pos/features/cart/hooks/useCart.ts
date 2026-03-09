// src/apps/pos/features/cart/hooks/useCart.ts
import { useState, useCallback } from "react";
import { Cart } from "../services/Cart";
import { storePurchase } from "../../../utils/api";
import { Product } from "../../product-list/types/product.types";
import { Purchase } from "../../../utils/api.schemas";
import { PaymentMethodData } from "../types/cart.types";
import { Guest } from "../../guestlist/types/guest.types";

interface EnrichedPurchase extends Purchase {
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

  const add = useCallback(
    (product: Product, count: number, listItem: Guest | null) => {
      setCart((prevCart) => prevCart.add(product, count, listItem));
    },
    [],
  );

  const remove = useCallback((product: Product) => {
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

      if (createdPurchase.status === "confirmed") {
        clear();
        setCheckoutProcessing(null);
        return createdPurchase;
      }

      throw new Error("Unknown purchase status: " + createdPurchase.status);
    } catch (error: unknown) {
      setCheckoutProcessing(null);
      throw error;
    }
  };

  return {
    cart,
    items: cart.items,
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
