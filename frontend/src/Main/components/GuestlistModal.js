import React, { useEffect, useState, useCallback } from "react";
import {
  FloatingLabel,
  Modal,
  Table,
  Avatar,
  Alert,
  Spinner,
} from "flowbite-react";
import { fetchGuestListByProductId } from "../hooks/Api";
import {
  HiShoppingCart,
  HiInformationCircle,
  HiXCircle,
  HiOutlineX,
} from "react-icons/hi";
import PropTypes from "prop-types";
import SidebarKeyboard from "./SidebarKeyboard";
import { useConfig } from "../../provider/ConfigProvider";
import { useAuth } from "../../Auth/provider/AuthProvider";
import MyButton from "./MyButton";

const GuestlistModal = ({
  isOpen,
  onClose,
  product,
  addToCart,
  hasListItem,
}) => {
  const [guestListEntries, setGuestListEntries] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [loadedSearchQuery, setLoadedSearchQuery] = useState("");
  const apiHost = useConfig().apiHost;
  const { token } = useAuth();

  let hasCodes = false;
  product.lists.forEach((list) => {
    if (list.typeCode) {
      hasCodes = true;
    }
  });

  const handleAddToCart = (listEntry, additionalGuests) => {
    addToCart(product, additionalGuests + 1, listEntry);
    onClose(); // close the modal
  };

  const handleManualAddToCart = () => {
    addToCart(product, 0, null);
    onClose(); // close the modal
  };

  const fetchGuestEntries = useCallback(
    async (query = "") => {
      setLoading(true);
      try {
        let response = await fetchGuestListByProductId(
          apiHost,
          token,
          product.id,
          searchQuery,
        );
        if (response === null) {
          response = [];
        }
        setGuestListEntries(response);
        setError(null);
        setLoading(false);
        setLoadedSearchQuery(searchQuery);
      } catch (error) {
        setError("Error fetching list entries: " + error.message);
        setGuestListEntries([]);
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
      <Modal.Body className="overflow-hidden">
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
            className="w-3/4 p-4 overflow-y-auto"
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
              {loading && (
                <div className="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10">
                  <Spinner size="xl" />
                </div>
              )}

              {error && (
                <Alert
                  className="my-3"
                  color="failure"
                  icon={HiInformationCircle}
                >
                  {error}
                </Alert>
              )}

              {!loading && guestListEntries.length === 0 && (
                <Alert className="my-3" color="warning" icon={HiXCircle}>
                  No entries found
                </Alert>
              )}

              {guestListEntries.length > 0 && (
                <div className="space-y-4">
                  <Table hoverable>
                    <Table.Head>
                      <Table.HeadCell className="w-1/12"></Table.HeadCell>
                      <Table.HeadCell className="w-5/12">Name</Table.HeadCell>
                      <Table.HeadCell className="w-6/12">Action</Table.HeadCell>
                    </Table.Head>
                    <Table.Body className="divide-y">
                      {guestListEntries.map((entry) => (
                        <Table.Row key={entry.id}>
                          <Table.Cell>
                            {!entry.code && (
                              <Avatar
                                placeholderInitials={getInitials(entry.name)}
                                size="md"
                                rounded
                              />
                            )}
                          </Table.Cell>
                          <Table.Cell className="">
                            {!entry.code && (
                              <>
                                <div className="text-xl">
                                  {highlightText(entry.name, loadedSearchQuery)}
                                </div>
                                <div className="text-sm">{entry.listName}</div>
                              </>
                            )}
                            {entry.code !== "" && (
                              <div className="text-3xl font-mono">
                                {highlightText(entry.code, loadedSearchQuery)}
                              </div>
                            )}
                          </Table.Cell>
                          <Table.Cell className="flex gap-5">
                            <MyButton
                              className="float"
                              key={0}
                              {...(hasListItem(entry.id)
                                ? { disabled: true }
                                : {})}
                              onClick={() => handleAddToCart(entry, 0)}
                            >
                              <HiShoppingCart />
                            </MyButton>
                            {Array.from(
                              { length: entry.additionalGuests },
                              (_, i) => (
                                <MyButton
                                  key={i + 1}
                                  className="float"
                                  {...(hasListItem(entry.id)
                                    ? { disabled: true }
                                    : {})}
                                  onClick={() => handleAddToCart(entry, i + 1)}
                                >
                                  <div className="text-xs">+{i + 1}</div>
                                </MyButton>
                              ),
                            )}
                          </Table.Cell>
                        </Table.Row>
                      ))}
                    </Table.Body>
                  </Table>
                </div>
              )}
            </div>
          </div>
        </div>
      </Modal.Body>
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

const highlightText = (text, highlight) => {
  if (!text) {
    return "";
  }
  if (!highlight.trim()) {
    return text;
  }
  const regex = new RegExp(`(${highlight})`, "gi");
  const parts = text.split(regex);
  return (
    <>
      {parts.map((part, i) =>
        regex.test(part) ? (
          <span key={i} className="font-bold underline">
            {part}
          </span>
        ) : (
          part
        ),
      )}
    </>
  );
};

const getInitials = (name) => {
  // Remove all non-alphabetical characters except spaces
  const cleanedName = name.replace(/[^a-zA-Z\s]/g, "").trim();

  // Split the cleaned name into words
  const words = cleanedName.split(" ").filter((word) => word.length > 0);

  // If there's only one word, take the first letter twice
  if (words.length === 1) {
    return words[0][0].toUpperCase();
  }

  // For multiple words, take the first letter of the first and last word
  const firstInitial = words[0][0].toUpperCase();
  const lastInitial = words[words.length - 1][0].toUpperCase();

  return firstInitial + lastInitial;
};

export default GuestlistModal;
