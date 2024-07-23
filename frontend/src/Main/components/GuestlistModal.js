import React, { useEffect, useState, useCallback } from "react";
import { FloatingLabel, Modal, Table, Avatar, Alert } from "flowbite-react";
import { fetchGuestListByProductId } from "../hooks/Api";
import { HiShoppingCart, HiInformationCircle, HiXCircle } from "react-icons/hi";
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
  const [error, setError] = useState(null);
  const [searchQuery, setSearchQuery] = useState("");
  const apiHost = useConfig().apiHost;
  const { token } = useAuth();

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
      } catch (error) {
        setError("Error fetching list entries: " + error.message);
        setGuestListEntries([]);
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
          <div className="w-1/4 bg-gray-100 p-4">
            <FloatingLabel
              variant="filled"
              label="Search..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              autoFocus={true}
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
          >
            <div className="text-xl mb-4">List for {product.name}</div>

            {error && (
              <Alert
                className="my-3"
                color="failure"
                icon={HiInformationCircle}
              >
                {error}
              </Alert>
            )}

            {guestListEntries.length === 0 && (
              <Alert className="my-3" color="warning" icon={HiXCircle}>
                No entries found
              </Alert>
            )}

            {guestListEntries.length > 0 && (
              <div className="space-y-4">
                <Table hoverable>
                  <Table.Head>
                    <Table.HeadCell></Table.HeadCell>
                    <Table.HeadCell>Name</Table.HeadCell>
                    <Table.HeadCell>Action</Table.HeadCell>
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
                                {highlightText(entry.name, searchQuery)}
                              </div>
                              <div className="text-sm">{entry.listName}</div>
                            </>
                          )}
                          {entry.code !== "" && (
                            <div className="text-3xl font-mono">
                              {highlightText(entry.code, searchQuery)}
                            </div>
                          )}
                        </Table.Cell>
                        <Table.Cell className="flex gap-5">
                          <MyButton
                            className="float"
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
                                key={i}
                                className="float"
                                {...(hasListItem(entry.id)
                                  ? { disabled: true }
                                  : {})}
                                onClick={() => handleAddToCart(entry, i + 1)}
                              >
                                +{i + 1}
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
  // Split the name by spaces
  const words = name.split(" ");
  let initials = "";

  // Iterate through each word and append the first letter to initials
  words.forEach((word) => {
    if (word.length > 0) {
      initials += word[0].toUpperCase();
    }
  });

  return initials;
};

export default GuestlistModal;