// src/apps/pos/features/product-list/hooks/useProducts.ts
import { useState, useEffect, useCallback } from "react";
import { createLogger } from "@core/logger/logger";
import { useToast } from "@pos/features/ui/toast/hooks/useToast";
import { usePosApi } from "@pos/api/usePosApi";
import { Product, ProductGroup, Section } from "@pos/api/schemas";

const log = createLogger("Product");

/**
 * Loads data the POS can do without: on failure it logs and resolves to null, so
 * the products still show up (just without headings or size rows).
 */
const optional = <T>(load: Promise<T>, what: string): Promise<T | null> =>
  load.catch((error: unknown) => {
    log.error(
      "Error fetching " + what,
      error instanceof Error ? { message: error.message } : { error },
    );
    return null;
  });

export const useProducts = () => {
  const [products, setProducts] = useState<Product[] | null>(null);
  const [sections, setSections] = useState<Section[]>([]);
  const [productGroups, setProductGroups] = useState<ProductGroup[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const { showToast } = useToast();
  const {
    fetchProducts,
    fetchSections,
    fetchProductGroups,
    addProductInterest,
  } = usePosApi();

  const loadProducts = useCallback(async () => {
    setLoading(true);
    try {
      const [fetchedProducts, fetchedSections, fetchedGroups] =
        await Promise.all([
          fetchProducts(),
          optional(fetchSections(), "sections"),
          optional(fetchProductGroups(), "product groups"),
        ]);

      if (fetchedSections !== null) {
        setSections(fetchedSections);
      }

      if (fetchedGroups !== null) {
        setProductGroups(fetchedGroups);
      }

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
  }, [fetchProducts, fetchSections, fetchProductGroups, showToast]);

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
    productGroups,
    loading,
    refreshProducts: loadProducts,
    addInterest,
  };
};

export default useProducts;
