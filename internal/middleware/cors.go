package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware allows cross-origin requests from the Next.js frontend.
// Origin whitelist dibaca dari env ALLOWED_ORIGINS (comma-separated).
// Jika ALLOWED_ORIGINS kosong, fallback ke wildcard "*" untuk dev.
func CORSMiddleware() gin.HandlerFunc {
	allowedEnv := os.Getenv("ALLOWED_ORIGINS")
	var whitelist []string
	if allowedEnv != "" {
		for _, o := range strings.Split(allowedEnv, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				whitelist = append(whitelist, trimmed)
			}
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if len(whitelist) == 0 {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if origin != "" && isAllowedOrigin(origin, whitelist) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func isAllowedOrigin(origin string, whitelist []string) bool {
	for _, allowed := range whitelist {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
