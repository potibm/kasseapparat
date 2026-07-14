import { ReactNode, useState, useEffect, useMemo } from "react";
import { authProvider } from "@core/auth/authProvider";
import { AuthContext } from "../context/AuthContext";

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
        // Wir starten den Motor: Prüfen, ob eine Session/Header da ist
        await authProvider.checkAuth({});

        // Wenn das geklappt hat, holen wir uns die Rolle
        const currentRole = await authProvider.getPermissions?.({});
        const identity = await authProvider.getIdentity?.();

        setIsAuthenticated(true);
        setRole(currentRole as string);
        setUsername(identity?.fullName as string);
      } catch {
        // 3. Variable "(error)" weggelassen, da sie nicht genutzt wird
        setIsAuthenticated(false);
        setRole(null);
        setUsername(null);
      } finally {
        // Egal ob Erfolg oder Fehler: Wir sind mit dem Laden fertig
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  const contextValue = useMemo(
    () => ({ isLoading, isAuthenticated, role, username }), // NEU
    [isLoading, isAuthenticated, role, username], // NEU
  );

  // Solange das Backend noch nicht geantwortet hat, zeigen wir nichts (oder einen Spinner)
  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        Loading Kasseapparat...
      </div>
    );
  }

  return <AuthContext value={contextValue}>{children}</AuthContext>;
};
