import { AuthProvider } from "react-admin";
import { z } from "zod";
import { createLogger } from "../logger/logger";

const logger = createLogger("Auth");

const AuthUserSchema = z.object({
  username: z.string(),
  role: z.string(),
});
type AuthUser = z.infer<typeof AuthUserSchema>;

export const createAuthProvider = (apiBaseUrl: string): AuthProvider => {
  let cachedUser: AuthUser | null = null;
  let pendingUser: Promise<AuthUser | null> | null = null;

  const redirectToLogin = () => {
    cachedUser = null;
    const returnTo = encodeURIComponent(
      window.location.pathname + window.location.search,
    );

    window.location.href = `${apiBaseUrl}/auth/login?returnTo=${returnTo}`;

    return new Promise<void>(() => {});
  };

  const requestUser = async (): Promise<AuthUser | null> => {
    const response = await fetch(`${apiBaseUrl}/auth/me`, {
      credentials: "include",
    });

    if (response.status === 401 || response.status === 403) {
      return null;
    }

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return AuthUserSchema.parse(await response.json());
  };

  /**
   * Resolves the current user, fetching `/auth/me` at most once.
   *
   * react-admin drives `checkAuth`, `getPermissions` and `getIdentity` from
   * three independent queries, so on a cold load all three start before any
   * fetch has resolved. They must therefore not depend on being called in a
   * particular order: awaiting one shared request keeps them consistent and
   * avoids fanning out into one request per caller. Callers that arrive while
   * the cache is cold join the in-flight request instead of starting another.
   */
  const ensureUser = (): Promise<AuthUser | null> => {
    if (cachedUser) {
      return Promise.resolve(cachedUser);
    }

    pendingUser ??= requestUser()
      .then((user) => {
        cachedUser = user;
        return user;
      })
      .finally(() => {
        pendingUser = null;
      });

    return pendingUser;
  };

  return {
    checkAuth: async (_params: unknown = {}) => {
      const user = await ensureUser();

      if (!user) {
        // Never settles: the browser is navigating to the login page.
        await redirectToLogin();
      }
    },

    getPermissions: async (_params: unknown = {}) => {
      const user = await ensureUser();

      if (!user) {
        await redirectToLogin();
      }

      return user?.role;
    },

    getIdentity: async () => {
      const user = await ensureUser();

      if (!user) {
        await redirectToLogin();
      }

      return {
        id: user?.username || "unknown",
        fullName: user?.username || "Unknown User",
      };
    },

    checkError: async (error) => {
      const status = error.status;
      if (status === 401 || status === 403) {
        return redirectToLogin();
      }
    },

    login: async () => {},

    logout: async () => {
      cachedUser = null;
      pendingUser = null;

      try {
        await fetch(`${apiBaseUrl}/auth/logout`, {
          method: "POST",
          credentials: "include",
        });
      } catch (error) {
        logger.warn("Backend logout failed", error);
      }
    },
  };
};
