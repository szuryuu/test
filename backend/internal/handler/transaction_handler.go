package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"kasiraiai/backend/internal/model"
	"kasiraiai/backend/internal/repository"
	"kasiraiai/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) List(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := repository.TransactionFilter{
		Type:      c.Query("type"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Page:      page,
		Limit:     limit,
	}

	txs, total, err := h.svc.FindByUmkmID(c.Request.Context(), umkmID, filter)
	if err != nil {
		slog.Error("gagal list transaksi", "umkm_id", umkmID, "error", err)
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar transaksi berhasil diambil", gin.H{
		"transactions": txs,
		"total":        total,
		"page":         page,
	})
}

func (h *TransactionHandler) Create(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	var req struct {
		Amount          int64  `json:"amount" binding:"required,gt=0"`
		Type            string `json:"type" binding:"required,oneof=income expense"`
		Description     string `json:"description" binding:"required"`
		TransactionDate string `json:"transaction_date" binding:"required"`
		CategoryID      string `json:"category_id,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, ErrInvalidInput)
		return
	}

	var catID *uuid.UUID
	if req.CategoryID != "" {
		parsed, err := uuid.Parse(req.CategoryID)
		if err == nil {
			catID = &parsed
		}
	}

	tx := &model.Transaction{
		UmkmID:          umkmID,
		Amount:          req.Amount,
		Type:            req.Type,
		Description:     req.Description,
		RawMessage:      req.Description,
		TransactionDate: req.TransactionDate,
		Source:          "manual",
		CategoryID:      catID,
	}

	if err := h.svc.Create(c.Request.Context(), tx); err != nil {
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusCreated, "Transaksi berhasil dibuat", gin.H{
		"transaction": tx,
	})
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequest(c, ErrInvalidInput)
		return
	}

	if err := h.svc.SoftDelete(c.Request.Context(), id, umkmID); err != nil {
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusOK, "Transaksi berhasil dihapus", nil)
}
