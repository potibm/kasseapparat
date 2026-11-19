import React, { useEffect, useRef, useState } from "react";
import { useThemeMode } from "flowbite-react";
import {
  HiBookOpen,
  HiChevronDown,
  HiLogout,
  HiMoon,
  HiOutlineUserCircle,
  HiShieldCheck,
  HiSun,
} from "react-icons/hi";
import useConfig from "@core/config/hooks/useConfig";
import { createAuthProvider } from "@core/auth/authProvider";
import { buildApiBaseUrl } from "@core/config/constants";
import Version from "@pos/components/Version";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3100";
const authProvider = createAuthProvider(buildApiBaseUrl(API_HOST));

const itemClass =
  "flex min-h-11 w-full items-center gap-3 rounded-lg px-3 text-left text-sm text-gray-800 hover:bg-gray-100 dark:text-gray-100 dark:hover:bg-gray-700";

interface UserMenuProps {
  username: string;
}

/**
 * The user menu in the POS header. It gathers everything that is rarely needed
 * while selling – admin, dark mode, manual, version and logout – behind one button,
 * so none of it can be hit by accident. Logout asks for confirmation.
 */
export const UserMenu: React.FC<UserMenuProps> = ({ username }) => {
  const { authMode } = useConfig();
  const { computedMode, toggleMode } = useThemeMode();
  const [open, setOpen] = useState(false);
  const [confirmLogout, setConfirmLogout] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);

  const close = () => {
    setOpen(false);
    setConfirmLogout(false);
  };

  useEffect(() => {
    if (!open) return;

    const handlePointerDown = (event: PointerEvent) => {
      if (rootRef.current && !rootRef.current.contains(event.target as Node)) {
        setOpen(false);
        setConfirmLogout(false);
      }
    };
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
        setConfirmLogout(false);
      }
    };

    document.addEventListener("pointerdown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [open]);

  const handleLogout = async () => {
    try {
      await authProvider.logout({});
    } catch {
      // noop: we leave the POS either way
    } finally {
      globalThis.location.href = "/";
    }
  };

  const isDark = computedMode === "dark";

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        aria-haspopup="menu"
        aria-expanded={open}
        aria-label={`User menu for ${username}`}
        onClick={() => (open ? close() : setOpen(true))}
        className="flex h-11 items-center gap-2 rounded-lg border border-gray-300 bg-white px-3 text-sm font-medium text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
      >
        <HiOutlineUserCircle className="h-5 w-5" aria-hidden="true" />
        <span className="max-w-40 truncate">{username}</span>
        <HiChevronDown className="h-4 w-4" aria-hidden="true" />
      </button>

      {open && (
        <div
          role="menu"
          aria-label="User menu"
          className="absolute right-0 top-full z-30 mt-2 w-64 rounded-xl border border-gray-200 bg-white p-2 shadow-lg dark:border-gray-700 dark:bg-gray-800"
        >
          <a
            role="menuitem"
            href="/admin"
            target="_blank"
            rel="noopener noreferrer"
            className={itemClass}
            onClick={close}
          >
            <HiShieldCheck className="h-5 w-5" aria-hidden="true" />
            Admin
          </a>
          <a
            role="menuitem"
            href="/manual.pdf"
            target="_blank"
            rel="noopener noreferrer"
            className={itemClass}
            onClick={close}
          >
            <HiBookOpen className="h-5 w-5" aria-hidden="true" />
            Manual
          </a>
          <button
            type="button"
            role="menuitem"
            className={itemClass}
            onClick={() => toggleMode()}
          >
            {isDark ? (
              <HiSun className="h-5 w-5" aria-hidden="true" />
            ) : (
              <HiMoon className="h-5 w-5" aria-hidden="true" />
            )}
            {isDark ? "Light mode" : "Dark mode"}
          </button>

          <div className="px-3 py-2 text-xs text-gray-500 dark:text-gray-400">
            <Version />
          </div>

          {authMode === "oidc" && (
            <div className="mt-1 border-t border-gray-200 pt-2 dark:border-gray-700">
              {confirmLogout ? (
                <div className="px-1">
                  <p className="px-2 pb-2 text-sm text-gray-800 dark:text-gray-100">
                    Really log out?
                  </p>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      className="min-h-11 flex-1 rounded-lg border border-gray-300 text-sm font-medium text-gray-800 dark:border-gray-600 dark:text-gray-100"
                      onClick={() => setConfirmLogout(false)}
                    >
                      Cancel
                    </button>
                    <button
                      type="button"
                      className="min-h-11 flex-1 rounded-lg bg-red-700 text-sm font-semibold text-white"
                      onClick={handleLogout}
                    >
                      Log out
                    </button>
                  </div>
                </div>
              ) : (
                <button
                  type="button"
                  role="menuitem"
                  className={`${itemClass} text-red-700 dark:text-red-400`}
                  onClick={() => setConfirmLogout(true)}
                >
                  <HiLogout className="h-5 w-5" aria-hidden="true" />
                  Logout
                </button>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default UserMenu;
