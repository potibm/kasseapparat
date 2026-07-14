import { AuthProvider } from "react-admin";

// 1. Der In-Memory Cache
// Wir speichern die Rolle hier in einer simplen Variable.
// So verhindern wir, dass getPermissions() bei jedem Rendern eines Buttons das Backend anpingt.
let cachedUser: { username: string; role: string } | null = null;

export const authProvider: AuthProvider = {
  // ----------------------------------------------------------------------
  // Die Kern-Methoden für beide Frontends
  // ----------------------------------------------------------------------

  checkAuth: async (_params: unknown = {}) => {
    // Wenn wir die Rolle schon kennen, winken wir den Request sofort durch.
    if (cachedUser) {
      return;
    }

    // Ansonsten fragen wir unser zustandsloses Backend.
    try {
      const response = await fetch("/api/v2/auth/me", {
        // WICHTIG: Damit Traefik-Header oder spätere Cookies mitgesendet werden!
        credentials: "omit", // Bei Headern via Proxy oft egal, bei Cookies (OIDC) später 'include'
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
    // Falls checkAuth noch nicht lief, garantieren wir hier, dass es nachgeholt wird.
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

  // ----------------------------------------------------------------------
  // Hilfs-Methoden (Primär für React Admin)
  // ----------------------------------------------------------------------

  checkError: async (error) => {
    // React Admin ruft das auf, wenn ein API-Call fehlschlägt.
    // Wenn das Backend 401/403 wirft, löschen wir den Cache und loggen den User aus.
    const status = error.status;
    if (status === 401 || status === 403) {
      cachedUser = null;
      throw new Error("Unauthorized");
    }
  },

  login: async () => {
    // Vorerst ein Platzhalter.
    // In Phase 2/3 (OIDC) machen wir hier ein: window.location.href = '/api/auth/login'
  },

  logout: async () => {
    cachedUser = null;
  },
};
