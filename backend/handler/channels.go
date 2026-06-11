package handler

import (
	"encoding/json"
	"strconv"
	"time"

	"ai-chat-platform/backend/model"
	"ai-chat-platform/backend/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ChannelHandler manages upstream provider channel records.
type ChannelHandler struct {
	DB             *gorm.DB
	DeepSeekClient *service.DeepSeekClient
}

// channelPayload is the JSON shape accepted by channel create/update endpoints.
type channelPayload struct {
	Name              string   `json:"name"`
	ProviderType      string   `json:"provider_type"`
	BaseURL           string   `json:"base_url"`
	APIKey            string   `json:"api_key"`
	Organization      string   `json:"organization"`
	GroupName         string   `json:"group_name"`
	ModelNames        []string `json:"model_names"`
	ModelMapping      string   `json:"model_mapping"`
	ExtraHeaders      string   `json:"extra_headers"`
	RequestTemplate   string   `json:"request_template"`
	ResponseTemplate  string   `json:"response_template"`
	Weight            int      `json:"weight"`
	Priority          int      `json:"priority"`
	Enabled           bool     `json:"enabled"`
	RateLimited       bool     `json:"rate_limited"`
	MaxRequestsMinute int      `json:"max_requests_minute"`
	TestModel         string   `json:"test_model"`
	Remark            string   `json:"remark"`
}

type channelTestPayload struct {
	Model       string   `json:"model"`
	Messages    []string `json:"messages"`
	Temperature *float64 `json:"temperature"`
	MaxTokens   int      `json:"max_tokens"`
}

// List returns provider channels with optional provider, group, and enabled filters.
func (h ChannelHandler) List(c *gin.Context) {
	var items []model.ProviderChannel
	query := h.DB.Order("priority desc, weight desc, id asc")

	if providerType := c.Query("provider_type"); providerType != "" {
		query = query.Where("provider_type = ?", providerType)
	}
	if groupName := c.Query("group_name"); groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}
	if enabled := c.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	if err := query.Find(&items).Error; err != nil {
		fail(c, 500, "failed to list channels")
		return
	}

	success(c, gin.H{"items": items})
}

// Get loads one provider channel by path id.
func (h ChannelHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid channel id")
		return
	}

	var item model.ProviderChannel
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "channel not found")
			return
		}
		fail(c, 500, "failed to get channel")
		return
	}

	success(c, item)
}

// Create validates required upstream settings and inserts a channel.
func (h ChannelHandler) Create(c *gin.Context) {
	var payload channelPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	if payload.Name == "" || payload.ProviderType == "" || payload.BaseURL == "" || payload.GroupName == "" {
		fail(c, 400, "name, provider_type, base_url and group_name are required")
		return
	}

	item := model.ProviderChannel{
		Name:              payload.Name,
		ProviderType:      payload.ProviderType,
		BaseURL:           payload.BaseURL,
		APIKey:            payload.APIKey,
		Organization:      payload.Organization,
		GroupName:         payload.GroupName,
		ModelNames:        payload.ModelNames,
		ModelMapping:      payload.ModelMapping,
		ExtraHeaders:      payload.ExtraHeaders,
		RequestTemplate:   payload.RequestTemplate,
		ResponseTemplate:  payload.ResponseTemplate,
		Weight:            payload.Weight,
		Priority:          payload.Priority,
		Enabled:           payload.Enabled,
		RateLimited:       payload.RateLimited,
		MaxRequestsMinute: payload.MaxRequestsMinute,
		TestModel:         payload.TestModel,
		Remark:            payload.Remark,
	}

	if err := h.DB.Create(&item).Error; err != nil {
		fail(c, 500, "failed to create channel")
		return
	}

	success(c, item)
}

// Update replaces editable settings for an existing provider channel.
func (h ChannelHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid channel id")
		return
	}

	var payload channelPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	var item model.ProviderChannel
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "channel not found")
			return
		}
		fail(c, 500, "failed to load channel")
		return
	}

	item.Name = payload.Name
	item.ProviderType = payload.ProviderType
	item.BaseURL = payload.BaseURL
	item.APIKey = payload.APIKey
	item.Organization = payload.Organization
	item.GroupName = payload.GroupName
	item.ModelNames = payload.ModelNames
	item.ModelMapping = payload.ModelMapping
	item.ExtraHeaders = payload.ExtraHeaders
	item.RequestTemplate = payload.RequestTemplate
	item.ResponseTemplate = payload.ResponseTemplate
	item.Weight = payload.Weight
	item.Priority = payload.Priority
	item.Enabled = payload.Enabled
	item.RateLimited = payload.RateLimited
	item.MaxRequestsMinute = payload.MaxRequestsMinute
	item.TestModel = payload.TestModel
	item.Remark = payload.Remark

	if err := h.DB.Save(&item).Error; err != nil {
		fail(c, 500, "failed to update channel")
		return
	}

	success(c, item)
}

// Delete removes one provider channel by id.
func (h ChannelHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid channel id")
		return
	}

	if err := h.DB.Delete(&model.ProviderChannel{}, id).Error; err != nil {
		fail(c, 500, "failed to delete channel")
		return
	}

	success(c, gin.H{"deleted": true})
}

// Test sends a lightweight chat request through the selected channel to verify upstream connectivity.
func (h ChannelHandler) Test(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid channel id")
		return
	}

	if h.DeepSeekClient == nil {
		fail(c, 500, "channel test client is not configured")
		return
	}

	var item model.ProviderChannel
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "channel not found")
			return
		}
		fail(c, 500, "failed to load channel")
		return
	}

	if item.ProviderType != "deepseek" {
		fail(c, 400, "channel test currently supports deepseek only")
		return
	}

	var payload channelTestPayload
	if err := c.ShouldBindJSON(&payload); err != nil && err.Error() != "EOF" {
		fail(c, 400, "invalid request body")
		return
	}

	modelName := payload.Model
	if modelName == "" {
		modelName = item.TestModel
	}
	if modelName == "" && len(item.ModelNames) > 0 {
		modelName = item.ModelNames[0]
	}
	if modelName == "" {
		fail(c, 400, "channel test model is required")
		return
	}

	rawMessages, err := buildChannelTestMessages(payload.Messages)
	if err != nil {
		fail(c, 400, "invalid test messages")
		return
	}
	if len(rawMessages) == 0 {
		rawMessages = []service.ChatCompletionMessage{
			{
				Role:    "user",
				Content: json.RawMessage(`"Ping from channel test. Reply briefly with OK."`),
			},
		}
	}

	startedAt := time.Now()
	response, err := h.DeepSeekClient.ChatCompletions(c.Request.Context(), item, service.ChatCompletionRequest{
		Model:       modelName,
		Messages:    rawMessages,
		Temperature: payload.Temperature,
		MaxTokens:   payload.MaxTokens,
		Stream:      false,
	})
	latencyMs := int(time.Since(startedAt).Milliseconds())
	if err != nil {
		fail(c, 502, err.Error())
		return
	}

	success(c, gin.H{
		"channel_id":    item.ID,
		"channel_name":  item.Name,
		"provider_type": item.ProviderType,
		"model":         modelName,
		"latency_ms":    latencyMs,
		"response":      response,
	})
}

func buildChannelTestMessages(messages []string) ([]service.ChatCompletionMessage, error) {
	items := make([]service.ChatCompletionMessage, 0, len(messages))
	for _, item := range messages {
		if item == "" {
			continue
		}
		content, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		items = append(items, service.ChatCompletionMessage{
			Role:    "user",
			Content: content,
		})
	}
	return items, nil
}
