package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response helpers — semua response API harus menggunakan fungsi ini.

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, status int, message string, errors ...string) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// Konstanta pesan error dalam Bahasa Indonesia
const (
	ErrInvalidInput       = "Data yang dikirim tidak valid"
	ErrPhoneAlreadyExists = "Nomor WhatsApp sudah terdaftar"
	ErrInvalidCredential  = "Nomor WhatsApp atau password salah"
	ErrUnauthorized       = "Anda belum login atau sesi telah habis"
	ErrNotFound           = "Data tidak ditemukan"
	ErrInternalServer     = "Terjadi kesalahan sistem, coba lagi nanti"
)

// BadRequest, InternalServerError, NotFound, Unauthorized — convenience shortcuts.
func BadRequest(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusBadRequest, msg)
}

func InternalServerError(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusInternalServerError, msg)
}

func NotFound(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusNotFound, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusUnauthorized, msg)
}

// getUmkmID mengekstrak umkm_id dari gin context (disimpan oleh middleware.Auth).
func getUmkmID(c *gin.Context) (uuid.UUID, error) {
	raw, exists := c.Get("umkm_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("umkm_id tidak ditemukan di context")
	}
	idStr, ok := raw.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("umkm_id bukan string")
	}
	return uuid.Parse(idStr)
}
