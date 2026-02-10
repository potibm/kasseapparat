import { getAdminData, updateToken } from "./authUtils";

export const refreshToken = () => {
  // get current token and expiry from local storage
  const adminData = getAdminData();
  const token = adminData?.token;
  const expiryDate = adminData?.expiryDate;

  if (!token || !expiryDate) {
    return Promise.reject(new Error("No token found. Please log in."));
  }

  const now = new Date();
  const expiry = new Date(expiryDate);

  // Check if token is still valid for at least 10 seconds
  if (expiry > now && expiry - now > 10000) {
    return Promise.resolve(token);
  }

  return updateToken().then(() => {
    const updatedAdminData = getAdminData();
    return updatedAdminData?.token;
  });
};
