package handler

import (
	"ai-chat-platform/backend/config"

	"github.com/gin-gonic/gin"
)

// HealthHandler exposes basic runtime information for probes.
type HealthHandler struct {
	Config config.Config
}

// Get returns a lightweight health response without touching the database.
func (h HealthHandler) Get(c *gin.Context) {
	success(c, gin.H{
		"app_name": h.Config.AppName,
		"env":      h.Config.AppEnv,
		"status":   "ok",
	})
}
