import React from "react";
import { Button, Tooltip, ButtonProps } from "flowbite-react";
import { useAuth } from "../../../auth/providers/auth-provider";
import { HiShieldCheck } from "react-icons/hi";
import { getSession, initializeSession } from "@admin/utils/auth-utils";

const AdminButton: React.FC<ButtonProps> = ({ ...props }) => {
  const { getSafeToken, userdata } = useAuth();

  const handleAdminClick = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();

    if (!getSession()) {
      const token = await getSafeToken();

      if (!userdata?.id || !userdata?.username) {
        console.error("Cannot initialize admin session: User data missing.");
        return;
      }

      const data = {
        ID: userdata.id,
        username: userdata.username,
        role: userdata.role ?? "user",
        gravatarUrl: userdata.gravatarUrl ?? "",
      };

      await initializeSession(data, token, 5);
    }

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
