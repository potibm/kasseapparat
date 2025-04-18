import jsonServerProvider from "ra-data-json-server";
import { fetchUtils } from "react-admin";

const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3001";

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

const dataProvider = jsonServerProvider(`${API_HOST}/api/v2`, httpClient);

const myDataProvider = {
  ...dataProvider,
  upload: (resource, params) => {
    const url = `${API_HOST}/api/v2/${resource}`;
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

export default myDataProvider;
