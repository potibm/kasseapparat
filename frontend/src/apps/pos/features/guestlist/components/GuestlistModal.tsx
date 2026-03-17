import React, { useEffect, useState, useCallback } from "react";
import { FloatingLabel, Modal, ModalBody } from "flowbite-react";
import { fetchGuestlistByProductId } from "../../../utils/api";
import { HiOutlineX } from "react-icons/hi";
import SidebarKeyboard from "./_internal/SidebarKeyboard";
import { useConfig } from "../../../../../core/config/providers/ConfigProvider";
import { useAuth } from "../../auth/providers/auth-provider";
import Button from "../../../components/Button";
import GuestlistResultTable from "./_internal/GuestlistResultTable";
import {
  Product as ProductType,
  Guest as GuestType,
} from "@pos/utils/api.schemas";
import GuestlistArrivalNoteModal from "./_internal/GuestlistArrivalNoteModal";

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
}

const GuestlistModal: React.FC<GuestlistModalProps> = ({
  isOpen,
  onClose,
  product,
  addToCart,
  hasListItem,
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

  const { apiHost } = useConfig();
  const { getSafeToken } = useAuth();

  const hasCodes = product.guestlists?.some((list) => list.typeCode) ?? false;

  const handleAddToCart = (listEntry: GuestType, additionalGuests: number) => {
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
    addToCart(product, 1, null);
    onClose();
  };

  const fetchGuestEntries = useCallback(
    async (query = "") => {
      setLoading(true);
      try {
        const token = await getSafeToken();

        const response = await fetchGuestlistByProductId(
          apiHost,
          token,
          product.id,
          query,
        );
        setGuestlistEntries(response);
        setError(null);
        setLoadedSearchQuery(query);
      } catch (error: unknown) {
        setError("Error fetching list entries: " + (error as Error).message);
        setGuestlistEntries([]);
        setLoadedSearchQuery("");
      } finally {
        setLoading(false);
      }
    },
    [product.id, apiHost, getSafeToken],
  );

  useEffect(() => {
    if (isOpen) {
      const handle = requestAnimationFrame(() => {
        fetchGuestEntries(searchQuery);
      });
      return () => cancelAnimationFrame(handle);
    }
  }, [isOpen, searchQuery, fetchGuestEntries]);

  useEffect(() => {
    if (!isOpen) {
      const handle = requestAnimationFrame(() => {
        setSearchQuery("");
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
              autoFocus={hasCodes}
            />

            <SidebarKeyboard term={searchQuery} setTerm={setSearchQuery} />

            <Button
              className="w-full mt-5"
              color="alternative"
              onClick={handleManualAddToCart}
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
