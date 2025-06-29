package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware(allowedKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		for _, allowed := range allowedKeys {
			if key == strings.TrimSpace(allowed) && key != "" {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
	}
}
