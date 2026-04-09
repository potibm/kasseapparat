import { use } from "react";
import { AuthContext } from "../context/AuthContext";
import { AuthContextType } from "../types/auth.types";

export const useAuth = (): AuthContextType => {
  const context = use(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export default useAuth;
