import { AuthProvider } from "react-admin";

let cachedUser: { username: string; role: string } | null = null;

export const authProvider: AuthProvider = {
  checkAuth: async (_params: unknown = {}) => {
    if (cachedUser) {
      return;
    }

    try {
      const response = await fetch("/api/v2/auth/me", {
        credentials: "omit",
      });

      if (!response.ok) {
        throw new Error("Not authenticated");
      }

      cachedUser = await response.json();

      return;
    } catch (error) {
      cachedUser = null;
      throw error;
    }
  },

  getPermissions: async (params: unknown = {}) => {
    if (!cachedUser) {
      await authProvider.checkAuth(params);
    }
    return cachedUser?.role;
  },

  getIdentity: async () => {
    if (!cachedUser) await authProvider.checkAuth({});
    return {
      id: cachedUser?.username || "unknown",
      fullName: cachedUser?.username || "Unknown User",
    };
  },

  checkError: async (error) => {
    const status = error.status;
    if (status === 401 || status === 403) {
      cachedUser = null;
      throw new Error("Unauthorized");
    }
  },

  login: async () => {},

  logout: async () => {
    cachedUser = null;
  },
};
