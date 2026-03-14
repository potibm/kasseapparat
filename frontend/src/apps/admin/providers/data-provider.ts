import jsonServerProvider from "ra-data-json-server";
import {
  fetchUtils,
  addRefreshAuthToDataProvider,
  DataProvider,
} from "react-admin";
import * as Sentry from "@sentry/react";
import { refreshToken } from "./refresh-token";
import { getSessionToken } from "../utils/auth-utils";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3001";

const resourceAlias: Record<string, string> = {
  sumupReaders: "sumup/readers",
  sumupTransactions: "sumup/transactions",
};

const resolveResource = (resource: string) =>
  resourceAlias[resource] || resource;

const httpClient = async (url: string, options: fetchUtils.Options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: "application/json" });
  } else if (!(options.headers instanceof Headers)) {
    options.headers = new Headers(options.headers);
  }

  const headers = options.headers as Headers;

  if (!(options as any).isUpload) {
    headers.set("Content-Type", "application/json");
  }

  const token = getSessionToken();
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  try {
    return await fetchUtils.fetchJson(url, options);
  } catch (error: unknown) {
    let normalizedMessage = "";

    if (error instanceof Error) {
      normalizedMessage = error.message.toLowerCase();
    } else if (
      typeof error === "object" &&
      error !== null &&
      "message" in error
    ) {
      normalizedMessage = String((error as any).message).toLowerCase();
    }

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

    throw error;
  }
};

const baseProvider = jsonServerProvider(`${API_HOST}/api/v2`, httpClient);

const dataProvider: DataProvider = {
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

  upload: async (resource: string, params: any) => {
    const url = `${API_HOST}/api/v2/${resolveResource(resource)}`;
    const options = {
      method: "POST",
      body: params.data,
      isUpload: true,
    };
    const { json } = await httpClient(url, options);
    return { data: json };
  },

  refund: async (resource: string, params: any) => {
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

    const { json } = await httpClient(url, options);
    return { data: json };
  },
};

export default addRefreshAuthToDataProvider(dataProvider, refreshToken);
