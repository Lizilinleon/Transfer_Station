package handler

import (
	"strconv"

	"ai-chat-platform/backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UsageLogHandler exposes usage records for admin inspection and manual entry.
type UsageLogHandler struct {
	DB *gorm.DB
}

// usageLogPayload is the JSON shape accepted by manual usage-log creation.
type usageLogPayload struct {
	APIKeyID         *uint   `json:"api_key_id"`
	ChannelID        *uint   `json:"channel_id"`
	ModelName        string  `json:"model_name"`
	RequestID        string  `json:"request_id"`
	PromptTokens     int64   `json:"prompt_tokens"`
	CompletionTokens int64   `json:"completion_tokens"`
	TotalTokens      int64   `json:"total_tokens"`
	InputCost        float64 `json:"input_cost"`
	OutputCost       float64 `json:"output_cost"`
	RequestCost      float64 `json:"request_cost"`
	TotalCost        float64 `json:"total_cost"`
	Status           string  `json:"status"`
	ErrorMessage     string  `json:"error_message"`
	LatencyMs        int     `json:"latency_ms"`
	ClientIP         string  `json:"client_ip"`
	UserAgent        string  `json:"user_agent"`
}

// List returns usage logs with optional model, status, and API key filters.
func (h UsageLogHandler) List(c *gin.Context) {
	var items []model.UsageLog
	query := h.DB.Preload("APIKey").Preload("Channel").Order("created_at desc, id desc")

	if modelName := c.Query("model_name"); modelName != "" {
		query = query.Where("model_name = ?", modelName)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if apiKeyID := c.Query("api_key_id"); apiKeyID != "" {
		query = query.Where("api_key_id = ?", apiKeyID)
	}

	if err := query.Find(&items).Error; err != nil {
		fail(c, 500, "failed to list usage logs")
		return
	}

	success(c, gin.H{"items": items})
}

// Create inserts a usage log record, defaulting status to success when omitted.
func (h UsageLogHandler) Create(c *gin.Context) {
	var payload usageLogPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	if payload.ModelName == "" {
		fail(c, 400, "model_name is required")
		return
	}

	item := model.UsageLog{
		APIKeyID:         payload.APIKeyID,
		ChannelID:        payload.ChannelID,
		ModelName:        payload.ModelName,
		RequestID:        payload.RequestID,
		PromptTokens:     payload.PromptTokens,
		CompletionTokens: payload.CompletionTokens,
		TotalTokens:      payload.TotalTokens,
		InputCost:        payload.InputCost,
		OutputCost:       payload.OutputCost,
		RequestCost:      payload.RequestCost,
		TotalCost:        payload.TotalCost,
		Status:           payload.Status,
		ErrorMessage:     payload.ErrorMessage,
		LatencyMs:        payload.LatencyMs,
		ClientIP:         payload.ClientIP,
		UserAgent:        payload.UserAgent,
	}

	if item.Status == "" {
		item.Status = "success"
	}

	if err := h.DB.Create(&item).Error; err != nil {
		fail(c, 500, "failed to create usage log")
		return
	}

	if err := h.DB.Preload("APIKey").Preload("Channel").First(&item, item.ID).Error; err != nil {
		fail(c, 500, "failed to load created usage log")
		return
	}

	success(c, item)
}

// Delete removes one usage log by id.
func (h UsageLogHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid usage log id")
		return
	}

	if err := h.DB.Delete(&model.UsageLog{}, id).Error; err != nil {
		fail(c, 500, "failed to delete usage log")
		return
	}

	success(c, gin.H{"deleted": true})
}
