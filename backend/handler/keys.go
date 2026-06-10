package handler

import (
	"strconv"
	"time"

	"ai-chat-platform/backend/model"
	"ai-chat-platform/backend/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KeyHandler struct {
	DB *gorm.DB
}

type keyPayload struct {
	Name             string   `json:"name"`
	GroupName        string   `json:"group_name"`
	ModelID          *uint    `json:"model_id"`
	ModelNames       []string `json:"model_names"`
	Enabled          *bool    `json:"enabled"`
	Quota            float64  `json:"quota"`
	UsedQuota        float64  `json:"used_quota"`
	UnlimitedQuota   bool     `json:"unlimited_quota"`
	BillingRule      string   `json:"billing_rule"`
	BillingConfig    string   `json:"billing_config"`
	InputTokenPrice  float64  `json:"input_token_price"`
	OutputTokenPrice float64  `json:"output_token_price"`
	RequestPrice     float64  `json:"request_price"`
	Remark           string   `json:"remark"`
	ExpiresAt        string   `json:"expires_at"`
}

func (h KeyHandler) List(c *gin.Context) {
	var items []model.APIKey
	query := h.DB.Preload("Model").Order("id desc")
	keyword := c.Query("keyword")

	if enabled := c.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	if groupName := c.Query("group_name"); groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}

	if modelID := c.Query("model_id"); modelID != "" {
		query = query.Where("model_id = ?", modelID)
	}

	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR access_key LIKE ? OR group_name LIKE ?", likeKeyword, likeKeyword, likeKeyword)
	}

	if err := query.Find(&items).Error; err != nil {
		fail(c, 500, "failed to list keys")
		return
	}

	success(c, gin.H{"items": items})
}

func (h KeyHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid key id")
		return
	}

	var item model.APIKey
	if err := h.DB.Preload("Model").First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "key not found")
			return
		}
		fail(c, 500, "failed to get key")
		return
	}

	success(c, item)
}

func (h KeyHandler) Create(c *gin.Context) {
	var payload keyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	if payload.Name == "" || payload.GroupName == "" {
		fail(c, 400, "name and group_name are required")
		return
	}

	item := model.APIKey{
		Name:             payload.Name,
		AccessKey:        util.GenerateAccessKey(),
		GroupName:        payload.GroupName,
		ModelID:          payload.ModelID,
		ModelNames:       payload.ModelNames,
		Enabled:          normalizeEnabled(payload.Enabled),
		Quota:            payload.Quota,
		UsedQuota:        payload.UsedQuota,
		RemainingQuota:   calculateRemainingQuota(payload.Quota, payload.UsedQuota, payload.UnlimitedQuota),
		UnlimitedQuota:   payload.UnlimitedQuota,
		BillingRule:      normalizeBillingRule(payload.BillingRule),
		BillingConfig:    payload.BillingConfig,
		InputTokenPrice:  payload.InputTokenPrice,
		OutputTokenPrice: payload.OutputTokenPrice,
		RequestPrice:     payload.RequestPrice,
		Remark:           payload.Remark,
	}

	if payload.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, payload.ExpiresAt)
		if err != nil {
			fail(c, 400, "expires_at must be RFC3339 format")
			return
		}
		item.ExpiresAt = &expiresAt
	}

	if err := h.DB.Create(&item).Error; err != nil {
		fail(c, 500, "failed to create key")
		return
	}

	if err := h.DB.Preload("Model").First(&item, item.ID).Error; err != nil {
		fail(c, 500, "failed to load created key")
		return
	}

	success(c, item)
}

func (h KeyHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid key id")
		return
	}

	var payload keyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, 400, "invalid request body")
		return
	}

	var item model.APIKey
	if err := h.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, 404, "key not found")
			return
		}
		fail(c, 500, "failed to load key")
		return
	}

	item.Name = payload.Name
	item.GroupName = payload.GroupName
	item.ModelID = payload.ModelID
	item.ModelNames = payload.ModelNames
	item.Enabled = normalizeEnabled(payload.Enabled)
	item.Quota = payload.Quota
	item.UsedQuota = payload.UsedQuota
	item.RemainingQuota = calculateRemainingQuota(payload.Quota, payload.UsedQuota, payload.UnlimitedQuota)
	item.UnlimitedQuota = payload.UnlimitedQuota
	item.BillingRule = normalizeBillingRule(payload.BillingRule)
	item.BillingConfig = payload.BillingConfig
	item.InputTokenPrice = payload.InputTokenPrice
	item.OutputTokenPrice = payload.OutputTokenPrice
	item.RequestPrice = payload.RequestPrice
	item.Remark = payload.Remark
	item.ExpiresAt = nil

	if payload.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, payload.ExpiresAt)
		if err != nil {
			fail(c, 400, "expires_at must be RFC3339 format")
			return
		}
		item.ExpiresAt = &expiresAt
	}

	if err := h.DB.Save(&item).Error; err != nil {
		fail(c, 500, "failed to update key")
		return
	}

	if err := h.DB.Preload("Model").First(&item, item.ID).Error; err != nil {
		fail(c, 500, "failed to load updated key")
		return
	}

	success(c, item)
}

func (h KeyHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fail(c, 400, "invalid key id")
		return
	}

	if err := h.DB.Delete(&model.APIKey{}, id).Error; err != nil {
		fail(c, 500, "failed to delete key")
		return
	}

	success(c, gin.H{"deleted": true})
}

func normalizeEnabled(input *bool) bool {
	if input == nil {
		return true
	}
	return *input
}

func normalizeBillingRule(input string) string {
	if input == "" {
		return "reserved"
	}
	return input
}

func calculateRemainingQuota(quota, usedQuota float64, unlimited bool) float64 {
	if unlimited {
		return 0
	}
	remaining := quota - usedQuota
	if remaining < 0 {
		return 0
	}
	return remaining
}
