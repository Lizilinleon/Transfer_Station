package handler

import (
	"strconv"

	"ai-chat-platform/backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModelHandler struct {
	DB *gorm.DB
}

type modelPayload struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Provider    string `json:"provider"`
	GroupName   string `json:"group_name"`
	BaseURL     string `json:"base_url"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	IsDefault   bool   `json:"is_default"`
	SortOrder   int    `json:"sort_order"`
}

func (h ModelHandler) List(c *gin.Context) {
	var models []model.AIModel
	query := h.DB.Order("sort_order asc, id asc")

	if enabled := c.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	if groupName := c.Query("group_name"); groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}

	if err := query.Find(&models).Error; err != nil {
		fail(c, 500, "failed to list models")
		return
	}

	success(c, gin.H{"items": models})
}

func (h ModelHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid model id")
		return
	}

	var item model.AIModel
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "model not found")
			return
		}
		fail(c, 500, "failed to get model")
		return
	}

	success(c, item)
}

func (h ModelHandler) Create(c *gin.Context) {
	var payload modelPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	if payload.Name == "" || payload.DisplayName == "" || payload.Provider == "" || payload.GroupName == "" {
		fail(c, 400, "name, display_name, provider and group_name are required")
		return
	}

	item := model.AIModel{
		Name:        payload.Name,
		DisplayName: payload.DisplayName,
		Provider:    payload.Provider,
		GroupName:   payload.GroupName,
		BaseURL:     payload.BaseURL,
		Description: payload.Description,
		Enabled:     payload.Enabled,
		IsDefault:   payload.IsDefault,
		SortOrder:   payload.SortOrder,
	}

	if err := h.DB.Create(&item).Error; err != nil {
		fail(c, 500, "failed to create model")
		return
	}

	success(c, item)
}

func (h ModelHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid model id")
		return
	}

	var payload modelPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	var item model.AIModel
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "model not found")
			return
		}
		fail(c, 500, "failed to load model")
		return
	}

	item.Name = payload.Name
	item.DisplayName = payload.DisplayName
	item.Provider = payload.Provider
	item.GroupName = payload.GroupName
	item.BaseURL = payload.BaseURL
	item.Description = payload.Description
	item.Enabled = payload.Enabled
	item.IsDefault = payload.IsDefault
	item.SortOrder = payload.SortOrder

	if err := h.DB.Save(&item).Error; err != nil {
		fail(c, 500, "failed to update model")
		return
	}

	success(c, item)
}

func (h ModelHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid model id")
		return
	}

	if err := h.DB.Delete(&model.AIModel{}, id).Error; err != nil {
		fail(c, 500, "failed to delete model")
		return
	}

	success(c, gin.H{"deleted": true})
}
