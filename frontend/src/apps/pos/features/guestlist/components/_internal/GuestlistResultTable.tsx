import React from "react";
import { Alert, Spinner } from "flowbite-react";
import { HiInformationCircle, HiXCircle } from "react-icons/hi";
import GuestlistResultTableRow from "./GuestlistResultTableRow";
import { Guest as GuestType } from "@pos/api/schemas";

interface GuestlistResultTableProps {
  loading: boolean;
  error: string | null;
  guestlistEntries: GuestType[];
  onAddToCart: (listEntry: GuestType, additionalGuests: number) => void;
  hasListItem: (guest: GuestType) => boolean;
  loadedSearchQuery: string;
  /** The guest whose arrival note is shown after the check-in. */
  noteGuestId?: number | null;
  onAcknowledgeNote?: () => void;
}

const GuestlistResultTable: React.FC<GuestlistResultTableProps> = ({
  loading,
  error,
  guestlistEntries,
  onAddToCart,
  hasListItem,
  loadedSearchQuery,
  noteGuestId = null,
  onAcknowledgeNote,
}) => {
  return (
    <>
      {loading && (
        <div className="absolute inset-0 z-10 flex items-center justify-center bg-white/75 dark:bg-gray-900/75">
          <Spinner size="xl" />
        </div>
      )}

      {error && (
        <Alert className="my-3" color="failure" icon={HiInformationCircle}>
          {error}
        </Alert>
      )}

      {!loading && guestlistEntries.length === 0 && (
        <Alert className="my-3" color="warning" icon={HiXCircle}>
          No matching guests in this guestlist
        </Alert>
      )}

      {guestlistEntries.length > 0 && (
        <div
          className="flex flex-col gap-3"
          data-testid="guestlist-result-table"
        >
          <div className="text-sm text-gray-600 dark:text-gray-300">
            <b>{guestlistEntries.length}</b>{" "}
            {guestlistEntries.length === 1 ? "guest" : "guests"}
            {loadedSearchQuery.trim() !== "" && (
              <> matching „{loadedSearchQuery.trim()}“</>
            )}
          </div>
          {guestlistEntries.map((entry) => (
            <GuestlistResultTableRow
              key={entry.id}
              entry={entry}
              onAddToCart={onAddToCart}
              hasListItem={hasListItem}
              loadedSearchQuery={loadedSearchQuery}
              showArrivalNote={entry.id === noteGuestId}
              onAcknowledgeNote={onAcknowledgeNote}
            />
          ))}
        </div>
      )}
    </>
  );
};

export default GuestlistResultTable;
