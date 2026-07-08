import { AuthProvider } from "react-admin";

export const authProvider: AuthProvider = {
  login: () => Promise.resolve(),

  logout: () => Promise.resolve(),

  checkAuth: () => Promise.resolve(),

  checkError: () => Promise.resolve(),

  getPermissions: () => Promise.resolve("admin"),
};

export default authProvider;
