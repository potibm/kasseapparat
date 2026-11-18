import React, { useEffect, useState, useCallback } from "react";
import { FloatingLabel, Modal, ModalBody } from "flowbite-react";
import { HiOutlineX } from "react-icons/hi";
import SidebarKeyboard from "./_internal/SidebarKeyboard";
import Button from "../../../components/Button";
import GuestlistResultTable from "./_internal/GuestlistResultTable";
import { Product as ProductType, Guest as GuestType } from "@pos/api/schemas";
import GuestlistArrivalNoteModal from "./_internal/GuestlistArrivalNoteModal";
import { usePosApi } from "@pos/api/usePosApi";
import { useConfig } from "@core/config/hooks/useConfig";
import { canAddMore, availableStock } from "@pos/features/stock/utils/stock";

const SEARCH_DEBOUNCE_MS = 200;

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
  const [noteToDisplay, setNoteToDisplay] = useState<{
    name: string;
    note: string;
  } | null>(null);

  const { fetchGuestlistByProductId } = usePosApi();
  const { outOfStockBehavior } = useConfig();

  const hasCodes = product.guestlists?.some((list) => list.typeCode) ?? false;

  const handleAddToCart = (listEntry: GuestType, additionalGuests: number) => {
    if (
      !canAddMore(product, quantityInCart, outOfStockBehavior) ||
      availableStock(product, quantityInCart) < additionalGuests + 1
    ) {
      return;
    }

    addToCart(product, additionalGuests + 1, listEntry);

    if (listEntry.arrivalNote) {
      setNoteToDisplay({
        name: listEntry.name,
        note: listEntry.arrivalNote,
      });
    } else {
      onClose();
    }
  };

  const handleCloseNote = () => {
    setNoteToDisplay(null);
    onClose();
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
      } catch (error: unknown) {
        if (signal.aborted) return;
        setError("Error fetching list entries: " + (error as Error).message);
        setGuestlistEntries([]);
        setLoadedSearchQuery("");
      } finally {
        if (!signal.aborted) setLoading(false);
      }
    },
    [product.id, fetchGuestlistByProductId],
  );

  useEffect(() => {
    if (!isOpen) return;

    // Debounce typing and cancel the previous request, so a slow response for an
    // older search term can never overwrite the results of the current one.
    const controller = new AbortController();
    const handle = setTimeout(
      () => {
        fetchGuestEntries(searchQuery, controller.signal);
        setNoteToDisplay(null);
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
        setNoteToDisplay(null);
      });
      return () => cancelAnimationFrame(handle);
    }
  }, [isOpen]);

  return (
    <Modal
      show={isOpen}
      onClose={onClose}
      position="top-center"
      size="7xl"
      dismissible
    >
      <ModalBody className="overflow-hidden">
        {noteToDisplay && (
          <GuestlistArrivalNoteModal
            isOpen={true}
            onClose={handleCloseNote}
            arrivalNote={noteToDisplay.note}
            name={noteToDisplay.name}
          />
        )}
        <div className="flex h-[calc(100vh-6rem)]">
          {/* Sidebar */}
          <div className="w-4/12 bg-gray-100 dark:bg-gray-900 p-4">
            <FloatingLabel
              className="mb-4"
              variant="filled"
              label="Search..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              data-testid="guestlist-search-input"
              autoFocus={hasCodes}
            />

            <SidebarKeyboard term={searchQuery} setTerm={setSearchQuery} />

            <Button
              className="w-full mt-5"
              color="alternative"
              aria-label="Add product to cart without guestlist entry"
              onClick={handleManualAddToCart}
              disabled={
                !canAddMore(product, quantityInCart, outOfStockBehavior)
              }
            >
              Manual
            </Button>
          </div>
          <div
            className="w-3/4 p-4 overflow-y-auto dark:text-white"
            id="results"
          >
            <div className="text-xl mb-4 flex justify-between items-center">
              <span>List for {product.name}</span>
              <Button onClick={onClose} color="gray">
                <HiOutlineX />
              </Button>
            </div>

            <div className="relative">
              <GuestlistResultTable
                loading={loading}
                error={error}
                guestlistEntries={guestlistEntries}
                onAddToCart={handleAddToCart}
                hasListItem={hasListItem}
                loadedSearchQuery={loadedSearchQuery}
              />
            </div>
          </div>
        </div>
      </ModalBody>
    </Modal>
  );
};

export default GuestlistModal;
