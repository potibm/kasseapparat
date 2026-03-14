import React, { useMemo } from "react";
import { Filter, SelectInput, FilterProps } from "react-admin";

const getOldestTimeOptions = () => {
  const now = new Date();
  now.setSeconds(0, 0);
  const baseTime = now.getTime();

  const toISO = (msAgo: number) => new Date(baseTime - msAgo).toISOString();

  return [
    { id: toISO(10 * 60 * 1000), name: "Last 10 minutes" },
    { id: toISO(60 * 60 * 1000), name: "Last 1 hour" },
    { id: toISO(24 * 60 * 60 * 1000), name: "Last 24 hours" },
    { id: toISO(48 * 60 * 60 * 1000), name: "Last 48 hours" },
    { id: toISO(7 * 24 * 60 * 60 * 1000), name: "Last 1 week" },
    { id: toISO(30 * 24 * 60 * 60 * 1000), name: "Last 1 month" },
  ];
};

export const SumupTransactionTimeRangeFilter: React.FC<
  React.PropsWithoutRef<Omit<FilterProps, "children">>
> = (props) => {
  const choices = useMemo(() => getOldestTimeOptions(), []);

  return (
    <Filter {...props}>
      <SelectInput
        source="oldest_time"
        label="Time Range"
        choices={choices}
        alwaysOn
      />
    </Filter>
  );
};

export default SumupTransactionTimeRangeFilter;
