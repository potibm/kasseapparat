export const fetchProducts = async (apiHost, jwtToken) => {
  return new Promise((resolve, reject) => {
    fetch(
      `${apiHost}/api/v1/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwtToken}`,
        },
      },
    )
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

export const fetchGuestListByProductId = async (
  apiHost,
  jwtToken,
  productId,
  query,
) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/products/${productId}/listEntries?q=${query}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwtToken}`,
      },
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

export const fetchPurchases = async (apiHost, jwtToken) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/purchases`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwtToken}`,
      },
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

export const addProductInterest = async (apiHost, jwtToken, productId) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/productInterests`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwtToken}`,
      },
      body: JSON.stringify({
        productId,
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
