import { useEffect, useState } from "react";
import { useDataProvider, RaRecord } from "react-admin";
import { createLogger } from "@core/logger/logger";
import {
  parseTimeBucket,
  calculateGranularity,
  generateTimeBuckets,
  formatTimeBucketString,
  snapToGranularity,
} from "../utils/formatTimeBucket";

const log = createLogger("Admin");

interface ChartDataPoint {
  hour: string;
  [key: string]: string | number;
}

interface UseHourlyChartDataOptions<T> {
  resource: string;
  dataKey: keyof T & string;
  aggregateFn: (acc: number, item: T) => number;
  initialValue?: number;
}

interface UseHourlyChartDataResult<T> {
  data: T[] | null;
  chartData: ChartDataPoint[];
  granularityMinutes: number;
  seriesKeys: string[];
  seriesNames: Record<string, string>;
}

export const useHourlyChartData = <
  T extends RaRecord & { timeBucket: string },
>({
  resource,
  dataKey,
  aggregateFn,
  initialValue = 0,
}: UseHourlyChartDataOptions<T>): UseHourlyChartDataResult<T> => {
  const [data, setData] = useState<T[] | null>(null);
  const dataProvider = useDataProvider();

  useEffect(() => {
    dataProvider
      .getList<T>(resource, {
        pagination: { page: 1, perPage: 1000 },
        sort: { field: "timeBucket", order: "ASC" },
        filter: {},
      })
      .then(({ data }) => {
        setData(data);
      })
      .catch((error) => {
        log.error(`${resource} fetch failed`, error);
        setData([]);
      });
  }, [dataProvider, resource]);

  if (!data || data.length === 0) {
    return {
      data,
      chartData: [],
      granularityMinutes: 60,
      seriesKeys: [],
      seriesNames: {},
    };
  }

  const timeBuckets = data.map((s) => parseTimeBucket(s.timeBucket));
  const minTime = new Date(Math.min(...timeBuckets.map((t) => t.getTime())));
  const maxTime = new Date(Math.max(...timeBuckets.map((t) => t.getTime())));

  const granularityMinutes = calculateGranularity(minTime, maxTime);
  const allBuckets = generateTimeBuckets(minTime, maxTime, granularityMinutes);

  const seriesKeys = Array.from(
    new Set(data.map((item) => String(item[dataKey]))),
  );
  const seriesNames: Record<string, string> = {};
  data.forEach((item) => {
    const key = String(item[dataKey]);
    if (!seriesNames[key]) {
      seriesNames[key] = key;
    }
  });

  const dataByBucket: Record<string, Record<string, number>> = {};
  allBuckets.forEach((bucket) => {
    dataByBucket[bucket] = {};
    seriesKeys.forEach((key) => {
      dataByBucket[bucket][key] = initialValue;
    });
  });

  data.forEach((item) => {
    const statTime = parseTimeBucket(item.timeBucket);
    const bucketTime = snapToGranularity(statTime, granularityMinutes);
    const bucketStr = formatTimeBucketString(
      bucketTime,
      granularityMinutes < 60,
    );

    if (dataByBucket[bucketStr]) {
      const key = String(item[dataKey]);
      dataByBucket[bucketStr][key] = aggregateFn(
        dataByBucket[bucketStr][key],
        item,
      );
    }
  });

  const chartData: ChartDataPoint[] = allBuckets.map((bucket) => {
    const point: ChartDataPoint = { hour: bucket };
    seriesKeys.forEach((key) => {
      point[key] = dataByBucket[bucket][key];
    });
    return point;
  });

  return {
    data,
    chartData,
    granularityMinutes,
    seriesKeys,
    seriesNames,
  };
};
