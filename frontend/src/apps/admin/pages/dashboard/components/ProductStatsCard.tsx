import * as React from "react";
import { useEffect, useState } from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { useDataProvider, RaRecord } from "react-admin";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
} from "@mui/material";
import { useConfig } from "@core/config/hooks/useConfig";
import Decimal from "decimal.js";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Admin");

interface ProductStat extends RaRecord {
  name: string;
  soldItems: number;
  totalNetPrice: string | number;
  totalGrossPrice: string | number;
}

const ProductStatsCard: React.FC = () => {
  const [stats, setStats] = useState<ProductStat[] | null>(null);
  const dataProvider = useDataProvider();
  const { currency } = useConfig(); // Destructuring ist cleaner

  useEffect(() => {
    dataProvider
      .getList<ProductStat>("productStats", {
        pagination: { page: 1, perPage: 100 },
        sort: { field: "name", order: "ASC" },
        filter: {},
      })
      .then(({ data }) => {
        setStats(data);
      })
      .catch((error) => {
        log.error("Dashboard fetch failed", error);
        setStats([]);
      });
  }, [dataProvider]);

  if (stats === null) return <Typography sx={{ p: 2 }}>Loading...</Typography>;
  if (stats.length === 0)
    return <Typography sx={{ p: 2 }}>No products yet.</Typography>;

  const totalNet = stats.reduce(
    (acc, stat) => acc.add(new Decimal(stat.totalNetPrice)),
    new Decimal(0),
  );

  const totalGross = stats.reduce(
    (acc, stat) => acc.add(new Decimal(stat.totalGrossPrice)),
    new Decimal(0),
  );

  return (
    <Card sx={{ mt: 2, boxShadow: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
          Product Sales Stats
        </Typography>

        <TableContainer
          component={Paper}
          elevation={0}
          sx={{ border: "1px solid", borderColor: "divider" }}
        >
          <Table size="small" aria-label="product stats table">
            <TableHead sx={{ backgroundColor: "action.hover" }}>
              <TableRow>
                <TableCell sx={{ fontWeight: "bold" }}>Product</TableCell>
                <TableCell align="right" sx={{ fontWeight: "bold" }}>
                  Units sold
                </TableCell>
                <TableCell align="right" sx={{ fontWeight: "bold" }}>
                  Revenue Net
                </TableCell>
                <TableCell align="right" sx={{ fontWeight: "bold" }}>
                  Revenue Gross
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {stats.map((stat) => (
                <TableRow
                  key={stat.id}
                  sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                >
                  <TableCell component="th" scope="row">
                    {stat.name}
                  </TableCell>
                  <TableCell align="right">{stat.soldItems}</TableCell>
                  <TableCell align="right">
                    {currency.format(
                      new Decimal(stat.totalNetPrice).toNumber(),
                    )}
                  </TableCell>
                  <TableCell align="right">
                    {currency.format(
                      new Decimal(stat.totalGrossPrice).toNumber(),
                    )}
                  </TableCell>
                </TableRow>
              ))}

              {/* Summary Row */}
              <TableRow sx={{ backgroundColor: "action.selected" }}>
                <TableCell sx={{ fontWeight: "bold" }}>Total</TableCell>
                <TableCell align="right">-</TableCell>
                <TableCell align="right" sx={{ fontWeight: "bold" }}>
                  {currency.format(totalNet.toNumber())}
                </TableCell>
                <TableCell align="right" sx={{ fontWeight: "bold" }}>
                  {currency.format(totalGross.toNumber())}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  );
};

export default ProductStatsCard;
