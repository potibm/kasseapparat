import React, { useState, useEffect, useMemo } from "react";
import {
  Navigate,
  createBrowserRouter,
  RouterProvider,
  RouteObject,
} from "react-router";
import { useAuth } from "../apps/pos/features/auth/providers/AuthProvider";
import { ProtectedRoute } from "./ProtectedRoute";
import Kasseapparat from "../apps/pos/pages/Pos";
import Admin from "../apps/admin/Admin";
import Logout from "../apps/pos/features/auth/components/Logout";
import Login from "../apps/pos/features/auth/components/Login";
import ChangePassword from "../apps/pos/features/auth/components/ChangePassword";
import NotFound from "../apps/pos/pages/NotFound";
import ForgotPassword from "../apps/pos/features/auth/components/ForgotPassword";
import { LoggedinErrorMessage } from "../apps/pos/features/auth/components/LoggedinErrorMessage";

const Routes: React.FC = () => {
  const { isLoggedIn } = useAuth() || {};

  const [loggedIn, setLoggedIn] = useState<boolean | null>(
    typeof isLoggedIn === "function" ? null : false,
  );

  useEffect(() => {
    if (!isLoggedIn) return;

    let isMounted = true; // Cleanup-Pattern, um Memory Leaks zu verhindern

    isLoggedIn()
      .then((status) => {
        if (isMounted) setLoggedIn(status);
      })
      .catch(() => {
        if (isMounted) setLoggedIn(false);
      });

    return () => {
      isMounted = false;
    };
  }, [isLoggedIn]);

  const router = useMemo(() => {
    // Define public routes accessible to all users
    const routesForPublic: RouteObject[] = [
      {
        path: "/admin/*",
        element: <Admin />,
      },
      {
        path: "/change-password",
        element: <ChangePassword />,
      },
    ];

    // Define routes accessible only to authenticated users
    const routesForAuthenticatedOnly: RouteObject[] = [
      {
        path: "/",
        element: <ProtectedRoute />,
        children: [
          {
            path: "/",
            element: <Kasseapparat />,
          },
          {
            path: "/logout",
            element: <Logout />,
          },
          {
            path: "/login",
            element: <Navigate to="/" replace />,
          },
          {
            path: "/forgot-password",
            element: <LoggedinErrorMessage />,
          },
        ],
      },
    ];

    // Define routes accessible only to non-authenticated users
    const routesForNotAuthenticatedOnly: RouteObject[] = [
      {
        path: "/login",
        element: <Login />,
      },
      {
        path: "/forgot-password",
        element: <ForgotPassword />,
      },
    ];

    const notFoundRoute: RouteObject[] = [
      {
        path: "*",
        element: <NotFound />,
      },
    ];

    return createBrowserRouter([
      ...routesForPublic,
      ...(loggedIn === true ? [] : routesForNotAuthenticatedOnly),
      ...routesForAuthenticatedOnly,
      ...notFoundRoute,
    ]);
  }, [loggedIn]);

  if (loggedIn === null) {
    return <div>⏳ Loading routes...</div>;
  }

  // Provide the router configuration using RouterProvider
  return <RouterProvider router={router} />;
};

export default Routes;
