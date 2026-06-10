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
			DisplayName:  "Default Admin",
			Status:       "active",
			GroupName:    cfg.DefaultModelGroup,
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
			Name:            cfg.DefaultModelName,
			DisplayName:     cfg.DefaultModelName,
			Provider:        cfg.DefaultModelProvider,
			GroupName:       cfg.DefaultModelGroup,
			BaseURL:         cfg.DeepSeekBaseURL,
			Description:     "Default model prepared for DeepSeek gateway integration.",
			Enabled:         true,
			IsDefault:       true,
			SupportsStream:  true,
			ContextWindow:   64000,
			MaxOutputTokens: 8000,
			SortOrder:       1,
		}
		if err := db.Create(&defaultModel).Error; err != nil {
			return fmt.Errorf("create default model: %w", err)
		}
	}

	defaultChannel := model.ProviderChannel{}
	if err := db.Where("name = ?", cfg.DeepSeekChannelName).First(&defaultChannel).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		defaultChannel = model.ProviderChannel{
			Name:             cfg.DeepSeekChannelName,
			ProviderType:     "deepseek",
			BaseURL:          cfg.DeepSeekBaseURL,
			APIKey:           cfg.DeepSeekAPIKey,
			GroupName:        cfg.DefaultModelGroup,
			ModelNames:       []string{cfg.DefaultModelName},
			Enabled:          true,
			Priority:         100,
			Weight:           100,
			TestModel:        cfg.DefaultModelName,
			RequestTemplate:  "",
			ResponseTemplate: "",
			Remark:           "Default DeepSeek channel. Keep API key empty until real upstream credentials are ready.",
		}
		if err := db.Create(&defaultChannel).Error; err != nil {
			return fmt.Errorf("create default deepseek channel: %w", err)
		}
	}

	defaultAbility := model.ModelAbility{}
	if err := db.Where("group_name = ? AND model_name = ? AND channel_id = ?", cfg.DefaultModelGroup, cfg.DefaultModelName, defaultChannel.ID).First(&defaultAbility).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		defaultAbility = model.ModelAbility{
			GroupName: cfg.DefaultModelGroup,
			ModelName: cfg.DefaultModelName,
			ChannelID: defaultChannel.ID,
			Enabled:   true,
			Priority:  100,
			Weight:    100,
		}
		if err := db.Create(&defaultAbility).Error; err != nil {
			return fmt.Errorf("create default model ability: %w", err)
		}
	}

	options := []model.SystemOption{
		{
			OptionKey:   "site_name",
			OptionValue: cfg.AppName,
			Category:    "general",
			Description: "Site name",
		},
		{
			OptionKey:   "default_model",
			OptionValue: cfg.DefaultModelName,
			Category:    "model",
			Description: "Default model name",
		},
		{
			OptionKey:   "default_group",
			OptionValue: cfg.DefaultModelGroup,
			Category:    "model",
			Description: "Default model group",
		},
		{
			OptionKey:   "deepseek_base_url",
			OptionValue: cfg.DeepSeekBaseURL,
			Category:    "gateway",
			Description: "Default DeepSeek upstream base URL",
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
