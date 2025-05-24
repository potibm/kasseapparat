import React from "react";
import { Button, ButtonGroup, DarkThemeToggle, Tooltip } from "flowbite-react";
import { HiCog, HiLogout, HiOutlineUserCircle } from "react-icons/hi";
import { Link } from "react-router";
import PropTypes from "prop-types";
import AdminButton from "./AdminButton";

const MainMenu = ({ username, ...props }) => {
  const buttonClass =
    "hover:bg-primary-800 dark:hover:bg-cyan-700 hover:text-white";

  return (
    <ButtonGroup className="mt-10" {...props}>
      <Button
        size="sm"
        aria-label={`Logged in as ${username}`}
        className="dark:hover:bg-cyan-700 border-none"
      >
        <Tooltip content={username}>
          <HiOutlineUserCircle className="h-5 w-5" />
        </Tooltip>
        <span className="ml-2 max-lg:hidden overflow-hidden max-w-10 text-nowrap text-sm">
          {username}
        </span>
      </Button>
      <Button as={Link} to="/logout" size="sm" className={buttonClass}>
        <Tooltip content="Logout">
          <HiLogout className="h-5 w-5" />
        </Tooltip>
        <span className="ml-2 max-xl:hidden text-sm">Logout</span>
      </Button>
      <AdminButton className={buttonClass} />
      <DarkThemeToggle
        aria-label="Toggle dark mode"
        className={`bg-primary-700 text-white dark:text-white rounded-l-none text-sm  px-3 py-1.5 ${buttonClass}`}
      />
    </ButtonGroup>
  );
};

MainMenu.propTypes = {
  username: PropTypes.string.isRequired,
};

export default MainMenu;
