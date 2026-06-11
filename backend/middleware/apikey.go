package middleware

import (
	"net/http"
	"strings"
	"time"

	"ai-chat-platform/backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIKeyContextKey is where the verified key record is stored in Gin context.
const APIKeyContextKey = "api_key_record"

// APIKeyAuth validates bearer tokens and blocks disabled, expired, or empty-quota keys.
func APIKeyAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "missing bearer token",
					"type":    "invalid_request_error",
				},
			})
			return
		}

		accessKey := strings.TrimSpace(authorization[len("Bearer "):])
		if accessKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "empty bearer token",
					"type":    "invalid_request_error",
				},
			})
			return
		}

		// Load the key record once so downstream handlers can reuse it.
		var apiKey model.APIKey
		if err := db.Preload("Model").Where("access_key = ?", accessKey).First(&apiKey).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "invalid api key",
					"type":    "invalid_request_error",
				},
			})
			return
		}

		if !apiKey.Enabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"message": "api key is disabled",
					"type":    "permission_error",
				},
			})
			return
		}

		if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"message": "api key has expired",
					"type":    "permission_error",
				},
			})
			return
		}

		if !apiKey.UnlimitedQuota && apiKey.RemainingQuota <= 0 {
			c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{
				"error": gin.H{
					"message": "api key has insufficient quota",
					"type":    "insufficient_quota",
				},
			})
			return
		}

		// Save the verified record for gateway handlers.
		c.Set(APIKeyContextKey, apiKey)
		c.Next()
	}
}
