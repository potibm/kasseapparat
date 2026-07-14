export interface AuthContextType {
  isLoading: boolean;
  isAuthenticated: boolean;
  role: string | null;
  username: string | null;
}
