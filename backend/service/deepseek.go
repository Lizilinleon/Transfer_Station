package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"ai-chat-platform/backend/model"
)

// DeepSeekClient wraps HTTP calls to the DeepSeek-compatible upstream API.
type DeepSeekClient struct {
	HTTPClient *http.Client
}

// ChatCompletionMessage is the role/content pair used by chat requests.
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest mirrors the non-streaming chat completion payload.
type ChatCompletionRequest struct {
	Model       string                  `json:"model"`
	Messages    []ChatCompletionMessage `json:"messages"`
	Temperature *float64                `json:"temperature,omitempty"`
	TopP        *float64                `json:"top_p,omitempty"`
	MaxTokens   int                     `json:"max_tokens,omitempty"`
	Stream      bool                    `json:"stream,omitempty"`
}

// ChatCompletionChoice represents one generated answer from the upstream model.
type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

// ChatCompletionUsage carries token counts returned by the upstream service.
type ChatCompletionUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

// ChatCompletionResponse is the upstream response shape returned to clients.
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   ChatCompletionUsage    `json:"usage"`
}

// NewDeepSeekClient creates an HTTP client with a configurable request timeout.
func NewDeepSeekClient(timeoutSeconds int) *DeepSeekClient {
	return &DeepSeekClient{
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

// ChatCompletions sends one non-streaming request through the selected channel.
func (client *DeepSeekClient) ChatCompletions(ctx context.Context, channel model.ProviderChannel, payload ChatCompletionRequest) (ChatCompletionResponse, error) {
	if channel.APIKey == "" {
		return ChatCompletionResponse{}, fmt.Errorf("upstream deepseek api key is empty")
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("marshal deepseek payload: %w", err)
	}

	// Normalize the base URL so channel configuration can include a trailing slash.
	endpoint := strings.TrimRight(channel.BaseURL, "/") + "/chat/completions"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("create deepseek request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+channel.APIKey)
	request.Header.Set("Content-Type", "application/json")
	if channel.Organization != "" {
		request.Header.Set("OpenAI-Organization", channel.Organization)
	}

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("send deepseek request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("read deepseek response: %w", err)
	}

	if response.StatusCode >= 400 {
		return ChatCompletionResponse{}, fmt.Errorf("deepseek upstream error: %s", strings.TrimSpace(string(body)))
	}

	var parsed ChatCompletionResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("decode deepseek response: %w", err)
	}

	return parsed, nil
}
