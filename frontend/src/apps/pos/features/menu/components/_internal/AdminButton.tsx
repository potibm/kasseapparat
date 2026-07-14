import React from "react";
import { Button, Tooltip, ButtonProps } from "flowbite-react";
import { HiShieldCheck } from "react-icons/hi";

const AdminButton: React.FC<ButtonProps> = ({ ...props }) => {
  const handleAdminClick = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();

    window.open("/admin", "_blank", "noopener,noreferrer");
  };

  return (
    <Button onClick={handleAdminClick} size="sm" {...props}>
      <Tooltip content="Admin">
        <HiShieldCheck className="h-5 w-5" />
      </Tooltip>
      <span className="ml-2 max-xl:hidden text-sm">Admin</span>
    </Button>
  );
};

export default AdminButton;
