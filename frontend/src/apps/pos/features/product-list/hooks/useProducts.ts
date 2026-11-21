// src/apps/pos/features/product-list/hooks/useProducts.ts
import { useState, useEffect, useCallback } from "react";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";
import { usePosApi } from "@pos/api/usePosApi";
import { useConfig } from "@core/config/hooks/useConfig";
import { Section } from "@pos/api/schemas";
import {
  DisplayProduct,
  keepVanishedProducts,
} from "../utils/vanishedProducts";

const log = createLogger("Product");

export const useProducts = () => {
  const [products, setProducts] = useState<DisplayProduct[] | null>(null);
  const [sections, setSections] = useState<Section[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const { showToast } = useToast();
  const { fetchProducts, fetchSections, addProductInterest } = usePosApi();
  const { outOfStockBehavior } = useConfig();

  const loadProducts = useCallback(async () => {
    setLoading(true);
    try {
      // Without sections the POS still works, it just shows no headings.
      const [fetchedProducts, fetchedSections] = await Promise.all([
        fetchProducts(),
        fetchSections().catch((error: unknown) => {
          log.error(
            "Error fetching sections",
            error instanceof Error ? { message: error.message } : { error },
          );
          return null;
        }),
      ]);

      if (fetchedSections !== null) {
        setSections(fetchedSections);
      }

      setProducts((previous) =>
        keepVanishedProducts(previous, fetchedProducts, outOfStockBehavior),
      );
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
  }, [fetchProducts, fetchSections, showToast, outOfStockBehavior]);

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

  return {
    products,
    sections,
    loading,
    refreshProducts: loadProducts,
    addInterest,
  };
};

export default useProducts;
