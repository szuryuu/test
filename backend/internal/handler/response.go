package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// fieldLabels memetakan nama field JSON ke label Bahasa Indonesia.
var fieldLabels = map[string]string{
	"Name":         "Nama Lengkap",
	"BusinessName": "Nama Usaha",
	"PhoneNumber":  "Nomor WhatsApp",
	"Password":     "Password",
	"Amount":       "Jumlah",
	"Type":         "Tipe",
	"Description":  "Deskripsi",
	"Year":         "Tahun",
	"Month":        "Bulan",
}

// tagMessages memetakan tag validator ke pesan Bahasa Indonesia.
var tagMessages = map[string]string{
	"required":   "wajib diisi",
	"min":        "minimal %s karakter",
	"max":        "maksimal %s karakter",
	"len":        "harus %s karakter",
	"startswith": "harus diawali %s",
	"numeric":    "harus berupa angka",
	"gte":        "minimal %s",
	"lte":        "maksimal %s",
	"oneof":      "harus salah satu dari: %s",
	"gt":         "harus lebih dari %s",
	"email":      "format email tidak valid",
}

// validationErrors mengekstrak error validasi dari Gin/validator menjadi
// pesan Bahasa Indonesia per-field. Mengembalikan nil jika err bukan
// validator.ValidationErrors.
func validationErrors(err error) []string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}

	messages := make([]string, 0, len(ve))
	for _, fe := range ve {
		label := fe.Field()
		if lbl, ok := fieldLabels[fe.Field()]; ok {
			label = lbl
		}

		tag := fe.Tag()
		param := fe.Param()

		var msg string
		if tmpl, ok := tagMessages[tag]; ok {
			if strings.Contains(tmpl, "%s") {
				msg = fmt.Sprintf(tmpl, param)
			} else {
				msg = tmpl
			}
		} else {
			msg = fmt.Sprintf("tidak valid (%s)", tag)
		}

		messages = append(messages, fmt.Sprintf("%s %s", label, msg))
	}
	return messages
}
