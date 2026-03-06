import React, { useEffect, useState, useCallback } from "react";
import { FloatingLabel, Modal, ModalBody } from "flowbite-react";
import { fetchGuestlistByProductId } from "../../../utils/api";
import { HiOutlineX } from "react-icons/hi";
import PropTypes from "prop-types";
import SidebarKeyboard from "./_internal/SidebarKeyboard";
import { useConfig } from "../../../../../core/config/providers/config-provider";
import { useAuth } from "../../auth/providers/auth-provider";
import Button from "../../../components/Button";
import GuestlistResultTable from "./_internal/ResultTable";

const GuestlistModal = ({
  isOpen,
  onClose,
  product,
  addToCart,
  hasListItem,
}) => {
  const [guestlistEntries, setGuestlistEntries] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [loadedSearchQuery, setLoadedSearchQuery] = useState("");
  const apiHost = useConfig().apiHost;
  const { getToken } = useAuth();

  const hasCodes = product.guestlists.some((list) => list.typeCode);

  const handleAddToCart = (listEntry, additionalGuests) => {
    addToCart(product, additionalGuests + 1, listEntry);
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
        let response = await fetchGuestlistByProductId(
          apiHost,
          await getToken(),
          product.id,
          query,
        );
        if (response === null) {
          response = [];
        }
        setGuestlistEntries(response);
        setError(null);
        setLoading(false);
        setLoadedSearchQuery(query);
      } catch (error) {
        setError("Error fetching list entries: " + error.message);
        setGuestlistEntries([]);
        setLoading(false);
        setLoadedSearchQuery("");
      }
    },
    [product.id, apiHost, getToken],
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
                error={error}
                loading={loading}
                loadedSearchQuery={loadedSearchQuery}
                guestlistEntries={guestlistEntries}
                hasListItem={hasListItem}
                onAddToCart={handleAddToCart}
                onClose={onClose}
              />
            </div>
          </div>
        </div>
      </ModalBody>
    </Modal>
  );
};

GuestlistModal.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  product: PropTypes.object.isRequired,
  addToCart: PropTypes.func.isRequired,
  hasListItem: PropTypes.func.isRequired,
};

export default GuestlistModal;
