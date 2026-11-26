import { createContext } from "react";
import { SumupReader } from "@pos/api/schemas";
import { ReaderProblem } from "../utils/readerProblem";

export interface SumupReaderContextValue {
  /**
   * Whether the POS offers card payments via SumUp: configured, and SumUp knows at
   * least one reader (or the list is not known yet). Pairing a reader is done in the
   * admin, so without any the POS shows neither the button nor a hint.
   */
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
