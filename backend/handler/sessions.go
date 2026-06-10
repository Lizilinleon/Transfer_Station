package handler

import (
	"strconv"
	"time"

	"ai-chat-platform/backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SessionHandler struct {
	DB *gorm.DB
}

type sessionPayload struct {
	Title        string  `json:"title"`
	SystemPrompt string  `json:"system_prompt"`
	UserID       *uint   `json:"user_id"`
	ModelName    string  `json:"model_name"`
	ChannelID    *uint   `json:"channel_id"`
	Temperature  float64 `json:"temperature"`
	TopP         float64 `json:"top_p"`
	MaxTokens    int     `json:"max_tokens"`
}

func (h SessionHandler) List(c *gin.Context) {
	var items []model.ChatSession
	query := h.DB.Order("updated_at desc, id desc")

	if modelName := c.Query("model_name"); modelName != "" {
		query = query.Where("model_name = ?", modelName)
	}

	if err := query.Find(&items).Error; err != nil {
		fail(c, 500, "failed to list sessions")
		return
	}

	success(c, gin.H{"items": items})
}

func (h SessionHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid session id")
		return
	}

	var item model.ChatSession
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "session not found")
			return
		}
		fail(c, 500, "failed to get session")
		return
	}

	success(c, item)
}

func (h SessionHandler) Create(c *gin.Context) {
	var payload sessionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	if payload.Title == "" {
		fail(c, 400, "title is required")
		return
	}

	item := model.ChatSession{
		Title:        payload.Title,
		SystemPrompt: payload.SystemPrompt,
		UserID:       payload.UserID,
		ModelName:    payload.ModelName,
		ChannelID:    payload.ChannelID,
		Temperature:  payload.Temperature,
		TopP:         payload.TopP,
		MaxTokens:    payload.MaxTokens,
	}

	if err := h.DB.Create(&item).Error; err != nil {
		fail(c, 500, "failed to create session")
		return
	}

	success(c, item)
}

func (h SessionHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid session id")
		return
	}

	var payload sessionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	var item model.ChatSession
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "session not found")
			return
		}
		fail(c, 500, "failed to load session")
		return
	}

	item.Title = payload.Title
	item.SystemPrompt = payload.SystemPrompt
	item.UserID = payload.UserID
	item.ModelName = payload.ModelName
	item.ChannelID = payload.ChannelID
	item.Temperature = payload.Temperature
	item.TopP = payload.TopP
	item.MaxTokens = payload.MaxTokens

	if err := h.DB.Save(&item).Error; err != nil {
		fail(c, 500, "failed to update session")
		return
	}

	success(c, item)
}

func (h SessionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid session id")
		return
	}

	tx := h.DB.Begin()
	if err := tx.Where("session_id = ?", id).Delete(&model.ChatMessage{}).Error; err != nil {
		tx.Rollback()
		fail(c, 500, "failed to delete session messages")
		return
	}
	if err := tx.Delete(&model.ChatSession{}, id).Error; err != nil {
		tx.Rollback()
		fail(c, 500, "failed to delete session")
		return
	}
	tx.Commit()

	success(c, gin.H{"deleted": true})
}

type MessageHandler struct {
	DB *gorm.DB
}

type messagePayload struct {
	Role           string `json:"role"`
	Content        string `json:"content"`
	InputTokens    int64  `json:"input_tokens"`
	OutputTokens   int64  `json:"output_tokens"`
	FinishReason   string `json:"finish_reason"`
	ResponseTimeMs int    `json:"response_time_ms"`
}

func (h MessageHandler) List(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid session id")
		return
	}

	var items []model.ChatMessage
	if err := h.DB.Where("session_id = ?", sessionID).Order("id asc").Find(&items).Error; err != nil {
		fail(c, 500, "failed to list messages")
		return
	}

	success(c, gin.H{"items": items})
}

func (h MessageHandler) Create(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid session id")
		return
	}

	var payload messagePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	if payload.Role == "" || payload.Content == "" {
		fail(c, 400, "role and content are required")
		return
	}

	item := model.ChatMessage{
		SessionID:      uint(sessionID),
		Role:           payload.Role,
		Content:        payload.Content,
		InputTokens:    payload.InputTokens,
		OutputTokens:   payload.OutputTokens,
		FinishReason:   payload.FinishReason,
		ResponseTimeMs: payload.ResponseTimeMs,
	}

	if err := h.DB.Create(&item).Error; err != nil {
		fail(c, 500, "failed to create message")
		return
	}

	now := time.Now()
	if err := h.DB.Model(&model.ChatSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]any{
			"message_count":   gorm.Expr("message_count + 1"),
			"last_message_at": &now,
		}).Error; err != nil {
		fail(c, 500, "failed to update session metadata")
		return
	}

	success(c, item)
}

func (h MessageHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid message id")
		return
	}

	if err := h.DB.Delete(&model.ChatMessage{}, id).Error; err != nil {
		fail(c, 500, "failed to delete message")
		return
	}

	success(c, gin.H{"deleted": true})
}
