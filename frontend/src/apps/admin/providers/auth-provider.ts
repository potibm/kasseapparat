import * as Sentry from "@sentry/react";
import {
  AuthProvider,
  UserIdentity,
  addRefreshAuthToAuthProvider,
} from "react-admin";
import { refreshToken } from "./refresh-token";
import {
  getSession,
  clearSession,
  initializeSession,
} from "../utils/auth-utils";
import { getJwtToken, logout } from "../../pos/features/auth/hooks/api";
import { LoginResponse as LoginResponseType } from "../../pos/features/auth/hooks/api.schemas";

const API_HOST = import.meta.env.VITE_API_HOST ?? "https://localhost:3000";

const authProvider: AuthProvider = {
  login: async ({ username, password }) => {
    try {
      const userData: LoginResponseType = await getJwtToken(
        API_HOST,
        username,
        password,
      );

      Sentry.setUser({
        id: userData.id.toString(),
        username: userData.username,
      });

      initializeSession(
        {
          ID: userData.id,
          username: userData.username,
          role: userData.role,
          gravatarUrl: userData.gravatarUrl,
        },
        userData.access_token,
        userData.expires_in,
      );
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : "Unknown error";
      throw new Error(`There was an error logging you in: ${errorMessage}`, {
        cause: error,
      });
    }
  },

  logout: async () => {
    try {
      await logout(API_HOST);
    } catch (error) {
      Sentry.captureException(error, { tags: { auth: "logout" } });
    } finally {
      clearSession();
    }
  },

  checkError: async ({ status }) => {
    if (status === 401 || status === 403) {
      clearSession();
      throw new Error("Authentication error. Please log in again.");
    }
  },

  checkAuth: async () => {
    const session = getSession();
    if (!session) {
      throw new Error("No session found");
    }
  },

  getPermissions: async () => {
    const session = getSession();
    if (!session) {
      throw new Error("No session found");
    }

    return session.role;
  },

  getIdentity: async (): Promise<UserIdentity> => {
    const session = getSession();
    if (!session) {
      throw new Error("No session found");
    }

    return {
      id: session.ID,
      fullName: session.username,
      avatar: session.gravatarUrl,
    };
  },
};

export default addRefreshAuthToAuthProvider(authProvider, refreshToken);
