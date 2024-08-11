import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import {
  Label,
  Button,
  TextInput,
  Alert,
  Modal,
  Spinner,
} from "flowbite-react";
import { changePassword } from "../hooks/Api";
import BaseCard from "../../components/BaseCard";
import { useConfig } from "../../provider/ConfigProvider";

const ChangePassword = () => {
  const [password, setPassword] = useState("");
  const [passwordRepeat, setPasswordRepeat] = useState("");
  const [validationMessage, setValidationMessage] = useState(null);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [disabled, setDisabled] = useState(false);
  const navigate = useNavigate();
  const apiHost = useConfig().apiHost;

  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const userId = searchParams.get("userId");
  const token = searchParams.get("token");

  if (
    userId === null ||
    isNaN(userId) ||
    token === null ||
    token.length !== 32
  ) {
    return (
      <BaseCard title="Change Password">
        <Alert color="failure">
          The link you were following is incorrect. Please, make sure you have
          the correct link.
        </Alert>
      </BaseCard>
    );
  }

  const handleChangePassword = (e) => {
    if (disabled) {
      return;
    }
    setDisabled(true);

    e.preventDefault();

    if (password.length < 8) {
      setValidationMessage({
        message: "The password must be at least 8 characters long.",
        details: null,
      });
      setDisabled(false);
      return;
    }
    if (password !== passwordRepeat) {
      setValidationMessage({
        message: "The passwords do not match.",
        details: null,
      });
      setDisabled(false);
      return;
    }

    changePassword(apiHost, userId, token, password)
      .then((auth) => {
        setShowSuccessModal(true);
      })
      .catch((error) => {
        setDisabled(false);
        if (error.message === "Token is invalid or has expired.") {
          setValidationMessage({
            message: "The password could not be changed. ",
            details: error.message,
            link: "/forgot-password",
            linkText: "Request new token",
          });
        } else {
          setValidationMessage({
            message: "The password could not be changed. Please try again.",
            details: error.message,
          });
        }
      });

    setValidationMessage(null); // Zurücksetzen der Validierungsnachricht
  };

  const handleModalClose = () => {
    setShowSuccessModal(false);
    navigate("/logout"); // Weiterleitung zur Login-Seite
  };

  return (
    <BaseCard title="Change Password" linkLogin={true}>
      {validationMessage && (
        <Alert color="failure" className="mb-2">
          <>
            <div>{validationMessage.message}</div>
            {validationMessage.details && (
              <div className="mt-2 font-mono">{validationMessage.details}</div>
            )}
            {validationMessage.link && (
              <Button
                color="failure"
                className="mt-2"
                onClick={() => navigate(validationMessage.link)}
              >
                {validationMessage.linkText}
              </Button>
            )}
          </>
        </Alert>
      )}

      <form className="flex flex-col gap-4" onSubmit={handleChangePassword}>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="password" value="Your new password" />
          </div>
          <TextInput
            id="password"
            type="password"
            placeholder="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value.trim())}
          />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="passwordRepeat" value="Repeat your password" />
          </div>
          <TextInput
            id="passwordRepeat"
            type="password"
            required
            value={passwordRepeat}
            onChange={(e) => setPasswordRepeat(e.target.value.trim())}
          />
        </div>
        <Button type="submit" disabled={disabled}>
          Change password {disabled && <Spinner className="ml-3" />}
        </Button>
        <Button
          type="cancel"
          disabled={disabled}
          color="warning"
          onClick={() => navigate("/")}
        >
          Cancel
        </Button>
      </form>
      <Modal show={showSuccessModal} onClose={handleModalClose}>
        <Modal.Header>Password Changed</Modal.Header>
        <Modal.Body>
          <p>
            Your password has been successfully changed. You will be redirected
            to the login page.
          </p>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={() => handleModalClose()}>Perfect</Button>
        </Modal.Footer>
      </Modal>
    </BaseCard>
  );
};

export default ChangePassword;
