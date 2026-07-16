import React from "react";
import { Outlet } from "react-router";
import { useAuth } from "../apps/pos/features/auth/hooks/useAuth";

export const ProtectedRoute: React.FC = () => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <div>⏳ Loading Permissions...</div>;
  }

  if (!isAuthenticated) {
    window.location.href = "/api/v2/auth/login";
    return;
  }

  return <Outlet />;
};
