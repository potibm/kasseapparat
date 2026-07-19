import { Tooltip, Button, ButtonProps } from "flowbite-react";
import { HiLogout } from "react-icons/hi";
import useConfig from "@core/config/hooks/useConfig";
import { createAuthProvider } from "@core/auth/authProvider";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3100";
const authProvider = createAuthProvider(API_HOST);

export const LogoutButton: React.FC<ButtonProps> = ({ ...props }) => {
  const { authMode } = useConfig();

  if (authMode !== "oidc") {
    return null;
  }

  const handleLogout = async () => {
    try {
      await authProvider.logout({});
    } catch {
      // noop
    } finally {
      window.location.href = "/";
    }
  };

  return (
    <Button onClick={handleLogout} size="sm" {...props}>
      <Tooltip content="Logout">
        <HiLogout className="h-5 w-5" />
      </Tooltip>
      <span className="ml-2 max-xl:hidden text-sm">Logout</span>
    </Button>
  );
};
