// src/apps/pos/features/product-list/hooks/useProducts.ts
import { useState, useEffect, useCallback } from "react";
import { fetchProducts, addProductInterest } from "../../../utils/api";
import { Product } from "../types/product.types";

export const useProducts = (
  apiHost: string,
  getToken: () => Promise<string>,
  onError: (msg: string) => void,
) => {
  const [products, setProducts] = useState<Product[] | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const loadProducts = useCallback(async () => {
    setLoading(true);
    try {
      const token = await getToken();
      const products = await fetchProducts(apiHost, token);

      setProducts(products);
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error
          ? "There was an error fetching the products: " + error.message
          : "An unknown error has occured";

      onError(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [apiHost, getToken, onError]);

  const addInterest = async (productId: number) => {
    try {
      const token = await getToken();
      await addProductInterest(apiHost, token, productId);
      await loadProducts();
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error
          ? "Error on saving the interest: " + error.message
          : "An unknown error has occured";

      onError(errorMessage);
    }
  };

  useEffect(() => {
    loadProducts();
  }, [loadProducts]);

  return { products, loading, refreshProducts: loadProducts, addInterest };
};

export default useProducts;
