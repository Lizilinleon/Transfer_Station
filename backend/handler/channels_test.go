package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ai-chat-platform/backend/model"
	"ai-chat-platform/backend/service"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestBuildChannelTestMessages(t *testing.T) {
	items, err := buildChannelTestMessages([]string{"hello", "world"})
	if err != nil {
		t.Fatalf("buildChannelTestMessages() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("buildChannelTestMessages() len = %d, want 2", len(items))
	}
	if string(items[0].Content) != `"hello"` {
		t.Fatalf("first message content = %s, want %q", string(items[0].Content), `"hello"`)
	}
}

func TestChannelHandlerTest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("request path = %s, want /chat/completions", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer ds-test-key" {
			t.Fatalf("authorization header = %q, want Bearer ds-test-key", auth)
		}
		if trace := r.Header.Get("X-Test-Trace"); trace != "channel-test" {
			t.Fatalf("X-Test-Trace header = %q, want channel-test", trace)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"chatcmpl-test",
			"object":"chat.completion",
			"created":1710000000,
			"model":"deepseek-chat",
			"choices":[{"index":0,"message":{"role":"assistant","content":"OK"},"finish_reason":"stop"}],
			"usage":{"prompt_tokens":5,"completion_tokens":2,"total_tokens":7}
		}`))
	}))
	defer upstream.Close()

	db := newTestChannelDB(t)
	channel := model.ProviderChannel{
		Name:         "DeepSeek Test Channel",
		ProviderType: "deepseek",
		BaseURL:      upstream.URL,
		APIKey:       "ds-test-key",
		GroupName:    "codex-group",
		ModelNames:   []string{"deepseek-chat"},
		Enabled:      true,
		TestModel:    "deepseek-chat",
		ExtraHeaders: `{"X-Test-Trace":"channel-test"}`,
	}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatalf("create channel: %v", err)
	}

	handler := ChannelHandler{
		DB:             db,
		DeepSeekClient: service.NewDeepSeekClient(10),
	}
	router := gin.New()
	router.POST("/api/channels/:id/test", handler.Test)

	request := httptest.NewRequest(http.MethodPost, "/api/channels/1/test", strings.NewReader(`{"messages":["hello from test"]}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want 200, body = %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			ChannelID uint   `json:"channel_id"`
			Model     string `json:"model"`
			Response  struct {
				ID string `json:"id"`
			} `json:"response"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != 0 {
		t.Fatalf("response code = %d, want 0", response.Code)
	}
	if response.Data.ChannelID != channel.ID {
		t.Fatalf("channel_id = %d, want %d", response.Data.ChannelID, channel.ID)
	}
	if response.Data.Model != "deepseek-chat" {
		t.Fatalf("model = %q, want deepseek-chat", response.Data.Model)
	}
	if response.Data.Response.ID != "chatcmpl-test" {
		t.Fatalf("response.id = %q, want chatcmpl-test", response.Data.Response.ID)
	}
}

func newTestChannelDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.ProviderChannel{}); err != nil {
		t.Fatalf("migrate channel: %v", err)
	}
	return db
}
