import jsonServerProvider from "ra-data-json-server";
import { fetchUtils } from "react-admin";
import * as Sentry from "@sentry/react";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3001";

const resourceAlias = {
  sumupReaders: "sumup/readers",
  sumupTransactions: "sumup/transactions",
};

const resolveResource = (resource) => resourceAlias[resource] || resource;

const httpClient = async (url, options = {}) => {
  if (!options.headers) {
    options.headers = new Headers();
    options.headers.set("Accept", "application/json");
  }
  if (!options.isUpload) {
    options.headers.set("Content-Type", "application/json");
  }
  // add your own headers here
  const adminData = localStorage.getItem("admin");
  if (adminData) {
    const parsedAdminData = JSON.parse(adminData);
    if (parsedAdminData?.token) {
      options.headers.set("Authorization", `Bearer ${parsedAdminData.token}`);
    }
  }

  try {
    return await fetchUtils.fetchJson(url, options);
  } catch (error) {
    const normalizedMessage = error?.message?.toLowerCase() || "";

    // 🧽 Filter known, non-critical messages
    const knownNonCritical = ["cookie token is empty"];

    const isExpected = knownNonCritical.some((msg) =>
      normalizedMessage.includes(msg),
    );

    if (!isExpected) {
      Sentry.captureException(error, {
        tags: {
          url,
          method: options.method || "GET",
        },
        extra: {
          request: {
            url,
            options,
          },
        },
      });
    }
  }
};

const baseProvider = jsonServerProvider(`${API_HOST}/api/v2`, httpClient);

const dataProvider = {
  ...baseProvider,

  getList: (resource, params) =>
    baseProvider.getList(resolveResource(resource), params),

  getOne: (resource, params) =>
    baseProvider.getOne(resolveResource(resource), params),

  getMany: (resource, params) =>
    baseProvider.getMany(resolveResource(resource), params),

  getManyReference: (resource, params) =>
    baseProvider.getManyReference(resolveResource(resource), params),

  create: (resource, params) =>
    baseProvider.create(resolveResource(resource), params),

  update: (resource, params) =>
    baseProvider.update(resolveResource(resource), params),

  updateMany: (resource, params) =>
    baseProvider.updateMany(resolveResource(resource), params),

  delete: (resource, params) =>
    baseProvider.delete(resolveResource(resource), params),

  deleteMany: (resource, params) =>
    baseProvider.deleteMany(resolveResource(resource), params),

  upload: (resource, params) => {
    const url = `${API_HOST}/api/v2/${resolveResource(resource)}`;
    const options = {
      method: "POST",
      body: params.data,
      isUpload: true,
    };
    return httpClient(url, options).then(({ json }) => ({
      data: json,
    }));
  },

  refund: (resource, params) => {
    if (!["purchases"].includes(resource)) {
      throw new Error(`Refund is not supported for resource: ${resource}`);
    }
    if (!params.id) {
      throw new Error("Refund requires an id");
    }

    const url = `${API_HOST}/api/v2/${resolveResource(resource)}/${params.id}/refund`;
    const options = {
      method: "POST",
      body: JSON.stringify(params.data),
    };
    return httpClient(url, options).then(({ json }) => ({
      data: json,
    }));
  },
};

export default dataProvider;
