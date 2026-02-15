import { Navigate, createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router-dom";
import { useAuth } from "../Auth/provider/AuthProvider";
import { ProtectedRoute } from "./ProtectedRoute";
import React from "react";
import Kasseapparat from "../Main/Kasseapparat";
import Admin from "../Admin/Admin";
import Logout from "../Auth/components/Logout";
import Login from "../Auth/components/Login";
import ChangePassword from "../Auth/components/ChangePassword";
import NotFound from "../components/NotFound";
import ForgotPassword from "../Auth/components/ForgotPassword";
import { LoggedinErrorMessage } from "./LoggedinErrorMessage";

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
