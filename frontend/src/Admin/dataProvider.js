import jsonServerProvider from "ra-data-json-server";
import { fetchUtils } from "react-admin";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3001";

const resourceAlias = {
  sumupReaders: "sumup/readers",
};

const resolveResource = (resource) => resourceAlias[resource] || resource;

const httpClient = (url, options = {}) => {
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

  return fetchUtils.fetchJson(url, options);
};

//const dataProvider = jsonServerProvider(`${API_HOST}/api/v2`, httpClient);

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
};

export default dataProvider;
