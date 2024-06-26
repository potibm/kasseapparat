import React, { useEffect, useState, useCallback } from "react";
import { FloatingLabel, Modal, Table } from "flowbite-react";
import { fetchGuestListByProductId } from "../hooks/Api";
import PropTypes from "prop-types";
import SidebarKeyboard from "./SidebarKeyboard";

const API_HOST = process.env.REACT_APP_API_HOST;

const GuestlistModal = ({ isOpen, onClose, product }) => {
  const [guestListEntries, setGuestListEntries] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");

  const fetchGuestEntries = useCallback(
    async (query = "") => {
      try {
        const response = await fetchGuestListByProductId(
          API_HOST,
          product.id,
          searchQuery,
        );
        setGuestListEntries(response);
      } catch (error) {
        console.error("Error fetching guest entries:", error);
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

            <SidebarKeyboard
              term={searchQuery}
              setTerm={setSearchQuery}
              /*              addToSearchTerm={handleAddToSearchTerm}
              removeFromSearchTerm={handleRemoveFromSearchTerm}
              removeSearchTerm={handleRemoveSearchTerm} */
            />
          </div>
          {/* Scrollable Content */}
          <div
            className="w-3/4 p-4 overflow-y-auto"
            style={{ maxHeight: "calc(100vh - 10rem)" }}
          >
            <h2 className="text-lg font-semibold mb-4">Content</h2>
            <div className="space-y-4">
              <Table hoverable>
                <Table.Head>
                  <Table.HeadCell>Id</Table.HeadCell>
                  <Table.HeadCell>Name</Table.HeadCell>
                  <Table.HeadCell>List</Table.HeadCell>
                  <Table.HeadCell>Code</Table.HeadCell>
                </Table.Head>
                <Table.Body className="divide-y">
                  {guestListEntries.map((entry) => (
                    <Table.Row key={entry.id}>
                      <Table.Cell>{entry.id}</Table.Cell>
                      <Table.Cell>
                        {highlightText(entry.name, searchQuery)}
                      </Table.Cell>
                      <Table.Cell>{entry.list}</Table.Cell>
                      <Table.Cell>{entry.code}</Table.Cell>
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

export default GuestlistModal;

/*

                <Table hoverable>
                  <Table.Head>
                    <Table.HeadCell>Id</Table.HeadCell>
                    <Table.HeadCell>Name</Table.HeadCell>
                    <Table.HeadCell>&nbsp;</Table.HeadCell>
                    <Table.HeadCell>List</Table.HeadCell>
                    <Table.HeadCell>Code</Table.HeadCell>

                    <Table.HeadCell>List</Table.HeadCell>
                  </Table.Head>
                  <Table.Body className="divide-y">
                    <Table.Row>
                      <Table.Cell>1</Table.Cell>
                      <Table.Cell>Hans Hase</Table.Cell>
                      <Table.Cell>+1</Table.Cell>
                      <Table.Cell>GuestList Poti</Table.Cell>
                      <Table.Cell>123456</Table.Cell>
                    </Table.Row>
                  </Table.Body>
                </Table>

                <Avatar placeholderInitials="RR" size="md" rounded>
                  <div className="space-y-1 font-medium dark:text-white">
                    <div>Hans Hase</div>
                    <div className="text-sm text-gray-500 dark:text-gray-400">Joined in August 2014</div>
                  </div>
                </Avatar>


                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
                <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
                <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
                <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
                <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
                <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
                <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
                <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
                <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
             
            </div>

            */
