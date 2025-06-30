package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware(allowedKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		log.Printf("[DEBUG] Received API key: %s", key)
		for _, allowed := range allowedKeys {
			if key == strings.TrimSpace(allowed) && key != "" {
				log.Printf("[DEBUG] API key is valid: %s", key)
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
	}
}
