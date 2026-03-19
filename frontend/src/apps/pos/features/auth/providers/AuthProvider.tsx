/* eslint-disable react-refresh/only-export-components */
import React, {
  createContext,
  useContext,
  useMemo,
  useRef,
  useState,
  ReactNode,
} from "react";
import { refreshJwtToken } from "../../../../../core/api/auth";
import { useConfig } from "../../../../../core/config/providers/ConfigProvider";
import {
  AuthContextType,
  AuthUser as AuthUserType,
  Session as SessionType,
} from "../types/auth.types";
import {
  getInitialSession,
  getInitialUser,
  clearAuthStorage,
  storeSession,
  storeUser,
} from "../services/auth-storage";
import { createLogger } from "@core/logger/logger";

const log = createLogger('Auth');

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const { apiHost } = useConfig();
  const [session, setSession] = useState<SessionType>(getInitialSession);
  const [user, setUser] = useState<AuthUserType | null>(getInitialUser);
  const refreshingPromise = useRef<Promise<string | null> | null>(null);

  const updateSession = React.useCallback(
    (token: string, expiresIn: number) => {
      const expiryDate = new Date(Date.now() + (expiresIn - 30) * 1000);
      log.debug("Updating session with new expiry date",  expiryDate.toISOString());

      setSession({ token: token, expiryDate: expiryDate });
      storeSession(token, expiryDate);

      return token;
    },
    [],
  );

  const removeSession = React.useCallback(() => {
    setSession({ token: null, expiryDate: null });
    setUser(null);
    clearAuthStorage();
    log.debug("User logged out, session cleared");
  }, []);

  const updateUser = React.useCallback((userdata: AuthUserType) => {
    setUser({ ...userdata });
    storeUser(userdata);
  }, []);

  const getToken = React.useCallback(async (): Promise<string | null> => {
    const now = new Date();
    const isTokenValid =
      session.token &&
      session.expiryDate &&
      session.expiryDate.getTime() - now.getTime() > 5000;

    if (isTokenValid) return session.token;
    if (!session.token) return null;

    if (refreshingPromise.current) return refreshingPromise.current;

    log.debug("Token expired or missing, starting refresh...");

    refreshingPromise.current = refreshJwtToken(apiHost)
      .then((res) => updateSession(res.access_token, res.expires_in))
      .catch(() => {
        removeSession();
        return null;
      })
      .finally(() => {
        refreshingPromise.current = null;
      });

    return refreshingPromise.current;
  }, [session, apiHost, updateSession, removeSession]);

  const getSafeToken = React.useCallback(async () => {
    const token = await getToken();
    if (!token) throw new Error("Authentication required");
    return token;
  }, [getToken]);

  const contextValue = useMemo<AuthContextType>(
    () => ({
      getToken,
      getSafeToken,
      isLoggedIn: async () => !!(await getToken()),
      setSession: updateSession,
      removeSession,
      userdata: user,
      setUserdata: updateUser,
      gravatarUrl: user?.gravatarUrl ?? "",
      role: user?.role ?? "user",
      username: user?.username ?? "unknown",
      id: user?.id ?? 0,
    }),
    [getToken, getSafeToken, updateSession, updateUser, removeSession, user],
  );

  // Provide the authentication context to the children components
  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export default AuthProvider;
