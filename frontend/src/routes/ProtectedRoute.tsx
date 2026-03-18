import React, { useEffect, useState } from "react";
import { Navigate, Outlet } from "react-router";
import { useAuth } from "../apps/pos/features/auth/providers/AuthProvider";

export const ProtectedRoute: React.FC = () => {
  const { isLoggedIn } = useAuth() || {};
  const [loggedIn, setLoggedIn] = useState<boolean | null>(null);

  useEffect(() => {
    if (!isLoggedIn) return;
    isLoggedIn()
      .then(setLoggedIn)
      .catch(() => setLoggedIn(false));
  }, [isLoggedIn]);

  // Check if the user is authenticated
  if (loggedIn === false) {
    // If not authenticated, redirect to the login page
    return <Navigate to="/login" />;
  }

  if (loggedIn === null) {
    // Still loading
    return null;
  }

  // If authenticated, render the child routes
  return <Outlet />;
};
