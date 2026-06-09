import api from "./axios";

export const dashboardAPI = {
  summary: (period, date) =>
    api.get("/dashboard/summary", { params: { period, date } }),

  categories: (period, type) =>
    api.get("/dashboard/categories", { params: { period, type } }),
};
