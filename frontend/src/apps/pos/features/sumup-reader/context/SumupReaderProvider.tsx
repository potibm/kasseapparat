import React, { useCallback, useEffect, useMemo, useState } from "react";
import { useConfig } from "@core/config/hooks/useConfig";
import { usePosApi } from "@pos/api/usePosApi";
import { SumupReader } from "@pos/api/schemas";
import {
  getCurrentReaderId,
  setCurrentReaderId,
} from "@core/localstorage/helper/local-storage-reader";
import { createLogger } from "@core/logger/logger";
import { SumupReaderContext } from "./SumupReaderContext";
import { findReaderProblem } from "../utils/readerProblem";
import SumupReaderDialog from "../components/SumupReaderDialog";

const log = createLogger("SumUp");

/** How often the reader status is refreshed while the POS is open. */
export const READER_REFRESH_MS = 60_000;

/**
 * Keeps track of the card reader this device uses: the selection (stored per
 * browser, like in the admin) and the readers SumUp knows with their status.
 * Renders the reader dialog, so it can be opened from anywhere in the POS.
 */
export const SumupReaderProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { sumupEnabled } = useConfig();
  const { fetchSumupReaders } = usePosApi();
  const [selectedReaderId, setSelectedReaderId] = useState(getCurrentReaderId);
  const [readers, setReaders] = useState<SumupReader[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);

  const refresh = useCallback(async () => {
    if (!sumupEnabled) return;

    setLoading(true);
    try {
      setReaders(await fetchSumupReaders());
    } catch (error: unknown) {
      log.error("Could not load the SumUp readers", error);
      setReaders(null);
    } finally {
      setLoading(false);
    }
  }, [sumupEnabled, fetchSumupReaders]);

  useEffect(() => {
    if (!sumupEnabled) return;

    const first = setTimeout(refresh, 0);
    const interval = setInterval(refresh, READER_REFRESH_MS);

    return () => {
      clearTimeout(first);
      clearInterval(interval);
    };
  }, [sumupEnabled, refresh]);

  const selectReader = useCallback((id: string) => {
    setCurrentReaderId(id);
    setSelectedReaderId(getCurrentReaderId());
  }, []);

  const openDialog = useCallback(() => {
    setDialogOpen(true);
    refresh();
  }, [refresh]);

  const closeDialog = useCallback(() => setDialogOpen(false), []);

  // No reader paired at all: nothing to choose here, that is set up in the admin.
  const available = sumupEnabled && readers?.length !== 0;

  const value = useMemo(
    () => ({
      enabled: available,
      selectedReaderId,
      selectedReader: readers?.find((r) => r.id === selectedReaderId),
      readers,
      loading,
      problem: available ? findReaderProblem(selectedReaderId, readers) : null,
      selectReader,
      refresh: () => void refresh(),
      dialogOpen,
      openDialog,
      closeDialog,
    }),
    [
      available,
      selectedReaderId,
      readers,
      loading,
      selectReader,
      refresh,
      dialogOpen,
      openDialog,
      closeDialog,
    ],
  );

  return (
    <SumupReaderContext value={value}>
      {children}
      {sumupEnabled && dialogOpen && <SumupReaderDialog />}
    </SumupReaderContext>
  );
};

export default SumupReaderProvider;
