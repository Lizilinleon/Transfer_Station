package model

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"size:128;not null"`
	DisplayName  string    `json:"display_name" gorm:"size:128"`
	Role         string    `json:"role" gorm:"size:32;not null;default:user"`
	Status       string    `json:"status" gorm:"size:32;not null;default:active"`
	GroupName    string    `json:"group_name" gorm:"size:128;not null;default:default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AIModel struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"size:128;uniqueIndex;not null"`
	DisplayName     string    `json:"display_name" gorm:"size:128;not null"`
	Provider        string    `json:"provider" gorm:"size:128;not null"`
	GroupName       string    `json:"group_name" gorm:"size:128;not null"`
	BaseURL         string    `json:"base_url" gorm:"size:255"`
	Description     string    `json:"description" gorm:"size:500"`
	Enabled         bool      `json:"enabled" gorm:"not null;default:true"`
	IsDefault       bool      `json:"is_default" gorm:"not null;default:false"`
	SupportsVision  bool      `json:"supports_vision" gorm:"not null;default:false"`
	SupportsTools   bool      `json:"supports_tools" gorm:"not null;default:false"`
	SupportsStream  bool      `json:"supports_stream" gorm:"not null;default:true"`
	ContextWindow   int       `json:"context_window" gorm:"not null;default:0"`
	MaxOutputTokens int       `json:"max_output_tokens" gorm:"not null;default:0"`
	InputPrice      float64   `json:"input_price" gorm:"not null;default:0"`
	OutputPrice     float64   `json:"output_price" gorm:"not null;default:0"`
	RequestPrice    float64   `json:"request_price" gorm:"not null;default:0"`
	SortOrder       int       `json:"sort_order" gorm:"not null;default:0"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ProviderChannel struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"size:128;not null"`
	ProviderType      string    `json:"provider_type" gorm:"size:128;not null"`
	BaseURL           string    `json:"base_url" gorm:"size:255;not null"`
	APIKey            string    `json:"api_key" gorm:"size:255"`
	Organization      string    `json:"organization" gorm:"size:128"`
	GroupName         string    `json:"group_name" gorm:"size:128;not null"`
	ModelNames        []string  `json:"model_names" gorm:"serializer:json"`
	ModelMapping      string    `json:"model_mapping" gorm:"type:text"`
	ExtraHeaders      string    `json:"extra_headers" gorm:"type:text"`
	RequestTemplate   string    `json:"request_template" gorm:"type:text"`
	ResponseTemplate  string    `json:"response_template" gorm:"type:text"`
	Weight            int       `json:"weight" gorm:"not null;default:0"`
	Priority          int       `json:"priority" gorm:"not null;default:0"`
	Enabled           bool      `json:"enabled" gorm:"not null;default:true"`
	RateLimited       bool      `json:"rate_limited" gorm:"not null;default:false"`
	MaxRequestsMinute int       `json:"max_requests_minute" gorm:"not null;default:0"`
	TestModel         string    `json:"test_model" gorm:"size:128"`
	Remark            string    `json:"remark" gorm:"size:500"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ModelAbility struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	GroupName string           `json:"group_name" gorm:"size:128;index;not null"`
	ModelName string           `json:"model_name" gorm:"size:128;index;not null"`
	ChannelID uint             `json:"channel_id" gorm:"index;not null"`
	Channel   *ProviderChannel `json:"channel,omitempty"`
	Enabled   bool             `json:"enabled" gorm:"not null;default:true"`
	Priority  int              `json:"priority" gorm:"not null;default:0"`
	Weight    int              `json:"weight" gorm:"not null;default:0"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type APIKey struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	Name             string     `json:"name" gorm:"size:128;not null"`
	AccessKey        string     `json:"access_key" gorm:"size:128;uniqueIndex;not null"`
	GroupName        string     `json:"group_name" gorm:"size:128;not null"`
	ModelID          *uint      `json:"model_id"`
	Model            *AIModel   `json:"model,omitempty"`
	ModelNames       []string   `json:"model_names" gorm:"serializer:json"`
	Enabled          bool       `json:"enabled" gorm:"not null;default:true"`
	Quota            float64    `json:"quota" gorm:"not null;default:0"`
	UsedQuota        float64    `json:"used_quota" gorm:"not null;default:0"`
	RemainingQuota   float64    `json:"remaining_quota" gorm:"not null;default:0"`
	UnlimitedQuota   bool       `json:"unlimited_quota" gorm:"not null;default:false"`
	BillingRule      string     `json:"billing_rule" gorm:"size:128;not null;default:reserved"`
	BillingConfig    string     `json:"billing_config" gorm:"type:text"`
	InputTokenPrice  float64    `json:"input_token_price" gorm:"not null;default:0"`
	OutputTokenPrice float64    `json:"output_token_price" gorm:"not null;default:0"`
	RequestPrice     float64    `json:"request_price" gorm:"not null;default:0"`
	LastUsedAt       *time.Time `json:"last_used_at"`
	Remark           string     `json:"remark" gorm:"size:500"`
	ExpiresAt        *time.Time `json:"expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type UsageLog struct {
	ID               uint             `json:"id" gorm:"primaryKey"`
	APIKeyID         *uint            `json:"api_key_id" gorm:"index"`
	APIKey           *APIKey          `json:"api_key,omitempty"`
	ChannelID        *uint            `json:"channel_id" gorm:"index"`
	Channel          *ProviderChannel `json:"channel,omitempty"`
	ModelName        string           `json:"model_name" gorm:"size:128;index;not null"`
	RequestID        string           `json:"request_id" gorm:"size:128;index"`
	PromptTokens     int64            `json:"prompt_tokens" gorm:"not null;default:0"`
	CompletionTokens int64            `json:"completion_tokens" gorm:"not null;default:0"`
	TotalTokens      int64            `json:"total_tokens" gorm:"not null;default:0"`
	InputCost        float64          `json:"input_cost" gorm:"not null;default:0"`
	OutputCost       float64          `json:"output_cost" gorm:"not null;default:0"`
	RequestCost      float64          `json:"request_cost" gorm:"not null;default:0"`
	TotalCost        float64          `json:"total_cost" gorm:"not null;default:0"`
	Status           string           `json:"status" gorm:"size:32;not null;default:success"`
	ErrorMessage     string           `json:"error_message" gorm:"size:500"`
	LatencyMs        int              `json:"latency_ms" gorm:"not null;default:0"`
	ClientIP         string           `json:"client_ip" gorm:"size:64"`
	UserAgent        string           `json:"user_agent" gorm:"size:255"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type SystemOption struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OptionKey   string    `json:"option_key" gorm:"size:128;uniqueIndex;not null"`
	OptionValue string    `json:"option_value" gorm:"size:2000;not null"`
	Category    string    `json:"category" gorm:"size:64;not null"`
	Description string    `json:"description" gorm:"size:255"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ChatSession struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Title         string     `json:"title" gorm:"size:255;not null"`
	SystemPrompt  string     `json:"system_prompt" gorm:"size:2000"`
	UserID        *uint      `json:"user_id"`
	ModelName     string     `json:"model_name" gorm:"size:128"`
	ChannelID     *uint      `json:"channel_id"`
	Temperature   float64    `json:"temperature" gorm:"not null;default:1"`
	TopP          float64    `json:"top_p" gorm:"not null;default:1"`
	MaxTokens     int        `json:"max_tokens" gorm:"not null;default:0"`
	MessageCount  int        `json:"message_count" gorm:"not null;default:0"`
	LastMessageAt *time.Time `json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ChatMessage struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	SessionID      uint      `json:"session_id" gorm:"index;not null"`
	Role           string    `json:"role" gorm:"size:32;not null"`
	Content        string    `json:"content" gorm:"type:text;not null"`
	InputTokens    int64     `json:"input_tokens" gorm:"not null;default:0"`
	OutputTokens   int64     `json:"output_tokens" gorm:"not null;default:0"`
	FinishReason   string    `json:"finish_reason" gorm:"size:64"`
	ResponseTimeMs int       `json:"response_time_ms" gorm:"not null;default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
