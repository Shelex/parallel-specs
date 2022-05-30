import jwtDecode from "jwt-decode";

const authKey = "auth";

export const auth = {
  set(token) {
    token && localStorage.setItem(authKey, token);
  },
  get() {
    return localStorage.getItem(authKey);
  },
  logout() {
    localStorage.removeItem(authKey);
  },
  info() {
    const token = this.get();
    if (!token) {
      return {};
    }
    return jwtDecode(token);
  },
};
