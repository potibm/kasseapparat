import React from "react";
import { Button, Tooltip } from "flowbite-react";
import { useAuth } from "../../../Auth/provider/AuthProvider";
import PropTypes from "prop-types";
import { HiShieldCheck } from "react-icons/hi";

const AdminButton = ({ className }) => {
  const { token, expiryDate, userdata } = useAuth();

  const handleAdminClick = (e) => {
    e.preventDefault();

    const localstorageKey = "admin";

    let currentUserInAdmin = localStorage.getItem(localstorageKey);
    if (!currentUserInAdmin) {
      currentUserInAdmin = JSON.stringify({
        ID: userdata.id,
        token: token,
        username: userdata.username,
        role: userdata.role,
        expire: expiryDate,
        gravatarUrl: userdata.gravatarUrl,
      });
      localStorage.setItem(localstorageKey, currentUserInAdmin);
    }

    window.open("/admin", "_blank", "noopener,noreferrer");
  };

  return (
    <Button onClick={handleAdminClick} size="sm" className={className}>
      <Tooltip content="Admin">
        <HiShieldCheck className="h-5 w-5" />
      </Tooltip>
      <span className="ml-2 max-xl:hidden text-sm">Admin</span>
    </Button>
  );
};

AdminButton.propTypes = {
  className: PropTypes.string,
};

export default AdminButton;
