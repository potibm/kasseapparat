import * as React from "react";
import { Box, Card, CardContent, Typography } from "@mui/material";
import { RaRecord } from "react-admin";
import Decimal from "decimal.js";
import { useConfig } from "@core/config/hooks/useConfig";
import { useStatsData } from "../hooks/useStatsData";

interface PaymentMethodStat extends RaRecord {
  purchaseCount: number;
  totalGrossPrice: string | number;
}

const Kpi: React.FC<{ label: string; value: string }> = ({ label, value }) => (
  <Card sx={{ flex: "1 1 12rem", boxShadow: 3 }}>
    <CardContent>
      <Typography variant="body2" color="text.secondary">
        {label}
      </Typography>
      <Typography variant="h4" component="p" sx={{ fontWeight: 700 }}>
        {value}
      </Typography>
    </CardContent>
  </Card>
);

/** The three numbers asked for most: revenue, purchases and average purchase. */
const KpiCards: React.FC = () => {
  const { data: stats } = useStatsData<PaymentMethodStat>("paymentMethodStats");
  const { currency } = useConfig();

  if (stats === null) return null;

  const revenue = stats.reduce(
    (acc, stat) => acc.add(new Decimal(stat.totalGrossPrice)),
    new Decimal(0),
  );
  const purchases = stats.reduce((acc, stat) => acc + stat.purchaseCount, 0);
  const average = purchases > 0 ? revenue.div(purchases) : new Decimal(0);

  return (
    <Box
      sx={{ mt: 2, display: "flex", gap: 2, flexWrap: "wrap" }}
      data-testid="dashboard-kpis"
    >
      <Kpi label="Revenue" value={currency.format(revenue.toNumber())} />
      <Kpi label="Purchases" value={String(purchases)} />
      <Kpi
        label="Average purchase"
        value={currency.format(average.toNumber())}
      />
    </Box>
  );
};

export default KpiCards;
