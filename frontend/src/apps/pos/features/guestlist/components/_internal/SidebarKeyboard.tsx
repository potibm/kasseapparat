import React from "react";
import { HiBackspace, HiOutlineMinusSm, HiOutlineX } from "react-icons/hi";
import SidebarKeyboardKey from "./SidebarKeyboardKey";

interface SidebarKeyboardProps {
  term: string;
  setTerm: (term: string) => void;
}

const SidebarKeyboard: React.FC<SidebarKeyboardProps> = ({ term, setTerm }) => {
  const addToSearchTerm = (letter: string) => {
    setTerm(term + letter);
  };

  const removeFromSearchTerm = () => {
    setTerm(term.slice(0, -1));
  };

  const removeSearchTerm = () => {
    setTerm("");
  };

  const alphabet = Array.from(new Array(26)).map((_e, i) =>
    String.fromCodePoint(i + 65),
  );

  return (
    <div className="flex flex-wrap gap-5">
      {alphabet.map((letter) => (
        <SidebarKeyboardKey
          key={letter}
          aria-label={"Add " + letter + " to search term"}
          onClick={() => addToSearchTerm(letter)}
        >
          {letter}
        </SidebarKeyboardKey>
      ))}
      <SidebarKeyboardKey
        onClick={() => addToSearchTerm(" ")}
        aria-label="Add space to search term"
      >
        <HiOutlineMinusSm />
      </SidebarKeyboardKey>
      <SidebarKeyboardKey
        onClick={() => removeFromSearchTerm()}
        aria-label="Remove last character from search term"
      >
        <HiBackspace />
      </SidebarKeyboardKey>
      <SidebarKeyboardKey
        onClick={() => removeSearchTerm()}
        aria-label="Clear search term"
      >
        <HiOutlineX />
      </SidebarKeyboardKey>
    </div>
  );
};

export default SidebarKeyboard;
