export interface AuthUser {
  id: number;
  username: string;
  email?: string;
  role: "admin" | "user" | "guest";
  gravatarUrl: string;
}

export interface AuthContextType {
  getToken: () => Promise<string | null>;
  getSafeToken: () => Promise<string>;
  isLoggedIn: () => Promise<boolean>;
  setSession: (token: string, expiresIn: number) => string;
  removeSession: () => void;
  userdata: AuthUser | null;
  setUserdata: (userdata: AuthUser) => void;
  gravatarUrl: string;
  role: string;
  username: string;
  id: number;
}

export interface Session {
  token: string | null;
  expiryDate: Date | null;
}

export interface LoginError {
  message: string;
  details?: string;
}
