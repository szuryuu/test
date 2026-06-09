import api from "./axios";

export const transactionAPI = {
  list: (params) => api.get("/transactions", { params }),

  create: (data) => api.post("/transactions", data),

  remove: (id) => api.delete(`/transactions/${id}`),
};
