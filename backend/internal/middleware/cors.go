package middleware

import (
	"time"

	"kasiraiai/backend/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS mengizinkan origin frontend. Di development, izinkan semua.
// Di production, batasi ke origin frontend spesifik.
func CORS(cfg *config.Config) gin.HandlerFunc {
	corsCfg := cors.DefaultConfig()

	if cfg.AppEnv == "production" {
		corsCfg.AllowOrigins = []string{cfg.FrontendURL}
	} else {
		corsCfg.AllowAllOrigins = true
	}

	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsCfg.ExposeHeaders = []string{"Content-Length"}
	corsCfg.AllowCredentials = true
	corsCfg.MaxAge = 12 * time.Hour

	return cors.New(corsCfg)
}
