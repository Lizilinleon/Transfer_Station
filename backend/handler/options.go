package handler

import (
	"strconv"

	"ai-chat-platform/backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OptionHandler manages editable system key-value options.
type OptionHandler struct {
	DB *gorm.DB
}

// optionPayload is the JSON shape accepted by option create/update endpoints.
type optionPayload struct {
	OptionKey   string `json:"option_key"`
	OptionValue string `json:"option_value"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// List returns options, optionally filtered by category.
func (h OptionHandler) List(c *gin.Context) {
	var items []model.SystemOption
	query := h.DB.Order("category asc, id asc")

	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&items).Error; err != nil {
		fail(c, 500, "failed to list options")
		return
	}

	success(c, gin.H{"items": items})
}

// Create validates and inserts a new system option.
func (h OptionHandler) Create(c *gin.Context) {
	var payload optionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	if payload.OptionKey == "" || payload.OptionValue == "" || payload.Category == "" {
		fail(c, 400, "option_key, option_value and category are required")
		return
	}

	item := model.SystemOption{
		OptionKey:   payload.OptionKey,
		OptionValue: payload.OptionValue,
		Category:    payload.Category,
		Description: payload.Description,
	}

	if err := h.DB.Create(&item).Error; err != nil {
		fail(c, 500, "failed to create option")
		return
	}

	success(c, item)
}

// Update replaces editable fields for an existing system option.
func (h OptionHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid option id")
		return
	}

	var payload optionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	var item model.SystemOption
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "option not found")
			return
		}
		fail(c, 500, "failed to load option")
		return
	}

	item.OptionKey = payload.OptionKey
	item.OptionValue = payload.OptionValue
	item.Category = payload.Category
	item.Description = payload.Description

	if err := h.DB.Save(&item).Error; err != nil {
		fail(c, 500, "failed to update option")
		return
	}

	success(c, item)
}

// Delete removes one system option by id.
func (h OptionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid option id")
		return
	}

	if err := h.DB.Delete(&model.SystemOption{}, id).Error; err != nil {
		fail(c, 500, "failed to delete option")
		return
	}

	success(c, gin.H{"deleted": true})
}
