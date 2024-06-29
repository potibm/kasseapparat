import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../provider/AuthProvider";
import { Label, Button, TextInput, Alert } from "flowbite-react";
import { getJwtToken } from "../hooks/Api";
import AuthCard from "./AuthCard";

const API_HOST = process.env.REACT_APP_API_HOST;

const Login = () => {
  const [error, setError] = useState(null);

  const { setToken, setUsername, setExpiryDate } = useAuth();
  const navigate = useNavigate();

  const handleLogin = (event) => {
    event.preventDefault();

    const username = event.target.username.value;
    const password = event.target.password.value;

    getJwtToken(API_HOST, username, password)
      .then((auth) => {
        const token = auth.token;
        const expiryDate = auth.expire;
        console.log(auth);
        const passwordChangeRequired = auth.passwordChangeRequired;
        console.log("Password change required", passwordChangeRequired);
        setToken(token);
        setUsername(username);
        setExpiryDate(expiryDate);
        navigate("/", { replace: true });
      })
      .catch((error) => {
        console.error("Login error:", error);
        setError("Invalid username or password. Please try again.");
      });
  };

  return (
    <AuthCard>
      {error && <Alert color="failure">{error}</Alert>}

      <form className="flex flex-col gap-4" onSubmit={handleLogin}>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="username" value="Your username" />
          </div>
          <TextInput
            id="username"
            type="text"
            placeholder="Username"
            required
          />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="password" value="Your password" />
          </div>
          <TextInput id="password" type="password" required />
        </div>
        <Button type="submit">Login</Button>
      </form>
    </AuthCard>
  );
};

export default Login;
