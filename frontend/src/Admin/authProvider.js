import { jwtDecode } from "jwt-decode";
import * as Sentry from "@sentry/react";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3001";
const ADMIN_STORAGE_KEY = "admin";

let updateTokenIntervalId = null;

class AuthorizationError extends Error {}

const setAdminData = (data) => {
  localStorage.setItem(ADMIN_STORAGE_KEY, JSON.stringify(data));
};

const updateAdminData = (data) => {
  const adminData = getAdminData();
  setAdminData({ ...adminData, ...data });
};

const getAdminData = () => {
  const adminData = localStorage.getItem(ADMIN_STORAGE_KEY);
  return adminData ? JSON.parse(adminData) : null;
};

const clearAdminData = () => {
  localStorage.removeItem(ADMIN_STORAGE_KEY);
};

const authProvider = {
  login: async ({ username, password }) => {
    const request = new Request(`${API_HOST}/login`, {
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
      const { id, token, role, username, gravatarUrl } = await response.json();
      const decodedToken = jwtDecode(token);
      const expire = new Date(decodedToken.exp * 1000);

      Sentry.setUser({
        id: id?.toString(),
        username,
      });

      setAdminData({ ID: id, token, username, role, expire, gravatarUrl });
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
  updateToken: async () => {
    const adminData = getAdminData();
    if (!adminData?.token) {
      return Promise.reject(new Error("No token found. Please log in."));
    }

    const request = new Request(`${API_HOST}/auth/refresh_token`, {
      method: "GET",
      headers: new Headers({
        "Content-Type": "application/json",
        Authorization: `Bearer ${adminData.token}`,
      }),
      credentials: "include",
    });

    try {
      const response = await fetch(request);
      if (response.status < 200 || response.status >= 300) {
        throw new Error(response.statusText);
      }

      const { token } = await response.json();

      const decodedToken = jwtDecode(token);
      const expire = new Date(decodedToken.exp * 1000);

      updateAdminData({ token, expire });
    } catch (error) {
      Sentry.captureException(error, {
        tags: { auth: "refresh_token" },
        extra: { token: adminData?.token },
      });

      console.error("Token refresh error:", error);
      clearAdminData();
      throw new Error("Token refresh error. Please log in again.");
    }
  },
  logout: () => {
    if (updateTokenIntervalId !== null) {
      clearInterval(updateTokenIntervalId);
      updateTokenIntervalId = null;
    }
    Sentry.setUser(null);
    clearAdminData();
    return Promise.resolve();
  },
  checkError: ({ status }) => {
    if (status === 401 || status === 403) {
      if (updateTokenIntervalId !== null) {
        clearInterval(updateTokenIntervalId);
        updateTokenIntervalId = null;
      }
      clearAdminData();
      return Promise.reject(
        new Error("Authentication error. Please log in again."),
      );
    }
    return Promise.resolve();
  },
  checkAuth: () => {
    const adminData = getAdminData();
    if (!adminData?.token) {
      return Promise.reject(new Error("No token found. Please log in."));
    }
    // check if the token is expired
    if (new Date(adminData.expire) <= new Date()) {
      return Promise.reject(new Error("Token has expired. Please log in."));
    }

    if (!updateTokenIntervalId) {
      updateTokenIntervalId = setInterval(
        () => {
          authProvider.updateToken().catch((error) => {
            clearInterval(updateTokenIntervalId);
            console.error("Token update error:", error);
          });
        },
        2 * 60 * 1000,
      );
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
export default authProvider;
export { getAdminData };
