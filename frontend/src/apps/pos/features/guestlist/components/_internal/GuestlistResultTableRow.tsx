import React, { useState, Fragment } from "react";
import { TableCell, TableRow } from "flowbite-react";
import { HiShoppingCart } from "react-icons/hi";
import Button from "../../../../components/Button";
import GuestlistArrivalNoteModal from "./GuestlistArrivalNoteModal";
import { Guest as GuestType } from "@pos/utils/api.schemas";
import { GuestlistAvatar } from "./GuestlistAvatar";

interface GuestlistResultTableRowProps {
  entry: GuestType;
  onAddToCart: (listEntry: GuestType, additionalGuests: number) => void;
  hasListItem: (guest: GuestType) => boolean;
  loadedSearchQuery: string;
}

const GuestlistResultTableRow: React.FC<GuestlistResultTableRowProps> = ({
  entry,
  onAddToCart,
  hasListItem,
  loadedSearchQuery,
}) => {
  const [showNote, setShowNote] = useState<boolean>(false);
  const isAlreadyInCart = hasListItem(entry);

  const handleAddToCart = (additionalGuests: number) => {
    onAddToCart(entry, additionalGuests);

    if (entry.arrivalNote) {
      setShowNote(true);
    }
  };

  const highlightText = (text: string, highlight: string) => {
    if (!text) return "";
    if (!highlight.trim()) return text;

    const regex = new RegExp(`(${highlight})`, "gi");
    const parts = text.split(regex);

    return parts.map((part, i) => (
      <Fragment key={`part-${i}`}>
        {regex.test(part) ? (
          <span className="font-bold underline">{part}</span>
        ) : (
          part
        )}
      </Fragment>
    ));
  };

  return (
    <>
      {entry.arrivalNote && (
        <GuestlistArrivalNoteModal
          isOpen={showNote}
          onClose={() => setShowNote(false)}
          arrivalNote={entry.arrivalNote}
          name={entry.name}
        />
      )}
      <TableRow key={entry.id}>
        <TableCell>
          {!entry.code && <GuestlistAvatar name={entry.name} />}
        </TableCell>
        <TableCell className="">
          {!entry.code && (
            <>
              <div className="text-xl">
                {highlightText(entry.name, loadedSearchQuery)}
              </div>
              <div className="text-sm">{entry.listName}</div>
            </>
          )}
          {entry.code !== null && entry.code !== "" && (
            <div className="text-3xl font-mono">
              {highlightText(String(entry.code), loadedSearchQuery)}
            </div>
          )}
        </TableCell>
        <TableCell className="flex gap-5">
          <Button
            className="float"
            key={0}
            disabled={isAlreadyInCart}
            onClick={() => handleAddToCart(0)}
          >
            <HiShoppingCart />
          </Button>
          {Array.from({ length: entry.additionalGuests }, (_, i) => {
            const count = i + 1;
            return (
              <Button
                key={`add-${entry.id}-${count}`}
                className="float"
                disabled={isAlreadyInCart}
                onClick={() => handleAddToCart(count)}
              >
                <div className="text-xs">+{count}</div>
              </Button>
            );
          })}
        </TableCell>
      </TableRow>
    </>
  );
};

export default GuestlistResultTableRow;
