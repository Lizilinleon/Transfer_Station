package handler

import (
	"strconv"

	"ai-chat-platform/backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AbilityHandler struct {
	DB *gorm.DB
}

type abilityPayload struct {
	GroupName string `json:"group_name"`
	ModelName string `json:"model_name"`
	ChannelID uint   `json:"channel_id"`
	Enabled   bool   `json:"enabled"`
	Priority  int    `json:"priority"`
	Weight    int    `json:"weight"`
}

func (h AbilityHandler) List(c *gin.Context) {
	var items []model.ModelAbility
	query := h.DB.Preload("Channel").Order("priority desc, weight desc, id asc")

	if groupName := c.Query("group_name"); groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}
	if modelName := c.Query("model_name"); modelName != "" {
		query = query.Where("model_name = ?", modelName)
	}
	if enabled := c.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	if err := query.Find(&items).Error; err != nil {
		fail(c, 500, "failed to list abilities")
		return
	}

	success(c, gin.H{"items": items})
}

func (h AbilityHandler) Create(c *gin.Context) {
	var payload abilityPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	if payload.GroupName == "" || payload.ModelName == "" || payload.ChannelID == 0 {
		fail(c, 400, "group_name, model_name and channel_id are required")
		return
	}

	item := model.ModelAbility{
		GroupName: payload.GroupName,
		ModelName: payload.ModelName,
		ChannelID: payload.ChannelID,
		Enabled:   payload.Enabled,
		Priority:  payload.Priority,
		Weight:    payload.Weight,
	}

	if err := h.DB.Create(&item).Error; err != nil {
		fail(c, 500, "failed to create ability")
		return
	}

	if err := h.DB.Preload("Channel").First(&item, item.ID).Error; err != nil {
		fail(c, 500, "failed to load created ability")
		return
	}

	success(c, item)
}

func (h AbilityHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid ability id")
		return
	}

	var payload abilityPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	var item model.ModelAbility
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "ability not found")
			return
		}
		fail(c, 500, "failed to load ability")
		return
	}

	item.GroupName = payload.GroupName
	item.ModelName = payload.ModelName
	item.ChannelID = payload.ChannelID
	item.Enabled = payload.Enabled
	item.Priority = payload.Priority
	item.Weight = payload.Weight

	if err := h.DB.Save(&item).Error; err != nil {
		fail(c, 500, "failed to update ability")
		return
	}

	if err := h.DB.Preload("Channel").First(&item, item.ID).Error; err != nil {
		fail(c, 500, "failed to load updated ability")
		return
	}

	success(c, item)
}

func (h AbilityHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid ability id")
		return
	}

	if err := h.DB.Delete(&model.ModelAbility{}, id).Error; err != nil {
		fail(c, 500, "failed to delete ability")
		return
	}

	success(c, gin.H{"deleted": true})
}
