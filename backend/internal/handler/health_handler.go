package handler

import (
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// IsReady di-set ke true setelah semua inisialisasi selesai.
// Digunakan oleh health check endpoint — Render memerlukan ini
// untuk menentukan apakah service siap menerima request.
var IsReady atomic.Bool

// HealthCheck mengembalikan status service.
// 200 OK jika siap, 503 jika masih starting.
func HealthCheck(c *gin.Context) {
	if !IsReady.Load() {
		c.JSON(503, gin.H{"status": "starting"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "kasiraiai",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}
