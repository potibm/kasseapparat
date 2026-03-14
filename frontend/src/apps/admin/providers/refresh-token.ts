import { getSession, updateToken } from "../utils/auth-utils";

export const refreshToken = async (): Promise<void> => {
  // get current token and expiry from local storage
  const session = getSession();
  if (!session) {
    throw new Error("No active session found. Please log in.");
  }

  const now = new Date();
  const expiry = new Date(session.expiryDate);

  // Check if token is still valid for at least 10 seconds
  if (expiry.getTime() - now.getTime() > 10 * 1000) {
    return;
  }

  await updateToken();
};
