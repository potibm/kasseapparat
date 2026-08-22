import { useEffect, useState } from "react";
import { useDataProvider } from "react-admin";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Admin");

interface UseStatsDataResult<T> {
  data: T[] | null;
  loading: boolean;
  error: boolean;
}

export const useStatsData = <T>(
  resource: string,
  sortField: string = "name",
  sortOrder: "ASC" | "DESC" = "ASC",
): UseStatsDataResult<T> => {
  const [data, setData] = useState<T[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const dataProvider = useDataProvider();

  useEffect(() => {
    setLoading(true);
    setError(false);

    dataProvider
      .getList<T>(resource, {
        pagination: { page: 1, perPage: 100 },
        sort: { field: sortField, order: sortOrder },
        filter: {},
      })
      .then(({ data }) => {
        setData(data);
        setLoading(false);
      })
      .catch((err) => {
        log.error(`${resource} fetch failed`, err);
        setData([]);
        setLoading(false);
        setError(true);
      });
  }, [dataProvider, resource, sortField, sortOrder]);

  return { data, loading, error };
};
