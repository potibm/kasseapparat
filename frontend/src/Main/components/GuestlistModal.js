import React, { useEffect, useState, useCallback } from "react";
import { FloatingLabel, Modal, Table, Avatar, Button } from "flowbite-react";
import { fetchGuestListByProductId } from "../hooks/Api";
import { HiShoppingCart } from "react-icons/hi";
import PropTypes from "prop-types";
import SidebarKeyboard from "./SidebarKeyboard";

const API_HOST = process.env.REACT_APP_API_HOST;

const GuestlistModal = ({
  isOpen,
  onClose,
  product,
  addToCart,
  hasListItem,
}) => {
  const [guestListEntries, setGuestListEntries] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");

  const handleAddToCart = (listEntry, additionalGuests) => {
    addToCart(product, additionalGuests + 1, listEntry);
    onClose(); // close the modal
  };

  const fetchGuestEntries = useCallback(
    async (query = "") => {
      try {
        let response = await fetchGuestListByProductId(
          API_HOST,
          product.id,
          searchQuery,
        );
        if (response === null) {
          response = [];
        }
        setGuestListEntries(response);
      } catch (error) {
        console.error("Error fetching guest entries:", error);
        setGuestListEntries([]);
      }
    },
    [product.id, searchQuery],
  );

  useEffect(() => {
    if (isOpen) {
      fetchGuestEntries(searchQuery);
    }
  }, [isOpen, searchQuery, fetchGuestEntries]);

  return (
    <Modal
      show={isOpen}
      onClose={onClose}
      position="top-center"
      size="7xl"
      dismissible
    >
      <Modal.Header>Gästeliste für {product.name}</Modal.Header>
      <Modal.Body>
        <div className="flex h-full">
          {/* Sidebar */}
          <div className="w-1/4 bg-gray-100 p-4">
            <h2 className="text-lg font-semibold mb-4">Search</h2>
            <FloatingLabel
              variant="filled"
              label="Search..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />

            <SidebarKeyboard term={searchQuery} setTerm={setSearchQuery} />
          </div>
          <div
            className="w-3/4 p-4 overflow-y-auto"
            style={{ maxHeight: "calc(100vh - 10rem)" }}
          >
            <h2 className="text-lg font-semibold mb-4">Content</h2>
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
                        <Button
                          className="float"
                          {...(hasListItem(entry.id) ? { disabled: true } : {})}
                          onClick={() => handleAddToCart(entry, 0)}
                        >
                          <HiShoppingCart />
                        </Button>
                        {Array.from(
                          { length: entry.additionalGuests },
                          (_, i) => (
                            <Button
                              key={i}
                              className="float"
                              {...(hasListItem(entry.id)
                                ? { disabled: true }
                                : {})}
                              onClick={() => handleAddToCart(entry, i + 1)}
                            >
                              +{i + 1}
                            </Button>
                          ),
                        )}
                      </Table.Cell>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table>
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
