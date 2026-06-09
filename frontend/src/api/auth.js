import api from "./axios";

export const authAPI = {
  login: (phoneNumber, password) =>
    api.post("/auth/login", { phone_number: phoneNumber, password }),

  register: (data) =>
    api.post("/auth/register", data),
};
