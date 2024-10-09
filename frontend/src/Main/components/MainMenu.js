import React from "react";
import { Button, DarkThemeToggle, Tooltip } from "flowbite-react";
import { HiCog, HiLogout, HiOutlineUserCircle } from "react-icons/hi";
import { Link } from "react-router-dom";
import PropTypes from "prop-types";

const MainMenu = ({ username, ...props }) => {
  return (
    <Button.Group className="mt-10">
      <Button size="sm">
        <Tooltip content={username}>
          <HiOutlineUserCircle className="h-5 w-5" />
        </Tooltip>
        <span className="ml-2 max-lg:hidden overflow-hidden max-w-10 text-nowrap text-sm">
          {username}
        </span>
      </Button>
      <Button
        as={Link}
        to="/logout"
        size="sm"
        className="hover:bg-cyan-800 dark:hover:bg-cyan-700 hover:text-white"
      >
        <Tooltip content="Logout">
          <HiLogout className="h-5 w-5" />
        </Tooltip>
        <span className="ml-2 max-xl:hidden text-sm">Logout</span>
      </Button>
      <Button
        as={Link}
        target="blank"
        to="/admin"
        size="sm"
        className="hover:bg-cyan-800 dark:hover:bg-cyan-700 hover:text-white"
      >
        <HiCog className="h-5 w-5" />
        <span className="ml-2 max-xl:hidden text-sm">Admin</span>
      </Button>
      <DarkThemeToggle className="hover:bg-cyan-800 dark:hover:bg-cyan-700 bg-cyan-700 dark:bg-cyan-600 text-white dark:text-white rounded-l-none text-sm  px-3 py-1.5" />
    </Button.Group>
  );
};

MainMenu.propTypes = {
  username: PropTypes.string.isRequired,
};

export default MainMenu;
