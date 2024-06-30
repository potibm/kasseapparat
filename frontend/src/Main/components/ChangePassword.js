import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../provider/AuthProvider";
import { Label, Button, TextInput, Alert } from "flowbite-react";
import { changePassword } from "../hooks/Api";
import AuthCard from "./AuthCard";

const API_HOST = process.env.REACT_APP_API_HOST;

const ChangePassword = () => {
  const [password, setPassword] = useState("");
  const [password2, setPassword2] = useState("");
  const [validationMessage, setValidationMessage] = useState("");
  const navigate = useNavigate();
  const { token } = useAuth();

  const handleChangePassword = (e) => {
    e.preventDefault();
    // Überprüfen, ob die Passwörter mindestens 8 Zeichen lang sind
    if (password.length < 8 || password2.length < 8) {
      setValidationMessage("The password must be at least 8 characters long.");
      return;
    }
    // Überprüfen, ob die Passwörter übereinstimmen
    if (password !== password2) {
      setValidationMessage("The passwords do not match.");
      return;
    }

    changePassword(API_HOST, token, password, password2)
      .then((auth) => {
        navigate("/logout", { replace: true });
      })
      .catch((error) => {
        console.error("Login error:", error);
        setValidationMessage(
          "The password could not be changed. Please try again.",
        );
      });

    setValidationMessage(""); // Zurücksetzen der Validierungsnachricht
  };

  return (
    <AuthCard title="Change Password">
      {validationMessage && <Alert color="failure">{validationMessage}</Alert>}

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
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="password2" value="Repeat your password" />
          </div>
          <TextInput
            id="password2"
            type="password"
            required
            value={password2}
            onChange={(e) => setPassword2(e.target.value)}
          />
        </div>
        <Button type="submit">Change password</Button>
        <Button
          type="cancel"
          color="warning"
          onClick={() => navigate("/logout")}
        >
          Cancel
        </Button>
      </form>
    </AuthCard>
  );
};

export default ChangePassword;
