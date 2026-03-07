// src/apps/pos/features/product-list/hooks/useProducts.ts
import { useState, useEffect, useCallback } from "react";
import Decimal from "decimal.js";
import { fetchProducts } from "../../../utils/api";
import { Product } from "../types/product.types";

export const useProducts = (apiHost: string, getToken: () => Promise<string>, onError: (msg: string) => void) => {
  const [products, setProducts] = useState<Product[] | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const loadProducts = useCallback(async () => {
    setLoading(true);
    try {
      const token = await getToken();
      const rawProducts = await fetchProducts(apiHost, token);
      
      // Hier findet die "Decimal"-Magie statt
      const converted = rawProducts.map((p: any) => ({
        ...p,
        netPrice: new Decimal(p.netPrice),
        grossPrice: new Decimal(p.grossPrice),
        vatAmount: new Decimal(p.vatAmount),
      }));
      
      setProducts(converted);
    } catch (error: any) {
      onError("There was an error fetching the products: " + error.message);
    } finally {
      setLoading(false);
    }
  }, [apiHost, getToken, onError]);

  useEffect(() => {
    loadProducts();
  }, [loadProducts]);

  return { products, loading, refreshProducts: loadProducts };
};

export default useProducts;