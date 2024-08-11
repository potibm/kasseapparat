export const getJwtToken = async (apiHost, login, password) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Credentials: "include",
      },
      body: JSON.stringify({ login, password }),
    })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((errorBody) => {
            throw new Error(errorBody.message || "Network response was not ok");
          });
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

export const refreshJwtToken = async (apiHost, refreshToken) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/auth/refresh_token`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${refreshToken}`,
      },
    })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((errorBody) => {
            throw new Error(errorBody.message || "Network response was not ok");
          });
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

export const changePassword = async (apiHost, userId, token, newPassword) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/auth/changePassword`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        userId: parseInt(userId),
        token,
        password: newPassword,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((errorBody) => {
            throw new Error(errorBody.error || "Network response was not ok");
          });
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

export const requestChangePasswordToken = async (apiHost, login) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/auth/changePasswordToken`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        login,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((errorBody) => {
            throw new Error(errorBody.error || "Network response was not ok");
          });
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};
