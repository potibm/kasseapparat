import React, { useEffect, useState, useCallback } from "react";
import { FloatingLabel, Modal, ModalBody } from "flowbite-react";
import { fetchGuestlistByProductId } from "../../hooks/Api";
import { HiOutlineX } from "react-icons/hi";
import PropTypes from "prop-types";
import SidebarKeyboard from "./components/SidebarKeyboard";
import { useConfig } from "../../../provider/ConfigProvider";
import { useAuth } from "../../../Auth/provider/AuthProvider";
import MyButton from "../MyButton";
import GuestlistResultTable from "./components/ResultTable";

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
  const { token } = useAuth();

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
    async () => {
      setLoading(true);
      try {
        let response = await fetchGuestlistByProductId(
          apiHost,
          token,
          product.id,
          searchQuery,
        );
        if (response === null) {
          response = [];
        }
        setGuestlistEntries(response);
        setError(null);
        setLoading(false);
        setLoadedSearchQuery(searchQuery);
      } catch (error) {
        setError("Error fetching list entries: " + error.message);
        setGuestlistEntries([]);
        setLoading(false);
        setLoadedSearchQuery("");
      }
    },
    [product.id, searchQuery, apiHost, token],
  );

  useEffect(() => {
    if (isOpen) {
      fetchGuestEntries(searchQuery);
    }
  }, [isOpen, searchQuery, fetchGuestEntries]);

  useEffect(() => {
    if (!isOpen) {
      setSearchQuery("");
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
        <div className="flex h-full">
          {/* Sidebar */}
          <div className="w-4/12 bg-gray-100 dark:bg-gray-900 p-4">
            <FloatingLabel
              variant="filled"
              label="Search..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              autoFocus={hasCodes}
            />

            <SidebarKeyboard term={searchQuery} setTerm={setSearchQuery} />

            <MyButton
              className="w-full mt-5"
              color="warning"
              onClick={handleManualAddToCart}
            >
              Manual
            </MyButton>
          </div>
          <div
            className="w-3/4 p-4 overflow-y-auto dark:text-white"
            style={{ maxHeight: "calc(100vh - 10rem)" }}
            id="results"
          >
            <div className="text-xl mb-4 flex justify-between items-center">
              <span>List for {product.name}</span>
              <MyButton onClick={onClose} color="gray">
                <HiOutlineX />
              </MyButton>
            </div>

            <div
              className="relative"
              style={{ maxHeight: "calc(100vh - 10rem)", minHeight: "200px" }}
            >
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
