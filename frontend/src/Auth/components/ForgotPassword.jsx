import React, { useState } from "react";
import { Label, Button, TextInput, Alert, Spinner } from "flowbite-react";
import BaseCard from "../../components/BaseCard";
import { useConfig } from "../../provider/ConfigProvider";
import { requestChangePasswordToken } from "../hooks/Api";

const RequestToken = () => {
  const [login, setLogin] = useState("");
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);
  const [disabled, setDisabled] = useState(false);

  const apiHost = useConfig().apiHost;

  const handleForgotPassword = (event) => {
    if (disabled) {
      return;
    }
    setDisabled(true);

    event.preventDefault();

    requestChangePasswordToken(apiHost, login)
      .then(() => {
        setSuccess(true);
      })
      .catch((error) => {
        console.error("Error:", error);
        setError("There was an unknown error to request a password change.");
        setDisabled(false);
      });
  };

  return (
    <BaseCard title="Forgot Password" linkLogin={true}>
      <>
        {success && (
          <Alert color="success">
            When a login exists, an email has been sent to your email address
            with instructions on how to reset your password. Please check your
            inbox and spam folder.
            <br />
            <br />
            When a login does not exist, no email has been sent.
          </Alert>
        )}
        {!success && (
          <>
            <p className="mb-2 dark:text-white">
              Enter your username or email to receive an email that will guide
              you through the process to set a new password.
            </p>

            {error && <Alert color="failure">{error}</Alert>}

            <form
              className="flex flex-col gap-4"
              onSubmit={handleForgotPassword}
            >
              <div>
                <div className="mb-2 block">
                  <Label htmlFor="login" value="Your login" />
                </div>
                <TextInput
                  id="login"
                  type="text"
                  placeholder="login"
                  required
                  value={login}
                  onChange={(e) => setLogin(e.target.value)}
                />
              </div>

              <Button type="submit" disabled={disabled}>
                Forgot my password {disabled && <Spinner className="ml-3" />}
              </Button>
            </form>
          </>
        )}
      </>
    </BaseCard>
  );
};

export default RequestToken;
