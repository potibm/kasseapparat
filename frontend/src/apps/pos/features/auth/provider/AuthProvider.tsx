import { ReactNode, useState, useEffect, useMemo } from "react";
import { createAuthProvider } from "@core/auth/authProvider";
import { AuthContext } from "../context/AuthContext";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3100";
const authProvider = createAuthProvider(API_HOST);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [role, setRole] = useState<string | null>(null);
  const [username, setUsername] = useState<string | null>(null);

  useEffect(() => {
    const initAuth = async () => {
      try {
        await authProvider.checkAuth({});

        const currentRole = await authProvider.getPermissions?.({});
        const identity = await authProvider.getIdentity?.();

        setIsAuthenticated(true);
        setRole(currentRole as string);
        setUsername(identity?.fullName as string);
      } catch {
        setIsAuthenticated(false);
        setRole(null);
        setUsername(null);
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  const contextValue = useMemo(
    () => ({ isLoading, isAuthenticated, role, username }), // NEU
    [isLoading, isAuthenticated, role, username], // NEU
  );

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        ⏳ Loading Kasseapparat...
      </div>
    );
  }

  return <AuthContext value={contextValue}>{children}</AuthContext>;
};
