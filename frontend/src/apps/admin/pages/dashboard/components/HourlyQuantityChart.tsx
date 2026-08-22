import * as React from "react";
import { useEffect, useState } from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { useDataProvider, RaRecord } from "react-admin";
import { Typography } from "@mui/material";
import { createLogger } from "@core/logger/logger";
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

const log = createLogger("Admin");

interface HourlyQuantityStat extends RaRecord {
  id: string;
  timeBucket: string;
  productId: number;
  productName: string;
  quantity: number;
}

interface ChartDataPoint {
  hour: string;
  [key: string]: string | number;
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
  const [stats, setStats] = useState<HourlyQuantityStat[] | null>(null);
  const dataProvider = useDataProvider();

  useEffect(() => {
    dataProvider
      .getList<HourlyQuantityStat>("hourlyQuantityStats", {
        pagination: { page: 1, perPage: 1000 },
        sort: { field: "hour", order: "ASC" },
        filter: {},
      })
      .then(({ data }) => {
        setStats(data);
      })
      .catch((error) => {
        log.error("Hourly quantity stats fetch failed", error);
        setStats([]);
      });
  }, [dataProvider]);

  if (stats === null) return <Typography sx={{ p: 2 }}>Loading...</Typography>;
  if (stats.length === 0)
    return <Typography sx={{ p: 2 }}>No data available.</Typography>;

  const parseTimeBucket = (bucket: string): Date => {
    return new Date(bucket.replace(" ", "T") + ":00Z");
  };

  const timeBuckets = stats.map((s) => parseTimeBucket(s.timeBucket));
  const minTime = new Date(Math.min(...timeBuckets.map((t) => t.getTime())));
  const maxTime = new Date(Math.max(...timeBuckets.map((t) => t.getTime())));

  const spanHours = (maxTime.getTime() - minTime.getTime()) / (1000 * 60 * 60);
  const granularityMinutes = spanHours < 12 ? 15 : spanHours < 24 ? 30 : 60;

  const allBuckets: string[] = [];
  const current = new Date(minTime);
  current.setUTCSeconds(0, 0);
  const currentMinutes = current.getUTCMinutes();
  current.setUTCMinutes(
    Math.floor(currentMinutes / granularityMinutes) * granularityMinutes,
  );

  while (current <= maxTime) {
    const bucketStr =
      current.getUTCFullYear() +
      "-" +
      String(current.getUTCMonth() + 1).padStart(2, "0") +
      "-" +
      String(current.getUTCDate()).padStart(2, "0") +
      " " +
      String(current.getUTCHours()).padStart(2, "0") +
      ":" +
      String(current.getUTCMinutes()).padStart(2, "0");
    allBuckets.push(bucketStr);
    current.setUTCMinutes(current.getUTCMinutes() + granularityMinutes);
  }

  const productSet = new Set<number>();
  const productNames: Record<number, string> = {};
  stats.forEach((stat) => {
    productSet.add(stat.productId);
    productNames[stat.productId] = stat.productName;
  });

  const dataByBucket: Record<string, Record<number, number>> = {};
  allBuckets.forEach((bucket) => {
    dataByBucket[bucket] = {};
    productSet.forEach((pid) => {
      dataByBucket[bucket][pid] = 0;
    });
  });

  stats.forEach((stat) => {
    const statTime = parseTimeBucket(stat.timeBucket);
    const bucketTime = new Date(statTime);
    const minutes = bucketTime.getUTCMinutes();
    bucketTime.setUTCMinutes(
      Math.floor(minutes / granularityMinutes) * granularityMinutes,
    );
    const bucketStr =
      bucketTime.getUTCFullYear() +
      "-" +
      String(bucketTime.getUTCMonth() + 1).padStart(2, "0") +
      "-" +
      String(bucketTime.getUTCDate()).padStart(2, "0") +
      " " +
      String(bucketTime.getUTCHours()).padStart(2, "0") +
      ":" +
      String(bucketTime.getUTCMinutes()).padStart(2, "0");

    if (dataByBucket[bucketStr]) {
      dataByBucket[bucketStr][stat.productId] += stat.quantity;
    }
  });

  const chartData: ChartDataPoint[] = allBuckets.map((bucket) => {
    const point: ChartDataPoint = { hour: bucket };
    productSet.forEach((pid) => {
      point[pid] = dataByBucket[bucket][pid];
    });
    return point;
  });

  const formatXAxis = (value: string) => {
    if (!value) return "";
    const parts = value.split(" ");
    if (parts.length < 2) return value;

    const date = new Date(value.replace(" ", "T") + ":00Z");
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    const hours = String(date.getHours()).padStart(2, "0");
    const minutes = String(date.getMinutes()).padStart(2, "0");

    if (granularityMinutes >= 60) {
      return `${month}-${day} ${hours}:00`;
    }
    return `${month}-${day} ${hours}:${minutes}`;
  };

  return (
    <Card sx={{ mt: 2, boxShadow: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
          Quantity Sold by Product (Last 3 Days)
        </Typography>
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="hour" tickFormatter={formatXAxis} />
            <YAxis />
            <Tooltip />
            <Legend />
            {Array.from(productSet).map((pid, index) => (
              <Area
                key={pid}
                type="monotone"
                dataKey={pid}
                stackId="1"
                name={productNames[pid] || `Product ${pid}`}
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
