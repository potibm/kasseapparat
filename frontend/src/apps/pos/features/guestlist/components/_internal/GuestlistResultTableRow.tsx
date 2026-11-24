import React from "react";
import { HiCheck, HiExclamation, HiExclamationCircle } from "react-icons/hi";
import { Guest as GuestType } from "@pos/api/schemas";
import { GuestlistAvatar } from "./GuestlistAvatar";
import HighlightedText from "./HighlightedText";

interface GuestlistResultTableRowProps {
  entry: GuestType;
  onAddToCart: (listEntry: GuestType, additionalGuests: number) => void;
  hasListItem: (guest: GuestType) => boolean;
  loadedSearchQuery: string;
  /** Set right after this guest was checked in and has an arrival note. */
  showArrivalNote?: boolean;
  onAcknowledgeNote?: () => void;
}

const actionClass =
  "pos-pressable h-14 rounded-xl font-bold disabled:cursor-not-allowed disabled:opacity-50";

/**
 * One guest in the search results: name (or code) with the matching part
 * highlighted, list name and allowed companions, and the check-in buttons. After a
 * check-in, the guest's arrival note opens inside the same card, so it is clearly
 * about this guest and not another search result.
 */
const GuestlistResultTableRow: React.FC<GuestlistResultTableRowProps> = ({
  entry,
  onAddToCart,
  hasListItem,
  loadedSearchQuery,
  showArrivalNote = false,
  onAcknowledgeNote,
}) => {
  const isAlreadyInCart = hasListItem(entry);
  const displayName = entry.code ?? entry.name;

  const details = [entry.listName];
  if (entry.additionalGuests > 0) {
    details.push(
      `up to ${entry.additionalGuests} additional guest${
        entry.additionalGuests === 1 ? "" : "s"
      }`,
    );
  }

  let cardClass = "border-gray-200 dark:border-gray-700";
  if (showArrivalNote) {
    cardClass = "border-2 border-orange-400 shadow-lg";
  } else if (isAlreadyInCart) {
    cardClass =
      "border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-950";
  }

  return (
    <div
      data-testid={"guestlist-result-id-" + entry.id}
      className={`overflow-hidden rounded-xl border ${cardClass}`}
    >
      <div className="flex items-center gap-4 p-4">
        {!entry.code && <GuestlistAvatar name={entry.name} />}
        <div className="min-w-0 flex-1">
          <div
            className={
              entry.code ? "font-mono text-3xl" : "truncate text-xl font-medium"
            }
          >
            <HighlightedText text={displayName} query={loadedSearchQuery} />
          </div>
          <div className="flex flex-wrap items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
            <span>{details.join(" · ")}</span>
            {entry.arrivalNote && (
              <span className="inline-flex items-center gap-1 rounded-full bg-orange-100 px-2 py-0.5 font-semibold text-orange-800 dark:bg-orange-900 dark:text-orange-100">
                <HiExclamationCircle aria-hidden="true" />
                Arrival note
              </span>
            )}
          </div>
        </div>

        {isAlreadyInCart ? (
          <span className="inline-flex h-11 items-center gap-2 rounded-full bg-green-100 px-4 font-semibold text-green-800 dark:bg-green-900 dark:text-green-100">
            <HiCheck aria-hidden="true" />
            {showArrivalNote ? "Added" : "In cart"}
          </span>
        ) : (
          <div className="flex shrink-0 flex-wrap justify-end gap-2">
            <button
              type="button"
              className={`${actionClass} bg-cyan-700 px-5 text-white`}
              onClick={() => onAddToCart(entry, 0)}
              aria-label={"Add " + displayName + " to cart"}
            >
              Check in
            </button>
            {Array.from({ length: entry.additionalGuests }, (_, i) => {
              const count = i + 1;
              return (
                <button
                  key={count}
                  type="button"
                  className={`${actionClass} w-16 border-2 border-cyan-700 bg-white text-lg text-cyan-700 dark:bg-gray-800 dark:text-cyan-300`}
                  onClick={() => onAddToCart(entry, count)}
                  aria-label={
                    "Add " +
                    displayName +
                    " with " +
                    count +
                    " additional guest(s) to cart"
                  }
                >
                  +{count}
                </button>
              );
            })}
          </div>
        )}
      </div>

      {showArrivalNote && entry.arrivalNote && (
        <div
          role="alert"
          data-testid="guestlist-arrival-note"
          className="flex items-start gap-3 border-t border-orange-200 bg-orange-50 p-4 dark:border-orange-800 dark:bg-orange-950"
        >
          <HiExclamation
            aria-hidden="true"
            className="h-7 w-7 shrink-0 text-orange-700 dark:text-orange-300"
          />
          <div className="flex-1">
            <div className="text-sm font-bold tracking-wide text-orange-800 uppercase dark:text-orange-200">
              Note for this guest
            </div>
            <div className="text-lg font-medium text-orange-950 dark:text-orange-50">
              {entry.arrivalNote}
            </div>
          </div>
          <button
            type="button"
            className="pos-pressable h-14 shrink-0 rounded-xl bg-orange-700 px-6 font-bold text-white"
            onClick={onAcknowledgeNote}
          >
            Got it
          </button>
        </div>
      )}
    </div>
  );
};

export default GuestlistResultTableRow;
