package handler

import (
	"net/http"

	"kasiraiai/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	svc service.ReportService
}

func NewReportHandler(svc service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) GenerateMonthly(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	var req struct {
		Year  int `json:"year" binding:"required"`
		Month int `json:"month" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, ErrInvalidInput)
		return
	}

	report, err := h.svc.GenerateMonthly(c.Request.Context(), umkmID, req.Year, req.Month)
	if err != nil {
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusOK, "Laporan berhasil dibuat", report)
}
