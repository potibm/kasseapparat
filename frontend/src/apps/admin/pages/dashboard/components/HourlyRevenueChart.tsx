import * as React from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { Typography } from "@mui/material";
import { useConfig } from "@core/config/hooks/useConfig";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";
import Decimal from "decimal.js";
import { useHourlyChartData } from "../hooks/useHourlyChartData";
import { formatXAxisLabel } from "../utils/formatTimeBucket";
import { paymentMethodColor } from "@admin/ui/colors";

interface HourlyRevenueStat {
  id: string;
  timeBucket: string;
  paymentMethod: string;
  name: string;
  totalGrossPrice: string | number;
}

const getColor = paymentMethodColor;

const HourlyRevenueChart: React.FC = () => {
  const { currency } = useConfig();
  const { data, chartData, granularityMinutes, seriesKeys, seriesNames } =
    useHourlyChartData<HourlyRevenueStat>({
      resource: "hourlyRevenueStats",
      dataKey: "paymentMethod",
      aggregateFn: (acc, item) =>
        acc + new Decimal(item.totalGrossPrice).toNumber(),
      initialValue: 0,
    });

  // Update seriesNames with actual names from data
  if (data) {
    data.forEach((item) => {
      seriesNames[item.paymentMethod] = item.name;
    });
  }

  if (data === null) {
    return <Typography sx={{ p: 2 }}>Loading...</Typography>;
  }

  if (chartData.length === 0) {
    return <Typography sx={{ p: 2 }}>No data available.</Typography>;
  }

  const formatTooltip = (value: unknown) => {
    if (typeof value === "number") {
      return currency.format(value);
    }
    return String(value);
  };

  return (
    <Card sx={{ mt: 2, boxShadow: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
          Revenue by Payment Method (Last 3 Days)
        </Typography>
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="hour"
              tickFormatter={(v) => formatXAxisLabel(v, granularityMinutes)}
            />
            <YAxis tickFormatter={(v) => currency.format(v)} />
            <Tooltip formatter={formatTooltip} />
            <Legend />
            {seriesKeys.map((pm, index) => (
              <Area
                key={pm}
                type="monotone"
                dataKey={pm}
                stackId="1"
                name={seriesNames[pm] || pm}
                fill={getColor(pm, index)}
                stroke={getColor(pm, index)}
              />
            ))}
          </AreaChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
};

export default HourlyRevenueChart;
