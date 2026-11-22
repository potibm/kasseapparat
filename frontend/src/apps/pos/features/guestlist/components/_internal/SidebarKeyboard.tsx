import React from "react";
import { HiBackspace } from "react-icons/hi";

interface SidebarKeyboardProps {
  term: string;
  setTerm: (term: string) => void;
}

const LETTERS = Array.from({ length: 26 }, (_e, i) =>
  String.fromCodePoint(i + 65),
);

const keyClass =
  "pos-pressable flex h-14 items-center justify-center rounded-lg border border-gray-300 bg-white text-xl font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-white";

/**
 * On-screen letters for the guest search: A–Z only. The search ignores diacritics,
 * so "MUL" also finds "Müller", and the iPad keyboard never has to open.
 */
const SidebarKeyboard: React.FC<SidebarKeyboardProps> = ({ term, setTerm }) => (
  <div className="grid grid-cols-6 gap-2" data-testid="guestlist-keyboard">
    {LETTERS.map((letter) => (
      <button
        key={letter}
        type="button"
        className={keyClass}
        aria-label={"Add " + letter + " to search term"}
        onClick={() => setTerm(term + letter)}
      >
        {letter}
      </button>
    ))}
    <button
      type="button"
      className={`${keyClass} col-span-2 text-base font-normal text-gray-600 dark:text-gray-300`}
      aria-label="Add space to search term"
      onClick={() => setTerm(term + " ")}
    >
      space
    </button>
    <button
      type="button"
      className={`${keyClass} col-span-2 bg-gray-200 dark:bg-gray-700`}
      aria-label="Remove last character from search term"
      onClick={() => setTerm(term.slice(0, -1))}
    >
      <HiBackspace className="h-6 w-6" aria-hidden="true" />
    </button>
  </div>
);

export default SidebarKeyboard;
