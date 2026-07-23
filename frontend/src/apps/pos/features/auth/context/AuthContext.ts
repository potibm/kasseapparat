import { createContext } from "react";
import { AuthContextType } from "../types/auth.types";

export const AuthContext = createContext<AuthContextType>({
  isLoading: true,
  isAuthenticated: false,
  role: null,
  username: null,
});

export default AuthContext;
