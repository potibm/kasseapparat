import React from "react";
import Logo from "@core/logo/Logo";
import UserMenu from "./UserMenu";

interface PosHeaderProps {
  username: string;
  environmentMessage?: string;
}

/**
 * Header bar above the product grid: app name, the environment message (e.g.
 * "Test environment") and the user menu. It frees the sidebar from the menu
 * buttons so the cart gets more room.
 */
export const PosHeader: React.FC<PosHeaderProps> = ({
  username,
  environmentMessage,
}) => (
  <header
    className="mb-3 flex min-h-14 items-center gap-3 rounded-lg bg-white px-3 py-2 shadow-sm dark:bg-gray-900"
    data-testid="pos-header"
  >
    <div className="flex shrink-0 items-center gap-2 text-lg font-bold text-gray-900 dark:text-white">
      <Logo className="h-6 w-6 text-[#ff3873]" aria-hidden="true" />
      Kasseapparat
    </div>
    {environmentMessage && (
      <span
        role="status"
        className="min-w-0 truncate rounded-full bg-blue-100 px-3 py-1 text-sm font-semibold text-blue-800 dark:bg-blue-900 dark:text-blue-100"
      >
        {environmentMessage}
      </span>
    )}
    <div className="flex-1" />
    <UserMenu username={username} />
  </header>
);

export default PosHeader;
