import * as Sentry from "@sentry/react";

const API_HOST = import.meta.env.VITE_API_HOST ?? "https://localhost:3000";
const LOCALSTORAGE_KEY = "kasseapparat.admin.auth";

export const setAdminData = (data) => {
  localStorage.setItem(LOCALSTORAGE_KEY, JSON.stringify(data));
};

export const updateAdminData = (data) => {
  const adminData = getAdminData();
  setAdminData({ ...adminData, ...data });
};

export const getAdminData = () => {
  const adminData = localStorage.getItem(LOCALSTORAGE_KEY);
  return adminData ? JSON.parse(adminData) : null;
};

export const clearAdminData = () => {
  localStorage.removeItem(LOCALSTORAGE_KEY);
};

export const updateSession = (token, expiresIn) => {
  const expiryDate = new Date(
    Date.now() + (expiresIn - 30) * 1000,
  ).toISOString();

  updateAdminData({ token, expiryDate });

  return token;
};

let refreshingPromise = null;

export const updateToken = async () => {
  if (refreshingPromise) {
    return refreshingPromise;
  }

  refreshingPromise = (async () => {
    const request = new Request(`${API_HOST}/api/v2/auth/refresh`, {
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

      const { access_token, expires_in } = await response.json();

      return updateSession(access_token, expires_in);
    } catch (error) {
      Sentry.captureException(error, {
        tags: { auth: "refresh_token" },
      });

      console.error("Token refresh error:", error);
      clearAdminData();
      throw new Error("Token refresh error. Please log in again.", {
        cause: error,
      });
    } finally {
      refreshingPromise = null;
    }
  })();

  return refreshingPromise;
};
