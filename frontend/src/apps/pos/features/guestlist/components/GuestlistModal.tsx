import React, { useEffect, useState, useCallback, useRef } from "react";
import { Modal, ModalBody } from "flowbite-react";
import { HiCheckCircle, HiOutlineX, HiQrcode } from "react-icons/hi";
import SidebarKeyboard from "./_internal/SidebarKeyboard";
import GuestlistResultTable from "./_internal/GuestlistResultTable";
import { Product as ProductType, Guest as GuestType } from "@pos/api/schemas";
import { usePosApi } from "@pos/api/usePosApi";
import { useConfig } from "@core/config/hooks/useConfig";
import { canAddMore, availableStock } from "@pos/features/stock/utils/stock";
import { findScannedGuest } from "../utils/scan";

const SEARCH_DEBOUNCE_MS = 200;

type ScanStatus =
  | { kind: "added"; code: string }
  | { kind: "inCart"; code: string }
  | { kind: "chooseGuests"; code: string }
  | { kind: "noStock"; code: string };

interface GuestlistModalProps {
  isOpen: boolean;
  onClose: () => void;
  product: ProductType;
  addToCart: (
    product: ProductType,
    quantity: number,
    listEntry: GuestType | null,
  ) => void;
  hasListItem: (guest: GuestType) => boolean;
  quantityInCart?: number;
}

