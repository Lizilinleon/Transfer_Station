package handler

import (
	"ai-chat-platform/backend/config"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	Config config.Config
}

func (h HealthHandler) Get(c *gin.Context) {
	success(c, gin.H{
		"app_name": h.Config.AppName,
		"env":      h.Config.AppEnv,
		"status":   "ok",
	})
}
