export const fetchProducts = async (apiHost) => {
  return new Promise((resolve, reject) => {
    fetch(
      `${apiHost}/api/v1/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true`,
    )
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

export const fetchGuestListByProductId = async (apiHost, productId, query) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/products/${productId}/listEntries?q=${query}`)
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

export const storePurchase = async (apiHost, jwtToken, cart) => {
  return new Promise((resolve, reject) => {
    // null the cart items list property to avoid unnecessary data transfer
    const cartPayload = cart;
    cartPayload.forEach((item) => {
      item.lists = null;
    });

    fetch(`${apiHost}/api/v1/purchases`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwtToken}`,
      },
      body: JSON.stringify({
        cart: cartPayload,
        totalPrice: cart.reduce((total, item) => total + item.totalPrice, 0),
      }), // : )
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

export const fetchPurchases = async (apiHost) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/purchases`)
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

export const deletePurchaseById = async (apiHost, jwtToken, purchaseId) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/purchases/${purchaseId}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwtToken}`,
      },
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

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
          throw new Error("Network response was not ok");
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
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};

export const changePassword = async (
  apiHost,
  jwtToken,
  newPassword,
  newPasswordRepeat,
) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/auth/change_password`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwtToken}`,
      },
      body: JSON.stringify({
        password1: newPassword,
        password2: newPasswordRepeat,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => resolve(data))
      .catch((error) => reject(error));
  });
};
