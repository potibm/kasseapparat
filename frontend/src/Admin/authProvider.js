import * as Sentry from "@sentry/react";
import { refreshToken } from "./refreshToken";
import { addRefreshAuthToAuthProvider } from "react-admin";
import {
  getAdminData,
  updateAdminData,
  clearAdminData,
  updateSession,
} from "./authUtils";

const API_HOST = import.meta.env.VITE_API_HOST ?? "https://localhost:3000";

class AuthorizationError extends Error {}

const authProvider = {
  login: async ({ username, password }) => {
    const request = new Request(`${API_HOST}/api/v2/auth/login`, {
      method: "POST",
      body: JSON.stringify({ login: username, password }),
      headers: new Headers({ "Content-Type": "application/json" }),
      credentials: "include",
    });

    try {
      const response = await fetch(request);
      if (response.status === 401) {
        const data = await response.json();
        throw new AuthorizationError(data.message);
      } else if (response.status < 200 || response.status >= 300) {
        throw new Error(response.statusText);
      }
      const { id, access_token, expires_in, role, username, gravatarUrl } =
        await response.json();

      Sentry.setUser({
        id: id?.toString(),
        username,
      });
      updateSession(access_token, expires_in);
      updateAdminData({ ID: id, username, role, gravatarUrl });
    } catch (error) {
      if (error instanceof AuthorizationError) {
        throw new Error(`There was an error logging you in: ${error.message}`);
      }

      Sentry.captureException(error, {
        tags: { auth: "login" },
        extra: { username },
      });

      throw new Error("Network error. Please try again.");
    }
  },
  logout: async () => {
    const request = new Request(`${API_HOST}/api/v2/auth/logout`, {
      method: "POST",
      headers: new Headers({
        "Content-Type": "application/json",
      }),
      credentials: "include",
    });

    try {
      const response = await fetch(request);
      if (response.status < 200 || response.status >= 300) {
        throw new Error(response.statusText);
      }
    } catch (error) {
      Sentry.captureException(error, {
        tags: { auth: "logout" },
      });

      console.error("Token logout error:", error);
    } finally {
      clearAdminData();
    }
  },
  checkError: ({ status }) => {
    if (status === 401 || status === 403) {
      clearAdminData();
      return Promise.reject(
        new Error("Authentication error. Please log in again."),
      );
    }
    return Promise.resolve();
  },
  checkAuth: () => {
    const { token, expiryDate } = getAdminData();

    if (!token || !expiryDate) {
      return Promise.reject(new Error("No token found. Please log in."));
    }

    return Promise.resolve();
  },
  getPermissions: () => {
    const adminData = getAdminData();
    const role = adminData ? adminData.role : null;
    return role
      ? Promise.resolve(role)
      : Promise.reject(new Error("No role found."));
  },
  getIdentity: () => {
    try {
      const adminData = getAdminData();
      const username = adminData ? adminData.username : null;
      const ID = adminData ? adminData.ID : null;
      const avatar = adminData ? adminData.gravatarUrl : null;
      return username
        ? Promise.resolve({ id: ID, fullName: username, avatar })
        : Promise.reject(new Error("No username found."));
    } catch (error) {
      return Promise.reject(error);
    }
  },
};
export default addRefreshAuthToAuthProvider(authProvider, refreshToken);
