import * as Sentry from "@sentry/react";

// 🔁 Shared error handler for failed fetch responses
const handleFetchError = async (response) => {
  let message = `HTTP ${response.status} ${response.statusText}`;
  let data = null;

  try {
    data = await response.json();
    message = data?.message || data?.error || message;
  } catch {
    // Ignore JSON parse errors
  }

  // Normalize message for comparison
  const normalizedMessage = message.toLowerCase();

  // 🧽 Filter known, non-critical messages
  const knownNonCritical = [
    "token is expired",
    "incorrect username or password",
  ];

  const isExpected = knownNonCritical.some((msg) =>
    normalizedMessage.includes(msg),
  );

  const error = new Error(message);

  if (!isExpected) {
    Sentry.captureException(error, {
      extra: {
        url: response.url,
        status: response.status,
        data,
      },
    });
  }

  throw data || error;
};

// 🔐 Authenticate user and retrieve JWT token
export const getJwtToken = async (apiHost, login, password) => {
  const response = await fetch(`${apiHost}/login`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ login, password }),
  });

  if (!response.ok) {
    await handleFetchError(response);
  }

  return response.json();
};

// 🔄 Refresh JWT token using refresh endpoint
export const refreshJwtToken = async (apiHost, refreshToken) => {
  const response = await fetch(`${apiHost}/auth/refresh_token`, {
    method: "GET",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${refreshToken}`,
    },
  });

  if (!response.ok) {
    await handleFetchError(response);
  }

  return response.json();
};

// 🔑 Change password using reset token
export const changePassword = async (apiHost, userId, token, newPassword) => {
  const response = await fetch(`${apiHost}/api/v2/auth/changePassword`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      userId: parseInt(userId),
      token,
      password: newPassword,
    }),
  });

  if (!response.ok) {
    await handleFetchError(response);
  }

  return response.json();
};

// 📧 Request password reset token for login name
export const requestChangePasswordToken = async (apiHost, login) => {
  const response = await fetch(`${apiHost}/api/v2/auth/changePasswordToken`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ login }),
  });

  if (!response.ok) {
    await handleFetchError(response);
  }

  return response.json();
};
