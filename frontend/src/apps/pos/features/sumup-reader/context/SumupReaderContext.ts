import { createContext } from "react";
import { SumupReader } from "@pos/api/schemas";
import { ReaderProblem } from "../utils/readerProblem";

export interface SumupReaderContextValue {
  /** Whether card payments via SumUp are configured at all. */
  enabled: boolean;
  selectedReaderId: string | undefined;
  selectedReader: SumupReader | undefined;
  /** null while unknown (loading or SumUp unreachable). */
  readers: SumupReader[] | null;
  loading: boolean;
  problem: ReaderProblem | null;
  selectReader: (id: string) => void;
  refresh: () => void;
  dialogOpen: boolean;
  openDialog: () => void;
  closeDialog: () => void;
}

export const SumupReaderContext = createContext<SumupReaderContextValue>({
  enabled: false,
  selectedReaderId: undefined,
  selectedReader: undefined,
  readers: null,
  loading: false,
  problem: null,
  selectReader: () => undefined,
  refresh: () => undefined,
  dialogOpen: false,
  openDialog: () => undefined,
  closeDialog: () => undefined,
});
