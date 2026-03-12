import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router";
import {
  Label,
  Button,
  TextInput,
  Alert,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  Spinner,
} from "flowbite-react";
import { changePassword } from "../../../../../core/api/auth";
import BaseCard from "../../../components/BaseCard";
import { useConfig } from "../../../../../core/config/providers/ConfigProvider";
import { z } from "zod";
import { AuthApiError as AuthApiErrorType } from "../../../utils/api.types";

interface ValidationMessageType {
  message: string;
  details: string | null;
  link?: string;
  linkText?: string;
}

const QueryParamsSchema = z.object({
  userId: z.string().regex(/^\d+$/).transform(Number),
  token: z.string().length(32),
});

const PasswordSchema = z
  .object({
    password: z
      .string()
      .min(8, "The password must be at least 8 characters long."),
    passwordRepeat: z.string(),
  })
  .refine((data) => data.password === data.passwordRepeat, {
    message: "The passwords do not match.",
    path: ["passwordRepeat"],
  });

const ChangePassword: React.FC = () => {
  const [formData, setFormData] = useState<{
    password: string;
    passwordRepeat: string;
  }>({
    password: "",
    passwordRepeat: "",
  });
  const [validationMessage, setValidationMessage] =
    useState<ValidationMessageType | null>(null);
  const [showSuccessModal, setShowSuccessModal] = useState<boolean>(false);
  const [disabled, setDisabled] = useState<boolean>(false);
  const navigate = useNavigate();
  const apiHost = useConfig().apiHost;

  const location = useLocation();

  const queryData = Object.fromEntries(new URLSearchParams(location.search));
  const queryResult = QueryParamsSchema.safeParse(queryData);

  if (!queryResult.success) {
    return (
      <BaseCard title="Change Password">
        <Alert color="failure">
          The link you were following is incorrect. Please, make sure you have
          the correct link.
        </Alert>
      </BaseCard>
    );
  }

  const { userId, token } = queryResult.data;

  const handleChangePassword = (event: React.SubmitEvent<HTMLFormElement>) => {
    if (disabled) {
      return;
    }

    event.preventDefault();
    setDisabled(true);
    setValidationMessage(null);

    const formResult = PasswordSchema.safeParse({
      password: formData.password,
      passwordRepeat: formData.passwordRepeat,
    });

    if (!formResult.success) {
      setValidationMessage({
        message: formResult.error.issues[0].message,
        details: null,
      });
      setDisabled(false);
      return;
    }

    changePassword(apiHost, userId, token, formData.password)
      .then(() => {
        setShowSuccessModal(true);
      })
      .catch((err: unknown) => {
        setDisabled(false);

        const error = err as AuthApiErrorType;

        const details = error.details || error.data?.details;
        const message = error.message || "An error occurred";

        if (details === "Token is invalid or has expired") {
          setValidationMessage({
            message: "The password could not be changed.",
            details: "The security token is no longer valid.",
            link: "/forgot-password",
            linkText: "Request new token",
          });
        } else {
          setValidationMessage({
            message: message,
            details: details || null,
          });
        }
      });

    setValidationMessage(null);
  };

  const handleModalClose = () => {
    setShowSuccessModal(false);
    navigate("/logout");
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
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
                color="red"
                className="mt-2"
                onClick={() => navigate(validationMessage.link!)}
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
            <Label htmlFor="password">Your new password</Label>
          </div>
          <TextInput
            id="password"
            name="password"
            type="password"
            placeholder="password"
            required
            value={formData.password}
            onChange={(e) => handleInputChange(e)}
          />
        </div>
        <div>
          <div className="mb-2 block">
            <Label htmlFor="passwordRepeat">Repeat your password</Label>
          </div>
          <TextInput
            id="passwordRepeat"
            name="passwordRepeat"
            type="password"
            required
            value={formData.passwordRepeat}
            onChange={(e) => handleInputChange(e)}
          />
        </div>
        <Button type="submit" disabled={disabled}>
          Change password {disabled && <Spinner className="ml-3" />}
        </Button>
        <Button
          type="reset"
          disabled={disabled}
          color="alternative"
          onClick={() => navigate("/")}
        >
          Cancel
        </Button>
      </form>
      <Modal show={showSuccessModal} onClose={handleModalClose}>
        <ModalHeader>Password Changed</ModalHeader>
        <ModalBody>
          <p>
            Your password has been successfully changed. You will be redirected
            to the login page.
          </p>
        </ModalBody>
        <ModalFooter>
          <Button onClick={() => handleModalClose()}>Perfect</Button>
        </ModalFooter>
      </Modal>
    </BaseCard>
  );
};

export default ChangePassword;
