import * as Sentry from "@sentry/react";
import { z } from "zod";
import {
  LoginResponseSchema,
  LoginResponse,
  SimpleResponseSchema,
  SimpleResponse,
  StringResponseSchema,
  StringResponse,
  RefreshTokenResponseSchema,
  RefreshTokenResponse,
} from "./api.schemas";

const authPost = async <S extends z.ZodTypeAny>(
  url: string,
  body: object | null,
  schema: S,
): Promise<z.infer<S>> => {
  const response = await fetch(url, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    await handleAuthError(response);
  }

  const data = await response.json();
  const result = schema.safeParse(data);

  if (!result.success) {
    console.error("Auth API Validation Error:", result.error);
    throw new Error("Invalid response format from Auth API");
  }

  return result.data;
};

// 🔁 Shared error handler for failed fetch responses
const handleAuthError = async (response: Response): Promise<never> => {
  let message = `HTTP ${response.status}`;
  let data: any = null;

  try {
    data = await response.json();
    message = data?.message || data?.error || message;
  } catch {
    /* ignore */
  }

  const normalizedMessage = message.toLowerCase();
  const knownNonCritical = [
    "token is expired",
    "incorrect username or password",
    "invalid credentials",
  ];
  const isExpected = knownNonCritical.some((msg) =>
    normalizedMessage.includes(msg),
  );

  if (!isExpected) {
    Sentry.captureException(new Error(message), {
      extra: { url: response.url, status: response.status, data },
    });
  }

  // Wir werfen ein Objekt, das die UI abfangen kann
  throw { message, status: response.status, data };
};

// 🔐 Authenticate user and retrieve JWT token
export const getJwtToken = (
  apiHost: string,
  login: string,
  password: string,
): Promise<LoginResponse> =>
  authPost(
    `${apiHost}/api/v2/auth/login`,
    { login, password },
    LoginResponseSchema,
  );

// 🔄 Refresh JWT token using refresh endpoint
export const refreshJwtToken = (
  apiHost: string,
): Promise<RefreshTokenResponse> =>
  authPost(`${apiHost}/api/v2/auth/refresh`, null, RefreshTokenResponseSchema);

// Logout user
export const logout = (apiHost: string): Promise<SimpleResponse> =>
  authPost(`${apiHost}/api/v2/auth/logout`, null, SimpleResponseSchema);

// 🔑 Change password using reset token
export const changePassword = (
  apiHost: string,
  userId: string | number,
  token: string,
  password: string,
): Promise<SimpleResponse> =>
  authPost(
    `${apiHost}/api/v2/auth/changePassword`,
    {
      userId: typeof userId === "string" ? parseInt(userId) : userId,
      token,
      password,
    },
    SimpleResponseSchema,
  );

// 📧 Request password reset token for login name
export const requestChangePasswordToken = (
  apiHost: string,
  login: string,
): Promise<StringResponse> =>
  authPost(
    `${apiHost}/api/v2/auth/changePasswordToken`,
    { login },
    StringResponseSchema,
  );
