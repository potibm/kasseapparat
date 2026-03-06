import { Navigate, createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router-dom";
import { useAuth } from "../apps/pos/features/auth/providers/auth-provider";
import { ProtectedRoute } from "./ProtectedRoute";
import React from "react";
import Kasseapparat from "../apps/pos/pages/Pos";
import Admin from "../apps/admin/Admin";
import Logout from "../apps/pos/features/auth/components/Logout";
import Login from "../apps/pos/features/auth/components/Login";
import ChangePassword from "../apps/pos/features/auth/components/ChangePassword";
import NotFound from "../apps/pos/pages/NotFound";
import ForgotPassword from "../apps/pos/features/auth/components/ForgotPassword";
import { LoggedinErrorMessage } from "../apps/pos/features/auth/components/LoggedinErrorMessage";

const Routes = () => {
  const { isLoggedIn } = useAuth() || {};
  const [loggedIn, setLoggedIn] = React.useState(null);

  React.useEffect(() => {
    if (!isLoggedIn) return;
    isLoggedIn()
      .then(setLoggedIn)
      .catch(() => setLoggedIn(false));
  }, [isLoggedIn]);

  // Define public routes accessible to all users
  const routesForPublic = [
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
  const routesForAuthenticatedOnly = [
    {
      path: "/",
      element: <ProtectedRoute />, // Wrap the component in ProtectedRoute
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
          element: <Navigate to="/" />,
        },
        {
          path: "/forgot-password",
          element: <LoggedinErrorMessage />,
        },
      ],
    },
  ];

  // Define routes accessible only to non-authenticated users
  const routesForNotAuthenticatedOnly = [
    {
      path: "/login",
      element: <Login />,
    },
    {
      path: "/forgot-password",
      element: <ForgotPassword />,
    },
  ];

  const notFoundRoute = [
    {
      path: "*",
      element: <NotFound />,
    },
  ];

  // Combine and conditionally include routes based on authentication status
  const router = createBrowserRouter([
    ...routesForPublic,
    ...(loggedIn === true ? [] : routesForNotAuthenticatedOnly),
    ...routesForAuthenticatedOnly,
    ...notFoundRoute,
  ]);

  // Provide the router configuration using RouterProvider
  return <RouterProvider router={router} />;
};

export default Routes;
