import React, { useMemo } from "react";
import { createBrowserRouter, RouterProvider, RouteObject } from "react-router";
import { useAuth } from "../apps/pos/features/auth/hooks/useAuth";
import { ProtectedRoute } from "./ProtectedRoute";
import Kasseapparat from "../apps/pos/pages/Pos";
import Admin from "../apps/admin/Admin";
import NotFound from "../apps/pos/pages/NotFound";
import Splash from "../apps/splash/pages/Splash";

const Routes: React.FC = () => {
  const { isAuthenticated, isLoading } = useAuth();

  const router = useMemo(() => {
    const routesForPublic: RouteObject[] = [
      {
        path: "/",
        element: <Splash />,
      },
      {
        path: "/admin/*",
        element: <Admin />,
      },
    ];

    const routesForAuthenticatedOnly: RouteObject[] = [
      {
        element: <ProtectedRoute />,
        children: [
          {
            path: "/pos",
            element: <Kasseapparat />,
          },
        ],
      },
    ];

    const routesForNotAuthenticatedOnly: RouteObject[] = [];

    const notFoundRoute: RouteObject[] = [
      {
        path: "*",
        element: <NotFound />,
      },
    ];

    return createBrowserRouter([
      ...routesForPublic,
      ...(isAuthenticated ? [] : routesForNotAuthenticatedOnly),
      ...routesForAuthenticatedOnly,
      ...notFoundRoute,
    ]);
  }, [isAuthenticated]);

  if (isLoading) {
    return <div>⏳ Loading routes...</div>;
  }

  return <RouterProvider router={router} />;
};

export default Routes;
