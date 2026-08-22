export const parseTimeBucket = (bucket: string): Date => {
  return new Date(bucket.replace(" ", "T") + ":00Z");
};

export const formatTimeBucketString = (
  date: Date,
  includeMinutes: boolean = false,
): string => {
  const year = date.getUTCFullYear();
  const month = String(date.getUTCMonth() + 1).padStart(2, "0");
  const day = String(date.getUTCDate()).padStart(2, "0");
  const hours = String(date.getUTCHours()).padStart(2, "0");
  const minutes = String(date.getUTCMinutes()).padStart(2, "0");

  if (includeMinutes) {
    return `${year}-${month}-${day} ${hours}:${minutes}`;
  }
  return `${year}-${month}-${day} ${hours}:00`;
};

export const formatXAxisLabel = (
  value: string,
  granularityMinutes: number,
): string => {
  if (!value) return "";
  const parts = value.split(" ");
  if (parts.length < 2) return value;

  const date = parseTimeBucket(value);
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");

  if (granularityMinutes >= 60) {
    return `${month}-${day} ${hours}:00`;
  }
  return `${month}-${day} ${hours}:${minutes}`;
};

export const calculateGranularity = (minTime: Date, maxTime: Date): number => {
  const spanHours = (maxTime.getTime() - minTime.getTime()) / (1000 * 60 * 60);
  if (spanHours < 12) return 15;
  if (spanHours < 24) return 30;
  return 60;
};

export const generateTimeBuckets = (
  minTime: Date,
  maxTime: Date,
  granularityMinutes: number,
): string[] => {
  const buckets: string[] = [];
  const current = new Date(minTime);
  current.setUTCSeconds(0, 0);
  const currentMinutes = current.getUTCMinutes();
  current.setUTCMinutes(
    Math.floor(currentMinutes / granularityMinutes) * granularityMinutes,
  );

  while (current <= maxTime) {
    buckets.push(formatTimeBucketString(current, granularityMinutes < 60));
    current.setUTCMinutes(current.getUTCMinutes() + granularityMinutes);
  }

  return buckets;
};

export const snapToGranularity = (
  date: Date,
  granularityMinutes: number,
): Date => {
  const snapped = new Date(date);
  const minutes = snapped.getUTCMinutes();
  snapped.setUTCMinutes(
    Math.floor(minutes / granularityMinutes) * granularityMinutes,
  );
  return snapped;
};
