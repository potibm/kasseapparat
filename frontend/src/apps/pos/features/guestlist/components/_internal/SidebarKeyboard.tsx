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
          onClick={() => addToSearchTerm(letter)}
        >
          {letter}
        </SidebarKeyboardKey>
      ))}
      <SidebarKeyboardKey onClick={() => addToSearchTerm(" ")}>
        <HiOutlineMinusSm />
      </SidebarKeyboardKey>
      <SidebarKeyboardKey onClick={() => removeFromSearchTerm()}>
        <HiBackspace />
      </SidebarKeyboardKey>
      <SidebarKeyboardKey onClick={() => removeSearchTerm()}>
        <HiOutlineX />
      </SidebarKeyboardKey>
    </div>
  );
};

export default SidebarKeyboard;
