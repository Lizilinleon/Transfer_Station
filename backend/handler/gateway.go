package handler

import (
	"net/http"
	"time"

	"ai-chat-platform/backend/middleware"
	"ai-chat-platform/backend/model"
	"ai-chat-platform/backend/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GatewayHandler struct {
	DB             *gorm.DB
	DeepSeekClient *service.DeepSeekClient
}

type ChatCompletionRequest struct {
	Model       string                          `json:"model"`
	Messages    []service.ChatCompletionMessage `json:"messages"`
	Temperature *float64                        `json:"temperature,omitempty"`
	TopP        *float64                        `json:"top_p,omitempty"`
	MaxTokens   int                             `json:"max_tokens,omitempty"`
	Stream      bool                            `json:"stream,omitempty"`
}

func (h GatewayHandler) ListModels(c *gin.Context) {
	apiKey := c.MustGet(middleware.APIKeyContextKey).(model.APIKey)
	models, err := h.collectAvailableModels(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "failed to load models",
				"type":    "server_error",
			},
		})
		return
	}

	type modelItem struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		OwnedBy string `json:"owned_by"`
	}

	items := make([]modelItem, 0, len(models))
	for _, item := range models {
		items = append(items, modelItem{
			ID:      item.Name,
			Object:  "model",
			OwnedBy: item.Provider,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data":   items,
	})
}

func (h GatewayHandler) ChatCompletions(c *gin.Context) {
	apiKey := c.MustGet(middleware.APIKeyContextKey).(model.APIKey)

	var payload ChatCompletionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "invalid request body",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	if payload.Model == "" || len(payload.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "model and messages are required",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	if payload.Stream {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": gin.H{
				"message": "stream mode is reserved for future integration",
				"type":    "not_implemented",
			},
		})
		return
	}

	if !h.modelAllowed(apiKey, payload.Model) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"message": "model is not allowed for this api key",
				"type":    "permission_error",
			},
		})
		return
	}

	channel, err := h.resolveChannel(apiKey.GroupName, payload.Model)
	if err != nil {
		message := err.Error()
		if err == gorm.ErrRecordNotFound {
			message = "no enabled deepseek channel found for this model and group"
		}
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"message": message,
				"type":    "upstream_error",
			},
		})
		return
	}

	startedAt := time.Now()
	upstreamResponse, err := h.DeepSeekClient.ChatCompletions(c.Request.Context(), channel, service.ChatCompletionRequest{
		Model:       payload.Model,
		Messages:    payload.Messages,
		Temperature: payload.Temperature,
		TopP:        payload.TopP,
		MaxTokens:   payload.MaxTokens,
		Stream:      false,
	})
	latencyMs := int(time.Since(startedAt).Milliseconds())
	if err != nil {
		_ = h.writeUsageLog(apiKey, channel, payload.Model, service.ChatCompletionUsage{}, 0, "error", err.Error(), latencyMs, c)
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "upstream_error",
			},
		})
		return
	}

	totalCost := h.calculateCost(payload.Model, upstreamResponse.Usage)
	if !apiKey.UnlimitedQuota && apiKey.RemainingQuota < totalCost {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error": gin.H{
				"message": "insufficient quota for this request",
				"type":    "insufficient_quota",
			},
		})
		return
	}

	if err := h.applyUsage(apiKey, channel, payload.Model, upstreamResponse.Usage, totalCost, latencyMs, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "failed to persist usage information",
				"type":    "server_error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, upstreamResponse)
}

func (h GatewayHandler) collectAvailableModels(apiKey model.APIKey) ([]model.AIModel, error) {
	var items []model.AIModel
	query := h.DB.Where("enabled = ? AND group_name = ?", true, apiKey.GroupName).Order("sort_order asc, id asc")
	if len(apiKey.ModelNames) > 0 {
		query = query.Where("name IN ?", apiKey.ModelNames)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (h GatewayHandler) modelAllowed(apiKey model.APIKey, modelName string) bool {
	if len(apiKey.ModelNames) == 0 {
		return true
	}
	for _, item := range apiKey.ModelNames {
		if item == modelName {
			return true
		}
	}
	return false
}

func (h GatewayHandler) resolveChannel(groupName, modelName string) (model.ProviderChannel, error) {
	var abilities []model.ModelAbility
	if err := h.DB.Preload("Channel").
		Where("group_name = ? AND model_name = ? AND enabled = ?", groupName, modelName, true).
		Order("priority desc, weight desc, id asc").
		Find(&abilities).Error; err != nil {
		return model.ProviderChannel{}, err
	}

	for _, ability := range abilities {
		if ability.Channel != nil && ability.Channel.Enabled && ability.Channel.ProviderType == "deepseek" {
			return *ability.Channel, nil
		}
	}

	var channels []model.ProviderChannel
	if err := h.DB.
		Where("group_name = ? AND provider_type = ? AND enabled = ?", groupName, "deepseek", true).
		Order("priority desc, weight desc, id asc").
		Find(&channels).Error; err != nil {
		return model.ProviderChannel{}, err
	}

	for _, channel := range channels {
		if len(channel.ModelNames) == 0 {
			return channel, nil
		}
		for _, item := range channel.ModelNames {
			if item == modelName {
				return channel, nil
			}
		}
	}

	return model.ProviderChannel{}, gorm.ErrRecordNotFound
}

func (h GatewayHandler) calculateCost(modelName string, usage service.ChatCompletionUsage) float64 {
	var selectedModel model.AIModel
	if err := h.DB.Where("name = ?", modelName).First(&selectedModel).Error; err != nil {
		return 0
	}

	inputCost := float64(usage.PromptTokens) * selectedModel.InputPrice
	outputCost := float64(usage.CompletionTokens) * selectedModel.OutputPrice
	requestCost := selectedModel.RequestPrice
	return inputCost + outputCost + requestCost
}

func (h GatewayHandler) applyUsage(apiKey model.APIKey, channel model.ProviderChannel, modelName string, usage service.ChatCompletionUsage, totalCost float64, latencyMs int, c *gin.Context) error {
	tx := h.DB.Begin()
	now := time.Now()

	updates := map[string]any{
		"last_used_at": &now,
	}
	if !apiKey.UnlimitedQuota {
		updates["used_quota"] = gorm.Expr("used_quota + ?", totalCost)
		updates["remaining_quota"] = gorm.Expr("CASE WHEN remaining_quota - ? < 0 THEN 0 ELSE remaining_quota - ? END", totalCost, totalCost)
	}

	if err := tx.Model(&model.APIKey{}).Where("id = ?", apiKey.ID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&model.UsageLog{
		APIKeyID:         &apiKey.ID,
		ChannelID:        &channel.ID,
		ModelName:        modelName,
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
		TotalCost:        totalCost,
		Status:           "success",
		LatencyMs:        latencyMs,
		ClientIP:         c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (h GatewayHandler) writeUsageLog(apiKey model.APIKey, channel model.ProviderChannel, modelName string, usage service.ChatCompletionUsage, totalCost float64, status, errorMessage string, latencyMs int, c *gin.Context) error {
	channelID := channel.ID
	return h.DB.Create(&model.UsageLog{
		APIKeyID:         &apiKey.ID,
		ChannelID:        &channelID,
		ModelName:        modelName,
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
		TotalCost:        totalCost,
		Status:           status,
		ErrorMessage:     errorMessage,
		LatencyMs:        latencyMs,
		ClientIP:         c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
	}).Error
}
