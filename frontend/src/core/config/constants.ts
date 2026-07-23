export const API_VERSION = "v3";

export const buildApiBaseUrl = (apiHost: string): string =>
  `${apiHost}/api/${API_VERSION}`;
