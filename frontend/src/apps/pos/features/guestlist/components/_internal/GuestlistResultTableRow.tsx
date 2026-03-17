import React, { Fragment } from "react";
import { TableCell, TableRow } from "flowbite-react";
import { HiShoppingCart } from "react-icons/hi";
import Button from "../../../../components/Button";
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
  const isAlreadyInCart = hasListItem(entry);

  const handleAddToCart = (additionalGuests: number) => {
    onAddToCart(entry, additionalGuests);
  };

  const highlightText = (text: string, highlight: string) => {
    if (!text) return "";
    if (!highlight.trim()) return text;

    const regex = new RegExp(`(${highlight})`, "gi");
    const parts = text.split(regex);

    return (
      <span>
        {parts.map((part, i) => (
          <Fragment key={`hl-${part}-${i}`}>
            {regex.test(part) ? (
              <span className="font-bold underline">{part}</span>
            ) : (
              part
            )}
          </Fragment>
        ))}
      </span>
    );
  };

  return (
    <TableRow
      key={entry.id}
      className="hover:bg-gray-100 dark:hover:bg-gray-700"
    >
      <TableCell>
        {!entry.code && <GuestlistAvatar name={entry.name} />}
      </TableCell>
      <TableCell>
        <div className="flex flex-col">
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
        </div>
      </TableCell>
      <TableCell className="flex gap-5">
        <Button
          className="float"
          key={`add-btn-${entry.id}-0`}
          disabled={isAlreadyInCart}
          onClick={() => handleAddToCart(0)}
        >
          <HiShoppingCart />
        </Button>
        {Array.from({ length: entry.additionalGuests }, (_, i) => {
          const count = i + 1;
          return (
            <Button
              key={`add-btn-${entry.id}-${count}`}
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
  );
};

export default GuestlistResultTableRow;
