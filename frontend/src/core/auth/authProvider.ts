import { AuthProvider } from "react-admin";
import { z } from "zod";
import { createLogger } from "../logger/logger";

const logger = createLogger("Auth");

const AuthUserSchema = z.object({
  username: z.string(),
  role: z.string(),
});
type AuthUser = z.infer<typeof AuthUserSchema>;

export const createAuthProvider = (apiHost: string): AuthProvider => {
  let cachedUser: AuthUser | null = null;

  return {
    checkAuth: async (_params: unknown = {}) => {
      if (cachedUser) {
        return;
      }

      try {
        const response = await fetch(`${apiHost}/api/v2/auth/me`, {
          credentials: "include",
        });

        if (response.status === 401) {
          throw new Error("Not authenticated");
        }

        if (!response.ok) {
          throw new Error(`API Error: ${response.statusText}`);
        }

        const rawData = await response.json();

        cachedUser = AuthUserSchema.parse(rawData);

        return;
      } catch (error) {
        cachedUser = null;
        throw error;
      }
    },

    getPermissions: async (_params: unknown = {}) => {
      if (!cachedUser) {
        throw new Error("checkAuth must be called first");
      }
      return cachedUser?.role;
    },

    getIdentity: async () => {
      if (!cachedUser) throw new Error("Not authenticated");
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

      try {
        await fetch(`${apiHost}/api/v2/auth/logout`, {
          method: "POST",
          credentials: "include",
        });
      } catch (error) {
        logger.warn("Backend logout failed", error);
      }
    },
  };
};
