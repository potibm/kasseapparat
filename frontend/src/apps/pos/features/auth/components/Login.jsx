import React, { useState } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../providers/auth-provider";
import { Label, Button, TextInput, Alert, Spinner } from "flowbite-react";
import { getJwtToken } from "../hooks/api";
import BaseCard from "../../../components/BaseCard";
import { useConfig } from "../../../../../core/config/providers/config-provider";

const Login = () => {
  const [error, setError] = useState(null);
  const [disabled, setDisabled] = useState(false);

  const { setSession, setUserdata } = useAuth();
  const navigate = useNavigate();
  const apiHost = useConfig().apiHost;

  const handleLogin = (event) => {
    if (disabled) {
      return;
    }
    setDisabled(true);
    event.preventDefault();

    const login = event.target.login.value;
    const password = event.target.password.value;

    getJwtToken(apiHost, login, password)
      .then((auth) => {
        const token = auth.access_token;
        const expiresIn = auth.expires_in;
        setSession(token, expiresIn);

        const userdata = auth;
        delete userdata.token;
        delete userdata.expire;
        delete userdata.code;
        setUserdata(userdata);
        console.log("Logged in, expires in: " + expiresIn);
        console.log("Userdata:", userdata);

        navigate("/", { replace: true });
      })
      .catch((error) => {
        setError({
          message: "There was an error logging you in.",
          details: error.message,
        });
        setDisabled(false);
      });
  };

  return (
    <BaseCard title="Login" linkForgotPassword={true}>
      {error && (
        <Alert color="failure" className="mb-2">
          <>
            <div>{error.message}</div>
            {error.details && (
              <div className="mt-2 font-mono">{error.details}</div>
            )}
          </>
        </Alert>
      )}

      <form className="flex flex-col gap-4" onSubmit={handleLogin}>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="username" value="Your username" />
          </div>
          <TextInput id="login" type="text" placeholder="Username" required />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="password" value="Your password" />
          </div>
          <TextInput
            id="password"
            type="password"
            placeholder="Password"
            required
          />
        </div>
        <Button type="submit" disabled={disabled}>
          Login {disabled && <Spinner className="ml-3" />}
        </Button>
      </form>
    </BaseCard>
  );
};

export default Login;
