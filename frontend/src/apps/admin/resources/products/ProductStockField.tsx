import React from "react";
import { Chip, Stack, Typography } from "@mui/material";
import { useRecordContext } from "react-admin";
import { summarizeStock, StockTone } from "./stock.utils";
import { Product } from "./types";

const chipColor: Record<StockTone, "default" | "warning" | "error"> = {
  unlimited: "default",
  normal: "default",
  low: "warning",
  out: "error",
};

/**
 * Stock at a glance: units left of the total (orange when low, red when gone) and
 * the units sold, plus how many people asked for a sold-out product.
 */
export const ProductStockField: React.FC<{ label?: string }> = () => {
  const product = useRecordContext<Product>();
  if (!product) return null;

  const stock = summarizeStock(product);
  const interest = product.soldOutRequestCount ?? 0;

  return (
    <Stack direction="row" spacing={1} alignItems="center">
      {stock.total === null ? (
        <Typography variant="body2" color="text.secondary">
          {stock.sold} sold
        </Typography>
      ) : (
        <Chip
          size="small"
          color={chipColor[stock.tone]}
          label={`${stock.left} / ${stock.total} left`}
          data-testid="product-stock"
        />
      )}
      {interest > 0 && (
        <Typography variant="body2" color="text.secondary">
          {interest} interested
        </Typography>
      )}
    </Stack>
  );
};

export default ProductStockField;
