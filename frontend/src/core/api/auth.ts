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
} from "./auth.schemas";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Auth");

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
    log.error("Auth API Validation Error", result.error);
    throw new Error("Invalid response format from Auth API");
  }

  return result.data;
};

// 🔁 Shared error handler for failed fetch responses
const handleAuthError = async (response: Response): Promise<never> => {
  let message = `HTTP ${response.status}`;
  let data: unknown = null;

  try {
    data = await response.json();
    if (data && typeof data === "object") {
      const errorObj = data as Record<string, unknown>;
      const rawMessage = errorObj.message ?? errorObj.error;

      if (typeof rawMessage === "string") {
        message = rawMessage;
      } else if (response.statusText) {
        message = response.statusText;
      }
    }
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

  const apiError = new ApiError(message, response.status, data);

  if (!isExpected) {
    Sentry.captureException(apiError, {
      extra: {
        url: response.url,
        status: response.status,
        data,
      },
    });
  }

  throw apiError;
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
): Promise<SimpleResponse> => {
  const normalizedUserId = typeof userId === "string" ? Number(userId) : userId;

  if (!Number.isInteger(normalizedUserId)) {
    throw new TypeError("Invalid userId");
  }

  return authPost(
    `${apiHost}/api/v2/auth/changePassword`,
    {
      userId: normalizedUserId,
      token,
      password,
    },
    SimpleResponseSchema,
  );
};

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

export class ApiError extends Error {
  status: number;
  data: unknown;
  details: unknown;

  constructor(message: string, status: number, data: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.data = data;

    if (data && typeof data === "object") {
      this.details = (data as Record<string, unknown>).details;
    }

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ApiError);
    }
  }
}
