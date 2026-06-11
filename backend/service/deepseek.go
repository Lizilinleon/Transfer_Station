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
	Role             string          `json:"role"`
	Content          json.RawMessage `json:"content,omitempty"`
	Name             string          `json:"name,omitempty"`
	ToolCallID       string          `json:"tool_call_id,omitempty"`
	ToolCalls        json.RawMessage `json:"tool_calls,omitempty"`
	ReasoningContent string          `json:"reasoning_content,omitempty"`
}

// ChatCompletionRequest mirrors the non-streaming chat completion payload.
type ChatCompletionRequest struct {
	Model               string                  `json:"model"`
	Messages            []ChatCompletionMessage `json:"messages"`
	Temperature         *float64                `json:"temperature,omitempty"`
	TopP                *float64                `json:"top_p,omitempty"`
	MaxTokens           int                     `json:"max_tokens,omitempty"`
	MaxCompletionTokens int                     `json:"max_completion_tokens,omitempty"`
	PresencePenalty     *float64                `json:"presence_penalty,omitempty"`
	FrequencyPenalty    *float64                `json:"frequency_penalty,omitempty"`
	Stop                json.RawMessage         `json:"stop,omitempty"`
	ResponseFormat      json.RawMessage         `json:"response_format,omitempty"`
	Tools               json.RawMessage         `json:"tools,omitempty"`
	ToolChoice          json.RawMessage         `json:"tool_choice,omitempty"`
	Seed                *int                    `json:"seed,omitempty"`
	User                string                  `json:"user,omitempty"`
	Logprobs            *bool                   `json:"logprobs,omitempty"`
	TopLogprobs         *int                    `json:"top_logprobs,omitempty"`
	Stream              bool                    `json:"stream,omitempty"`
}

// ChatCompletionChoice represents one generated answer from the upstream model.
type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
	Logprobs     json.RawMessage       `json:"logprobs,omitempty"`
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

	endpoint := buildChatCompletionsEndpoint(channel.BaseURL)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("create deepseek request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+channel.APIKey)
	request.Header.Set("Content-Type", "application/json")
	if channel.Organization != "" {
		request.Header.Set("OpenAI-Organization", channel.Organization)
	}
	if err := applyExtraHeaders(request, channel.ExtraHeaders); err != nil {
		return ChatCompletionResponse{}, err
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

func buildChatCompletionsEndpoint(baseURL string) string {
	normalized := strings.TrimRight(baseURL, "/")
	if strings.HasSuffix(normalized, "/chat/completions") {
		return normalized
	}
	return normalized + "/chat/completions"
}

func applyExtraHeaders(request *http.Request, rawHeaders string) error {
	rawHeaders = strings.TrimSpace(rawHeaders)
	if rawHeaders == "" {
		return nil
	}

	var headers map[string]string
	if err := json.Unmarshal([]byte(rawHeaders), &headers); err != nil {
		return fmt.Errorf("parse channel extra_headers: %w", err)
	}
	for key, value := range headers {
		if strings.TrimSpace(key) == "" {
			continue
		}
		request.Header.Set(key, value)
	}
	return nil
}
