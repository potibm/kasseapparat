import React, { useState } from "react";
import { Label, Button, TextInput, Alert, Spinner } from "flowbite-react";
import BaseCard from "../../../components/BaseCard";
import { useConfig } from "../../../../../core/config/providers/ConfigProvider";
import { requestChangePasswordToken } from "../hooks/api";

const RequestToken: React.FC = () => {
  const [loginInput, setLoginInput] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);
  const [disabled, setDisabled] = useState<boolean>(false);

  const { apiHost } = useConfig();

  const handleForgotPassword = (event: React.SubmitEvent<HTMLFormElement>) => {
    if (disabled) {
      return;
    }
    setDisabled(true);

    event.preventDefault();

    requestChangePasswordToken(apiHost, loginInput)
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
                  <Label htmlFor="login">Your login</Label>
                </div>
                <TextInput
                  id="login"
                  type="text"
                  placeholder="login"
                  required
                  value={loginInput}
                  onChange={(e) => setLoginInput(e.target.value)}
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
