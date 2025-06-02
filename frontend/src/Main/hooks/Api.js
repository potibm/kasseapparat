import Decimal from "decimal.js";

export const fetchProducts = async (apiHost, jwtToken) => {
  return new Promise((resolve, reject) => {
    fetch(
      `${apiHost}/api/v2/products?_end=1000&_sort=pos&_order=asc&_filter_hidden=true`,
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

export const fetchGuestlistByProductId = async (
  apiHost,
  jwtToken,
  productId,
  query,
) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v2/products/${productId}/guests?q=${query}`, {
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

export const storePurchase = async (
  apiHost,
  jwtToken,
  cart,
  paymentMethodCode,
  paymentMethodData = {},
) => {
  return new Promise((resolve, reject) => {
    // null the cart items list property to avoid unnecessary data transfer
    const cartPayload = cart;
    cartPayload.forEach((item) => {
      item.lists = null;
    });

    let payloadData = {
      paymentMethod: paymentMethodCode,
      cart: cartPayload,
      totalGrossPrice: cart.reduce(
        (total, item) => total.add(item.totalGrossPrice),
        new Decimal(0),
      ),
      totalNetPrice: cart.reduce(
        (total, item) => total.add(item.totalNetPrice),
        new Decimal(0),
      ),
    };

    // merge additional payment method data if provided
    if (Object.keys(paymentMethodData).length > 0) {
      payloadData = { ...payloadData, ...paymentMethodData };
    }

    fetch(`${apiHost}/api/v2/purchases`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwtToken}`,
      },
      body: JSON.stringify(payloadData),
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

export const fetchPurchases = async (apiHost, jwtToken, userId) => {
  return new Promise((resolve, reject) => {
    fetch(
      `${apiHost}/api/v2/purchases?createdById=${encodeURIComponent(userId)}`,
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

export const deletePurchaseById = async (apiHost, jwtToken, purchaseId) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v2/purchases/${purchaseId}`, {
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
    fetch(`${apiHost}/api/v2/productInterests`, {
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
