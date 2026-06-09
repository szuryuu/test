import api from "./axios";

export const kurAPI = {
  getScore: () => api.get("/kur/score"),
  recalculate: () => api.post("/kur/recalculate"),
};

export const reportAPI = {
  generate: (year, month) => api.post("/reports/monthly", { year, month }),
};
