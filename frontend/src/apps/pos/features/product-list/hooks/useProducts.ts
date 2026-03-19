// src/apps/pos/features/product-list/hooks/useProducts.ts
import { useState, useEffect, useCallback } from "react";
import { fetchProducts, addProductInterest } from "../../../utils/api";
import { Product as ProductType } from "../../../utils/api.schemas";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Product");

export const useProducts = (
  apiHost: string,
  getToken: () => Promise<string>,
  onError: (msg: string) => void,
) => {
  const [products, setProducts] = useState<ProductType[] | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const loadProducts = useCallback(async () => {
    setLoading(true);
    try {
      const token = await getToken();
      const fetchedProducts = await fetchProducts(apiHost, token);

      setProducts(fetchedProducts);
      log.debug("Products fetched successfully", {
        productCount: fetchedProducts.length,
      });
    } catch (error: unknown) {
      log.error(
        "Error fetching products",
        error instanceof Error ? { message: error.message } : { error },
      );
      const errorMessage =
        error instanceof Error
          ? "There was an error fetching the products: " + error.message
          : "An unknown error has occurred";

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
          : "An unknown error has occurred";

      onError(errorMessage);
    }
  };

  useEffect(() => {
    loadProducts();
  }, [loadProducts]);

  return { products, loading, refreshProducts: loadProducts, addInterest };
};

export default useProducts;
