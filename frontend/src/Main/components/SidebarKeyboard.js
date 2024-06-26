import React from "react";
import { Button } from "flowbite-react";
import { HiBackspace, HiOutlineMinusSm, HiOutlineX } from "react-icons/hi";
import PropTypes from "prop-types";

const SidebarKeyboardKey = ({ children, onClick, ...props }) => {
  return (
    <Button
      onClick={onClick}
      size="md"
      className="min-w-12 max-w-12 mr-3.5 mb-3.5 me-0"
      {...props}
    >
      {children}
    </Button>
  );
};

SidebarKeyboardKey.propTypes = {
  children: PropTypes.node.isRequired,
  onClick: PropTypes.func.isRequired,
};

const SidebarKeyboard = ({ term, setTerm }) => {
  const addToSearchTerm = (letter) => {
    setTerm(term + letter);
  };

  const removeFromSearchTerm = () => {
    setTerm(term.slice(0, -1));
  };

  const removeSearchTerm = () => {
    setTerm("");
  };

  const alphabet = Array.from(Array(26)).map((e, i) =>
    String.fromCharCode(i + 65),
  );

  return (
    <div className="flex flex-wrap gap-1">
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

SidebarKeyboard.propTypes = {
  term: PropTypes.string.isRequired,
  setTerm: PropTypes.func.isRequired,
};

export default SidebarKeyboard;
