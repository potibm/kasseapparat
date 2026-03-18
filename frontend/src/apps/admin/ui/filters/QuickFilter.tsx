import { Chip, ChipProps } from "@mui/material";

// T ist der Typ des Filterwerts (z.B. boolean, string, etc.)
// Wir nutzen Omit, um einen sauberen ChipProps-Typ ohne Konflikte zu haben
interface QuickFilterProps<T = unknown> extends Omit<
  ChipProps,
  "defaultValue"
> {
  label: string;
  source?: string;
  defaultValue?: T;
}

export const QuickFilter = <T,>({
  label,
  source: _source,
  defaultValue: _defaultValue,
  ...rest
}: QuickFilterProps<T>) => {
  return <Chip sx={{ marginBottom: 1 }} label={label} {...rest} />;
};

export default QuickFilter;
