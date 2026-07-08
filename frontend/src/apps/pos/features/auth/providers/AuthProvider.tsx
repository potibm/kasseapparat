import React, { useMemo, ReactNode } from "react";
import { AuthContextType, AuthUser as AuthUserType } from "../types/auth.types";
import { AuthContext } from "../context/AuthContext";

// Phase 1: Dummy auth bypass - always authenticated
const DUMMY_USER: AuthUserType = {
  id: 1,
  username: "dummy",
  email: "dummy@example.com",
  role: "admin",
  gravatarUrl: "",
};

const DUMMY_TOKEN = "dummy-token-for-phase-1";

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const getToken = React.useCallback(async (): Promise<string | null> => {
    return DUMMY_TOKEN;
  }, []);

  const getSafeToken = React.useCallback(async () => {
    return DUMMY_TOKEN;
  }, []);

  const isLoggedIn = React.useCallback(async () => true, []);

  const contextValue = useMemo<AuthContextType>(
    () => ({
      getToken,
      getSafeToken,
      isLoggedIn,
      setSession: () => DUMMY_TOKEN,
      removeSession: () => {},
      userdata: DUMMY_USER,
      setUserdata: () => {},
      gravatarUrl: DUMMY_USER.gravatarUrl,
      role: DUMMY_USER.role,
      username: DUMMY_USER.username,
      id: DUMMY_USER.id,
    }),
    [getToken, getSafeToken, isLoggedIn],
  );

  return <AuthContext value={contextValue}>{children}</AuthContext>;
};

export default AuthProvider;
