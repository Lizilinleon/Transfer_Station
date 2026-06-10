package database

import (
	"fmt"

	"ai-chat-platform/backend/config"
	"ai-chat-platform/backend/model"
	"ai-chat-platform/backend/util"

	"gorm.io/gorm"
)

func seedDefaults(db *gorm.DB, cfg config.Config) error {
	admin := model.User{}
	if err := db.Where("username = ?", cfg.DefaultAdminUsername).First(&admin).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		admin = model.User{
			Username:     cfg.DefaultAdminUsername,
			PasswordHash: util.HashSHA256(cfg.DefaultAdminPassword),
			Role:         "admin",
			DisplayName:  "默认管理员",
			Status:       "active",
		}
		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("create default admin: %w", err)
		}
	}

	defaultModel := model.AIModel{}
	if err := db.Where("name = ?", cfg.DefaultModelName).First(&defaultModel).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		defaultModel = model.AIModel{
			Name:        cfg.DefaultModelName,
			DisplayName: cfg.DefaultModelName,
			Provider:    cfg.DefaultModelProvider,
			GroupName:   cfg.DefaultModelGroup,
			BaseURL:     "https://api.example.com/v1",
			Description: "默认模型，占位用于本地开发联调。",
			Enabled:     true,
			IsDefault:   true,
			SortOrder:   1,
		}
		if err := db.Create(&defaultModel).Error; err != nil {
			return fmt.Errorf("create default model: %w", err)
		}
	}

	options := []model.SystemOption{
		{
			OptionKey:   "site_name",
			OptionValue: cfg.AppName,
			Category:    "general",
			Description: "站点名称",
		},
		{
			OptionKey:   "default_model",
			OptionValue: cfg.DefaultModelName,
			Category:    "model",
			Description: "默认模型名称",
		},
		{
			OptionKey:   "default_group",
			OptionValue: cfg.DefaultModelGroup,
			Category:    "model",
			Description: "默认模型分组",
		},
	}

	for _, option := range options {
		current := model.SystemOption{}
		if err := db.Where("option_key = ?", option.OptionKey).First(&current).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			if err := db.Create(&option).Error; err != nil {
				return fmt.Errorf("create default option %s: %w", option.OptionKey, err)
			}
		}
	}

	return nil
}
