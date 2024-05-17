import REACT_APP_API_HOST from 'react'
import './App.css'
import Kasseapparat from './Main/Kasseapparat'
import Admin from './Admin/Admin'
import { createBrowserRouter, RouterProvider, BrowserRouter as Router, Route, Switch, Link } from 'react-router-dom';

function App () {

  var router = createBrowserRouter([
    {
      path: "/",
      element: <Kasseapparat />,
    },
    {
      path: "admin/*",
      element: <Admin />,
    }
  ])

  return (
    <div>
        <RouterProvider router={router} />
    </div>
  )
}

export default App