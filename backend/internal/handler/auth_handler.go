package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"kasiraiai/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, ErrInvalidInput)
		return
	}

	umkm, token, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrPhoneAlreadyExists) {
			ErrorResponse(c, http.StatusConflict, ErrPhoneAlreadyExists)
			return
		}
		slog.Error("gagal register UMKM", "error", err)
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusCreated, "Pendaftaran berhasil", gin.H{
		"token": token,
		"umkm":  umkm,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, ErrInvalidInput)
		return
	}

	umkm, token, err := h.svc.Login(c.Request.Context(), req.PhoneNumber, req.Password)
	if err != nil {
		Unauthorized(c, ErrInvalidCredential)
		return
	}

	SuccessResponse(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
		"umkm":  umkm,
	})
}
