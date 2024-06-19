import jsonServerProvider from "ra-data-json-server";
import { fetchUtils } from "react-admin";

const API_HOST = process.env.REACT_APP_API_HOST ?? "http://localhost:3001";

const httpClient = (url, options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: "application/json" });
  }
  // add your own headers here
  const adminData = localStorage.getItem("admin");
  if (adminData) {
    const parsedAdminData = JSON.parse(adminData);
    if (parsedAdminData && parsedAdminData.token) {
      options.headers.set("Authorization", `Bearer ${parsedAdminData.token}`);
    }
  }

  return fetchUtils.fetchJson(url, options);
};

const dataProvider = jsonServerProvider(API_HOST + "/api/v1", httpClient);

export default dataProvider;
