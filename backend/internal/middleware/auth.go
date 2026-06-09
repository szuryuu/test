package middleware

import (
	"net/http"
	"strings"

	"kasiraiai/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// Auth adalah middleware validasi JWT.
// Token diambil dari header Authorization: Bearer <token>.
// Token divalidasi melalui AuthService.ValidateToken.
// umkm_id disimpan di context untuk digunakan handler selanjutnya.
func Auth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Anda belum login atau sesi telah habis",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Anda belum login atau sesi telah habis",
			})
			return
		}

		tokenStr := parts[1]
		umkmID, err := authService.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Anda belum login atau sesi telah habis",
			})
			return
		}

		c.Set("umkm_id", umkmID.String())
		c.Next()
	}
}


