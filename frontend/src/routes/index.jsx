import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import { useAuth } from '../provider/authProvider'
import { ProtectedRoute } from './ProtectedRoute'
import React from 'react'
import Kasseapparat from '../Main/Kasseapparat'
import Admin from '../Admin/Admin'
import Logout from '../Main/components/Logout'
import Login from '../Main/components/Login'

const Routes = () => {
  const { token } = useAuth()

  // Define public routes accessible to all users
  const routesForPublic = [
    {
      path: '/admin/*',
      element: <Admin />
    }
  ]

  // Define routes accessible only to authenticated users
  const routesForAuthenticatedOnly = [
    {
      path: '/',
      element: <ProtectedRoute />, // Wrap the component in ProtectedRoute
      children: [
        {
          path: '/',
          element: <Kasseapparat />
        },
        {
          path: '/logout',
          element: <Logout />
        }
      ]
    }
  ]

  // Define routes accessible only to non-authenticated users
  const routesForNotAuthenticatedOnly = [
    {
      path: '/login',
      element: <Login />
    }
  ]

  // Combine and conditionally include routes based on authentication status
  const router = createBrowserRouter([
    ...routesForPublic,
    ...(!token ? routesForNotAuthenticatedOnly : []),
    ...routesForAuthenticatedOnly
  ])

  // Provide the router configuration using RouterProvider
  return <RouterProvider router={router} />
}

export default Routes
