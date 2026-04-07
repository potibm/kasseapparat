import { createContext } from "react";
import { AuthContextType } from "../types/auth.types";

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);

export default AuthContext;
