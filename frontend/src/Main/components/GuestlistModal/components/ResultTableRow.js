import React, { useState } from "react";
import { Table, Avatar } from "flowbite-react";
import { HiShoppingCart } from "react-icons/hi";
import PropTypes from "prop-types";
import MyButton from "../../MyButton";
import GuestListArrivalNoteModal from "./ArrivalNoteModal";

const GuestListResultTableRow = ({
  entry,
  onAddToCart,
  hasListItem,
  loadedSearchQuery,
}) => {
  const showArrivalModal = (name, arrivalNote) => {
    return new Promise((resolve) => {
      const handleClose = () => {
        setArrivalModalContent(null);
        resolve();
      };

      setArrivalModalContent(
        <GuestListArrivalNoteModal
          isOpen={true}
          onClose={handleClose}
          arrivalNote={arrivalNote}
          name={name}
        />,
      );
    });
  };

  const [arrivalModalContent, setArrivalModalContent] = useState(null);

  const handleAddToCart = async (listEntry, additionalGuests) => {
    if (listEntry.arrivalNote) {
      await showArrivalModal(listEntry.name, listEntry.arrivalNote);
    }
    onAddToCart(listEntry, additionalGuests);
  };

  return (
    <>
      {arrivalModalContent}
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
            {...(hasListItem(entry.id) ? { disabled: true } : {})}
            onClick={() => handleAddToCart(entry, 0)}
          >
            <HiShoppingCart />
          </MyButton>
          {Array.from({ length: entry.additionalGuests }, (_, i) => (
            <MyButton
              key={i + 1}
              className="float"
              {...(hasListItem(entry.id) ? { disabled: true } : {})}
              onClick={() => handleAddToCart(entry, i + 1)}
            >
              <div className="text-xs">+{i + 1}</div>
            </MyButton>
          ))}
        </Table.Cell>
      </Table.Row>
    </>
  );
};

let uniqueIdCounter = 0; // Initialisierung eines Zählers

const generateUniqueId = () => {
  uniqueIdCounter += 1; // Erhöhe den Zähler
  return `key-${Date.now()}-${uniqueIdCounter}`; // Kombiniere Zeitstempel und Zähler
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
      {parts.map((part, i) => {
        const key = generateUniqueId();
        return regex.test(part) ? (
          <span key={key} className="font-bold underline">
            {part}
          </span>
        ) : (
          part
        );
      })}
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

GuestListResultTableRow.propTypes = {
  entry: PropTypes.object.isRequired,
  onAddToCart: PropTypes.func.isRequired,
  hasListItem: PropTypes.func.isRequired,
  loadedSearchQuery: PropTypes.string,
};

export default GuestListResultTableRow;