const GuestlistModal: React.FC<GuestlistModalProps> = ({
  isOpen,
  onClose,
  product,
  addToCart,
  hasListItem,
  quantityInCart = 0,
}) => {
  const [guestlistEntries, setGuestlistEntries] = useState<GuestType[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [loadedSearchQuery, setLoadedSearchQuery] = useState<string>("");
  // The guest just checked in whose arrival note is shown in their result card.
  const [noteGuestId, setNoteGuestId] = useState<number | null>(null);
  const [noteFromScan, setNoteFromScan] = useState(false);
  const hasCodes = product.guestlists?.some((list) => list.typeCode) ?? false;

  const [scanStatus, setScanStatus] = useState<ScanStatus | null>(null);
  const [searchFocused, setSearchFocused] = useState(hasCodes);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const { fetchGuestlistByProductId } = usePosApi();
  const { outOfStockBehavior } = useConfig();

  // Scanners type into the search field, so it must be focused for them to work.
  // (Done here instead of the autoFocus attribute, which hurts accessibility.)
  useEffect(() => {
    if (!hasCodes) return;
    const id = window.setTimeout(() => searchInputRef.current?.focus(), 0);
    return () => window.clearTimeout(id);
  }, [hasCodes]);

  const handleAddToCart = (listEntry: GuestType, additionalGuests: number) => {
    if (
      !canAddMore(product, quantityInCart, outOfStockBehavior) ||
      availableStock(product, quantityInCart) < additionalGuests + 1
    ) {
      return;
    }

    addToCart(product, additionalGuests + 1, listEntry);

    if (listEntry.arrivalNote) {
      setNoteGuestId(listEntry.id);
    } else {
      onClose();
    }
  };

  /**
   * A scanned ticket checks its guest in right away and keeps the list open for the
   * next scan. Guests who may bring companions still need a tap on +1/+2.
   */
  const checkInScannedGuest = useCallback(
    (guest: GuestType) => {
      const code = guest.code ?? guest.name;

      if (hasListItem(guest)) {
        setScanStatus({ kind: "inCart", code });
        return;
      }
      if (guest.additionalGuests > 0) {
        setScanStatus({ kind: "chooseGuests", code });
        return;
      }
      if (
        !canAddMore(product, quantityInCart, outOfStockBehavior) ||
        availableStock(product, quantityInCart) < 1
      ) {
        setScanStatus({ kind: "noStock", code });
        return;
      }

      addToCart(product, 1, guest);

      if (guest.arrivalNote) {
        setNoteGuestId(guest.id);
        setNoteFromScan(true);
      } else {
        setScanStatus({ kind: "added", code });
        setSearchQuery("");
      }
    },
    [addToCart, hasListItem, product, quantityInCart, outOfStockBehavior],
  );

  // The search callback reads the latest version through this ref, so adding a
  // guest to the cart does not restart the search (and hide an arrival note).
  const checkInScannedRef = useRef(checkInScannedGuest);
  useEffect(() => {
    checkInScannedRef.current = checkInScannedGuest;
  }, [checkInScannedGuest]);

  const handleAcknowledgeNote = () => {
    setNoteGuestId(null);

    if (noteFromScan) {
      // Back to scanning: the guest is in the cart, the list stays open.
      setNoteFromScan(false);
      setSearchQuery("");
      searchInputRef.current?.focus();
    } else {
      onClose();
    }
  };

  const handleManualAddToCart = () => {
    if (!canAddMore(product, quantityInCart, outOfStockBehavior)) {
      return;
    }

    addToCart(product, 1, null);
    onClose();
  };

  const fetchGuestEntries = useCallback(
    async (query: string, signal: AbortSignal) => {
      setLoading(true);
      try {
        const response = await fetchGuestlistByProductId(
          product.id,
          query,
          signal,
        );
        // A newer search has started meanwhile: drop this (stale) result.
        if (signal.aborted) return;
        setGuestlistEntries(response);
        setError(null);
        setLoadedSearchQuery(query);

        const scannedGuest = hasCodes
          ? findScannedGuest(response, query)
          : null;
        if (scannedGuest) checkInScannedRef.current(scannedGuest);
      } catch (error: unknown) {
        if (signal.aborted) return;
        setError("Error fetching list entries: " + (error as Error).message);
        setGuestlistEntries([]);
        setLoadedSearchQuery("");
      } finally {
        if (!signal.aborted) setLoading(false);
      }
    },
    [product.id, fetchGuestlistByProductId, hasCodes],
  );

  useEffect(() => {
    if (!isOpen) return;

    // Debounce typing and cancel the previous request, so a slow response for an
    // older search term can never overwrite the results of the current one.
    const controller = new AbortController();
    const handle = setTimeout(
      () => {
        fetchGuestEntries(searchQuery, controller.signal);
        setNoteGuestId(null);
      },
      searchQuery === "" ? 0 : SEARCH_DEBOUNCE_MS,
    );

    return () => {
      clearTimeout(handle);
      controller.abort();
    };
  }, [isOpen, searchQuery, fetchGuestEntries]);

  useEffect(() => {
    if (!isOpen) {
      const handle = requestAnimationFrame(() => {
        setSearchQuery("");
        setNoteGuestId(null);
        setNoteFromScan(false);
        setScanStatus(null);
      });
      return () => cancelAnimationFrame(handle);
    }
  }, [isOpen]);

  const canSellWithoutEntry = canAddMore(
    product,
    quantityInCart,
    outOfStockBehavior,
  );

  return (
    <Modal
      show={isOpen}
      onClose={onClose}
      position="top-center"
      size="7xl"
      dismissible
    >
      <ModalBody className="overflow-hidden p-0">
        <div className="flex h-[calc(100vh-6rem)] flex-col dark:text-white">
          <div className="flex h-16 shrink-0 items-center justify-between border-b border-gray-200 pr-4 pl-6 dark:border-gray-700">
            <h2 className="text-xl font-bold">List for {product.name}</h2>
            <button
              type="button"
              aria-label="Close guest list"
              onClick={onClose}
              className="pos-pressable flex h-12 w-12 items-center justify-center rounded-xl border border-gray-300 dark:border-gray-600"
            >
              <HiOutlineX className="h-6 w-6" aria-hidden="true" />
            </button>
          </div>

          <div className="flex min-h-0 flex-1">
            {/* Search and letters */}
            <div className="flex w-[26rem] shrink-0 flex-col gap-4 bg-gray-100 p-5 dark:bg-gray-900">
              <div className="relative">
                <label
                  htmlFor="guestlist-search"
                  className="absolute top-2 left-4 text-xs font-semibold text-gray-600 dark:text-gray-300"
                >
                  {hasCodes ? "Scan ticket or search name" : "Search name"}
                </label>
                <input
                  ref={searchInputRef}
                  id="guestlist-search"
                  type="text"
                  autoComplete="off"
                  value={searchQuery}
                  onChange={(e) => {
                    setSearchQuery(e.target.value);
                    setScanStatus(null);
                  }}
                  onKeyDown={(e) => {
                    // Scanners may end a code with Enter; nothing to submit here.
                    if (e.key === "Enter") e.preventDefault();
                  }}
                  onFocus={() => setSearchFocused(true)}
                  onBlur={() => setSearchFocused(false)}
                  data-testid="guestlist-search-input"
                  className="h-16 w-full rounded-xl border-2 border-cyan-700 bg-white pt-5 pr-14 pl-4 text-2xl font-semibold tracking-wide text-gray-900 dark:bg-gray-800 dark:text-white"
                />
                {searchQuery !== "" && (
                  <button
                    type="button"
                    aria-label="Clear search term"
                    onClick={() => setSearchQuery("")}
                    className="absolute top-2.5 right-2 flex h-11 w-11 items-center justify-center rounded-lg bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-200"
                  >
                    <HiOutlineX className="h-5 w-5" aria-hidden="true" />
                  </button>
                )}
              </div>

              {hasCodes && (
                <ScannerStatus
                  status={scanStatus}
                  listening={searchFocused}
                  onFocusSearch={() => searchInputRef.current?.focus()}
                />
              )}

              <SidebarKeyboard term={searchQuery} setTerm={setSearchQuery} />

              <div className="flex-1" />

              <button
                type="button"
                className="pos-pressable h-14 rounded-xl border border-dashed border-gray-400 font-semibold text-gray-700 disabled:cursor-not-allowed disabled:opacity-50 dark:text-gray-200"
                aria-label="Add product to cart without guestlist entry"
                onClick={handleManualAddToCart}
                disabled={!canSellWithoutEntry}
              >
                Sell without list entry
              </button>
            </div>

            {/* Results */}
            <div
              className="relative min-w-0 flex-1 overflow-y-auto p-5"
              id="results"
            >
              <GuestlistResultTable
                loading={loading}
                error={error}
                guestlistEntries={guestlistEntries}
                onAddToCart={handleAddToCart}
                hasListItem={hasListItem}
                loadedSearchQuery={loadedSearchQuery}
                noteGuestId={noteGuestId}
                onAcknowledgeNote={handleAcknowledgeNote}
              />
            </div>
          </div>
        </div>
      </ModalBody>
    </Modal>
  );
};

interface ScannerStatusProps {
  status: ScanStatus | null;
  listening: boolean;
  onFocusSearch: () => void;
}

/**
 * Tells the cashier whether a scan will arrive: the scanner types into the search
 * field, so it is ready exactly while that field has the focus.
 */
const ScannerStatus: React.FC<ScannerStatusProps> = ({
  status,
  listening,
  onFocusSearch,
}) => {
  if (status) {
    const messages: Record<ScanStatus["kind"], [string, string]> = {
      added: [
        `${status.code} added – scan the next ticket`,
        "text-green-800 dark:text-green-300",
      ],
      inCart: [
        `${status.code} is already in the cart`,
        "text-amber-800 dark:text-amber-300",
      ],
      chooseGuests: [
        `${status.code}: choose the number of guests`,
        "text-cyan-800 dark:text-cyan-300",
      ],
      noStock: [
        `${status.code}: no tickets left`,
        "text-red-800 dark:text-red-300",
      ],
    };
    const [message, color] = messages[status.kind];

    return (
      <output
        data-testid="guestlist-scan-status"
        className={`-mt-2 flex items-center gap-2 text-sm font-semibold ${color}`}
      >
        <HiCheckCircle className="h-5 w-5 shrink-0" aria-hidden="true" />
        {message}
      </output>
    );
  }

  if (listening) {
    return (
      <output
        data-testid="guestlist-scan-status"
        className="-mt-2 flex items-center gap-2 text-sm text-green-800 dark:text-green-300"
      >
        <HiQrcode className="h-5 w-5 shrink-0" aria-hidden="true" />
        <span>
          <b>Ready to scan</b> · the search field is listening
        </span>
      </output>
    );
  }

  return (
    <button
      type="button"
      data-testid="guestlist-scan-status"
      onClick={onFocusSearch}
      className="-mt-2 flex items-center gap-2 text-left text-sm font-semibold text-amber-800 dark:text-amber-300"
    >
      <HiQrcode className="h-5 w-5 shrink-0" aria-hidden="true" />
      Not ready to scan – tap here first
    </button>
  );
};

export default GuestlistModal;
