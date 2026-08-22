import * as React from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { Typography } from "@mui/material";
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
import { useHourlyChartData } from "../hooks/useHourlyChartData";
import { formatXAxisLabel } from "../utils/formatTimeBucket";

interface HourlyQuantityStat {
  id: string;
  timeBucket: string;
  productId: number;
  productName: string;
  quantity: number;
}

const PRODUCT_COLORS = [
  "#1f77b4",
  "#ff7f0e",
  "#2ca02c",
  "#d62728",
  "#9467bd",
  "#8c564b",
  "#e377c2",
  "#7f7f7f",
  "#bcbd22",
  "#17becf",
];

const HourlyQuantityChart: React.FC = () => {
  const { data, chartData, granularityMinutes, seriesKeys, seriesNames } =
    useHourlyChartData<HourlyQuantityStat>({
      resource: "hourlyQuantityStats",
      dataKey: "productId",
      aggregateFn: (acc, item) => acc + item.quantity,
      initialValue: 0,
    });

  // Update seriesNames with actual names from data
  if (data) {
    data.forEach((item) => {
      seriesNames[String(item.productId)] = item.productName;
    });
  }

  if (data === null) {
    return <Typography sx={{ p: 2 }}>Loading...</Typography>;
  }

  if (chartData.length === 0) {
    return <Typography sx={{ p: 2 }}>No data available.</Typography>;
  }

  return (
    <Card sx={{ mt: 2, boxShadow: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
          Quantity Sold by Product (Last 3 Days)
        </Typography>
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="hour"
              tickFormatter={(v) => formatXAxisLabel(v, granularityMinutes)}
            />
            <YAxis />
            <Tooltip />
            <Legend />
            {seriesKeys.map((pid, index) => (
              <Area
                key={pid}
                type="monotone"
                dataKey={pid}
                stackId="1"
                name={seriesNames[pid] || `Product ${pid}`}
                fill={PRODUCT_COLORS[index % PRODUCT_COLORS.length]}
                stroke={PRODUCT_COLORS[index % PRODUCT_COLORS.length]}
              />
            ))}
          </AreaChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
};

export default HourlyQuantityChart;
