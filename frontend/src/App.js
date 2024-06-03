import REACT_APP_API_HOST from 'react'
import './App.css'
import Kasseapparat from './Main/Kasseapparat'
import Admin from './Admin/Admin'
import { createBrowserRouter, RouterProvider, BrowserRouter as Router, Route, Switch, Link } from 'react-router-dom';
import AuthProvider from "./provider/authProvider"
import Routes from "./routes";

function App () {

  return (
    <AuthProvider>
      <Routes />
    </AuthProvider>
  )
}

export default App