import * as React from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { RaRecord } from "react-admin";
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
import { useStatsData } from "../hooks/useStatsData";

interface PaymentMethodStat extends RaRecord {
  id: string;
  paymentMethod: string;
  name: string;
  purchaseCount: number;
  totalNetPrice: string | number;
  totalGrossPrice: string | number;
}

const PaymentMethodStatsCard: React.FC = () => {
  const { data: stats } = useStatsData<PaymentMethodStat>("paymentMethodStats");
  const { currency } = useConfig();

  if (stats === null) {
    return <Typography sx={{ p: 2 }}>Loading...</Typography>;
  }

  if (stats.length === 0) {
    return <Typography sx={{ p: 2 }}>No purchases yet.</Typography>;
  }

  const totalNet = stats.reduce(
    (acc, stat) => acc.add(new Decimal(stat.totalNetPrice)),
    new Decimal(0),
  );

  const totalGross = stats.reduce(
    (acc, stat) => acc.add(new Decimal(stat.totalGrossPrice)),
    new Decimal(0),
  );

  const totalPurchases = stats.reduce(
    (acc, stat) => acc + stat.purchaseCount,
    0,
  );

  return (
    <Card sx={{ mt: 2, boxShadow: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
          Payment Method Stats
        </Typography>

        <TableContainer
          component={Paper}
          elevation={0}
          sx={{ border: "1px solid", borderColor: "divider" }}
        >
          <Table size="small" aria-label="payment method stats table">
            <TableHead sx={{ backgroundColor: "action.hover" }}>
              <TableRow>
                <TableCell sx={{ fontWeight: "bold" }}>
                  Payment Method
                </TableCell>
                <TableCell align="right" sx={{ fontWeight: "bold" }}>
                  Purchases
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
                  <TableCell align="right">{stat.purchaseCount}</TableCell>
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

              <TableRow sx={{ backgroundColor: "action.selected" }}>
                <TableCell sx={{ fontWeight: "bold" }}>Total</TableCell>
                <TableCell align="right" sx={{ fontWeight: "bold" }}>
                  {totalPurchases}
                </TableCell>
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

export default PaymentMethodStatsCard;
