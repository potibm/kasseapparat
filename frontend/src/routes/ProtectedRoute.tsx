import React from "react";
import { Outlet } from "react-router";
import { useAuth } from "../apps/pos/features/auth/hooks/useAuth";
import useConfig from "@core/config/hooks/useConfig";

export const ProtectedRoute: React.FC = () => {
  const { isAuthenticated, isLoading } = useAuth();
  const { apiBaseUrl } = useConfig();

  if (isLoading) {
    return <div>⏳ Loading Permissions...</div>;
  }

  if (!isAuthenticated) {
    window.location.href = `${apiBaseUrl}/auth/login`;
    return;
  }

  return <Outlet />;
};
