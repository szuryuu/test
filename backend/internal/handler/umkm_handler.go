package handler

import (
	"kasiraiai/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type UmkmHandler struct {
	repo repository.UmkmRepo
}

func NewUmkmHandler(repo repository.UmkmRepo) *UmkmHandler {
	return &UmkmHandler{repo: repo}
}

// Placeholder — akan diisi di Step 4
func (h *UmkmHandler) GetProfile(c *gin.Context) {}
