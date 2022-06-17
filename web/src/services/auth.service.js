import jwtDecode from "jwt-decode";

export const authChangedEvent = "auth_changed";
const authKey = "auth";

const triggerAuthChange = () =>
  window.dispatchEvent(new Event(authChangedEvent));

export const auth = {
  set(token) {
    token && localStorage.setItem(authKey, token);
    triggerAuthChange();
  },
  get() {
    return localStorage.getItem(authKey);
  },
  logout() {
    localStorage.removeItem(authKey);
    triggerAuthChange();
  },
  info() {
    const token = this.get();
    if (!token) {
      return {};
    }
    return jwtDecode(token);
  },
};
