// src/apps/pos/features/product-list/hooks/useProducts.ts
import { useState, useEffect, useCallback } from "react";
import { Product as ProductType } from "../../../api/schemas";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";
import { usePosApi } from "@pos/api/usePosApi";

const log = createLogger("Product");

export const useProducts = () => {
  const [products, setProducts] = useState<ProductType[] | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const { showToast } = useToast();
  const { fetchProducts, addProductInterest } = usePosApi();

  const loadProducts = useCallback(async () => {
    setLoading(true);
    try {
      const fetchedProducts = await fetchProducts();

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
      showToast({ severity: "error", message: errorMessage, autoClose: false });
    } finally {
      setLoading(false);
    }
  }, [fetchProducts, showToast]);

  const addInterest = async (productId: number, productName: string) => {
    try {
      await addProductInterest(productId);
      showToast({
        severity: "success",
        message: `Interest in ${productName} registered successfully!`,
      });
      await loadProducts();
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error
          ? "Error on saving the interest: " + error.message
          : "An unknown error has occurred";

      showToast({ severity: "error", message: errorMessage, autoClose: false });
    }
  };

  useEffect(() => {
    loadProducts();
  }, [loadProducts]);

  return { products, loading, refreshProducts: loadProducts, addInterest };
};

export default useProducts;
