export const fetchProducts = async (apiHost) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/products?_end=1000&_sort=pos&_order=asc`)
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => resolve(data))
    .catch(error => reject(error));
});
}

export const storePurchase = async (apiHost, cart, totalPrice) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/purchases`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ cart, totalPrice: cart.reduce((total, item) => total + item.totalPrice, 0) }) // : )
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => resolve(data))
    .catch(error => reject(error));
});
}

export const fetchPurchases = async (apiHost) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/purchases`)
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => resolve(data))
    .catch(error => reject(error));
});
}

export const deletePurchaseById = async (apiHost, purchaseId) => {
  return new Promise((resolve, reject) => {
    fetch(`${apiHost}/api/v1/purchases/${purchaseId}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => resolve(data))
    .catch(error => reject(error));
});
}