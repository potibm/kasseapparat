import React from "react";
import {
  Card,
  CardContent,
  List as MuiList,
  ListItem,
  ListItemText,
  Typography,
} from "@mui/material";
import { useGetList } from "react-admin";
import {
  productsByInterest,
  ProductWithInterest,
} from "./interestSummary.utils";

/**
 * The interest per product at a glance – what to restock first – next to the list
 * of single requests.
 */
export const ProductInterestSummary: React.FC = () => {
  const { data, isPending } = useGetList<ProductWithInterest & { id: number }>(
    "products",
    {
      pagination: { page: 1, perPage: 1000 },
      sort: { field: "name", order: "ASC" },
    },
  );

  const products = productsByInterest(data ?? []);

  return (
    <Card sx={{ ml: 2, mt: 8, minWidth: 240, alignSelf: "flex-start" }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Per product
        </Typography>
        {!isPending && products.length === 0 && (
          <Typography variant="body2" color="text.secondary">
            Nobody asked for a sold-out product yet.
          </Typography>
        )}
        <MuiList dense data-testid="interest-summary">
          {products.map((product) => (
            <ListItem
              key={product.id}
              disableGutters
              secondaryAction={
                <Typography sx={{ fontWeight: 700 }}>
                  {product.soldOutRequestCount}
                </Typography>
              }
            >
              <ListItemText primary={product.name} />
            </ListItem>
          ))}
        </MuiList>
      </CardContent>
    </Card>
  );
};

export default ProductInterestSummary;
