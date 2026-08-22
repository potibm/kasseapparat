import * as React from "react";
import { useEffect, useState } from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { useDataProvider, RaRecord } from "react-admin";
import { Typography } from "@mui/material";
import { useConfig } from "@core/config/hooks/useConfig";
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
import Decimal from "decimal.js";

const log = createLogger("Admin");

interface HourlyRevenueStat extends RaRecord {
  id: string;
  timeBucket: string;
  paymentMethod: string;
  name: string;
  totalGrossPrice: string | number;
}

interface ChartDataPoint {
  hour: string;
  [key: string]: string | number;
}

const PAYMENT_METHOD_COLORS: Record<string, string> = {
  CASH: "#4caf50",
  CC: "#2196f3",
  SUMUP: "#ff9800",
  VOUCHER: "#9c27b0",
};

const getColor = (paymentMethod: string, index: number): string => {
  return (
    PAYMENT_METHOD_COLORS[paymentMethod] ||
    ["#e91e63", "#00bcd4", "#ff5722", "#607d8b", "#795548"][index % 5]
  );
};

const HourlyRevenueChart: React.FC = () => {
  const [stats, setStats] = useState<HourlyRevenueStat[] | null>(null);
  const dataProvider = useDataProvider();
  const { currency } = useConfig();

  useEffect(() => {
    dataProvider
      .getList<HourlyRevenueStat>("hourlyRevenueStats", {
        pagination: { page: 1, perPage: 1000 },
        sort: { field: "hour", order: "ASC" },
        filter: {},
      })
      .then(({ data }) => {
        setStats(data);
      })
      .catch((error) => {
        log.error("Hourly revenue stats fetch failed", error);
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

  const paymentMethodSet = new Set<string>();
  const paymentMethodNames: Record<string, string> = {};
  stats.forEach((stat) => {
    paymentMethodSet.add(stat.paymentMethod);
    paymentMethodNames[stat.paymentMethod] = stat.name;
  });

  const dataByBucket: Record<string, Record<string, Decimal>> = {};
  allBuckets.forEach((bucket) => {
    dataByBucket[bucket] = {};
    paymentMethodSet.forEach((pm) => {
      dataByBucket[bucket][pm] = new Decimal(0);
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
      dataByBucket[bucketStr][stat.paymentMethod] = dataByBucket[bucketStr][
        stat.paymentMethod
      ].add(new Decimal(stat.totalGrossPrice));
    }
  });

  const chartData: ChartDataPoint[] = allBuckets.map((bucket) => {
    const point: ChartDataPoint = { hour: bucket };
    paymentMethodSet.forEach((pm) => {
      point[pm] = dataByBucket[bucket][pm].toNumber();
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
            <XAxis dataKey="hour" tickFormatter={formatXAxis} />
            <YAxis tickFormatter={(v) => currency.format(v)} />
            <Tooltip formatter={formatTooltip} />
            <Legend />
            {Array.from(paymentMethodSet).map((pm, index) => (
              <Area
                key={pm}
                type="monotone"
                dataKey={pm}
                stackId="1"
                name={paymentMethodNames[pm] || pm}
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
