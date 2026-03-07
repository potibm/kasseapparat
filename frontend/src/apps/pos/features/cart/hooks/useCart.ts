// src/apps/pos/features/cart/hooks/useCart.ts
import { useState, useCallback } from "react";
import { Cart } from "../services/Cart";
import { storePurchase } from "../../../utils/api";
import { Product } from "../../product-list/types/product.types";

export const useCart = (apiHost: string, getToken) => {
  const [cart, setCart] = useState(new Cart());
  const [isPolling, setIsPolling] = useState(false);
  const [pendingPurchase, setPendingPurchase] = useState(null);

  const add = useCallback((product: Product, count: number, listItem) => {
    setCart((prevCart) => prevCart.add(product, count, listItem));
  }, []);

  const remove = useCallback((product: Product) => {
    setCart((prevCart) => prevCart.remove(product.id));
  }, []);

  const clear = useCallback(() => {
    setCart(new Cart());
  }, []);

  const checkout = async (
    paymentMethodCode: string,
    paymentMethodData: any,
  ) => {
    const token = await getToken();
    const payload = cart.toApiPayload(paymentMethodCode, paymentMethodData);
    const createdPurchase = await storePurchase(apiHost, token, payload);

    if (createdPurchase.status === "pending") {
      setPendingPurchase(createdPurchase);
      setIsPolling(true);

      return new Promise((resolve, reject) => {
        // Diese Resolver-Funktion geben wir später an das Modal weiter
        createdPurchase.onComplete = (success) => {
          setIsPolling(false);
          if (success) {
            clear();
            resolve(createdPurchase);
          } else {
            reject(new Error("Payment failed or cancelled"));
          }
        };
      });
    }

    if (createdPurchase.status === "confirmed") {
      clear();
      return createdPurchase;
    }

    throw new Error("Unknown purchase status: " + createdPurchase.status);
  };

  return {
    cart,
    items: cart.items,
    add,
    remove,
    clear,
    checkout,
    isPolling,
    pendingPurchase,
    setIsPolling,
  };
};

export default useCart;
