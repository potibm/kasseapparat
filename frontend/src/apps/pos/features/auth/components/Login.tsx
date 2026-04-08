import React, { useState } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";
import { Label, Button, TextInput, Alert, Spinner } from "flowbite-react";
import { getJwtToken } from "../../../../../core/api/auth";
import BaseCard from "../../../components/BaseCard";
import { useConfig } from "@core/config/hooks/useConfig";
import {
  LoginError as LoginErrorType,
  AuthUser as AuthUserType,
} from "../types/auth.types";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Auth");

const Login: React.FC = () => {
  const [formData, setFormData] = useState<{
    login: string;
    password: string;
  }>({
    login: "",
    password: "",
  });
  const [error, setError] = useState<LoginErrorType | null>(null);
  const [disabled, setDisabled] = useState(false);

  const { setSession, setUserdata } = useAuth();
  const navigate = useNavigate();
  const apiHost = useConfig().apiHost;

  const handleLogin = (event: React.SubmitEvent<HTMLFormElement>) => {
    if (disabled) {
      return;
    }
    setDisabled(true);
    event.preventDefault();
    setError(null);

    getJwtToken(apiHost, formData.login, formData.password)
      .then((auth) => {
        const { access_token, expires_in, ...userdata } = auth;

        setSession(access_token, expires_in);
        setUserdata(userdata as AuthUserType);
        log.debug("Logged in, expires in", expires_in);

        navigate("/", { replace: true });
      })
      .catch((error: Error) => {
        log.error("Login failed", error);
        setError({
          message: "There was an error logging you in.",
          details: error.message,
        });
        setDisabled(false);
      });
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
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
            <Label htmlFor="login">Your username</Label>
          </div>
          <TextInput
            id="login"
            name="login"
            type="text"
            placeholder="Username"
            required
            value={formData.login}
            onChange={(e) => handleInputChange(e)}
          />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="password">Your password</Label>
          </div>
          <TextInput
            id="password"
            type="password"
            name="password"
            placeholder="Password"
            required
            value={formData.password}
            onChange={(e) => handleInputChange(e)}
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
