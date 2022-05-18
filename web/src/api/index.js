const host = process.env.REACT_APP_API_HOST || "127.0.0.1:8080";
const protocol = process.env.REACT_APP_API_PROTOCOL || "http";
const isSecured = protocol.endsWith("s");
const baseUrl = `${protocol}://${host}`;

const apiPath = `/api`

const apiUrl = `${baseUrl}${apiPath}/`;
const baseWs = `ws${isSecured ? "s" : ""}://${host}`;
const wsUrl = `${baseWs}${apiPath}/listen`;

export const url = {
  ws: wsUrl,
  api: apiUrl,
};

export const endpoints = {
  login: "auth",
  register: "register",
  apiKeys: "keys",
  projects: "projects",
  session: "session",
  spec: "spec",
  next: (id, opts = {}) =>
    `/${id}/next${opts ? `?${new URLSearchParams(opts)}` : ""}`,
  sessions: (id, pagination = {}) =>
    `/${id}/sessions${pagination ? `?${new URLSearchParams(pagination)}` : ""}`,
  listen: (to) => `${wsUrl}?${new URLSearchParams(to)}`,
};




