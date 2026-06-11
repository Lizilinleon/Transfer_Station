package service

import (
	"net/http"
	"testing"
)

func TestBuildChatCompletionsEndpoint(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
	}{
		{
			name:    "deepseek official base",
			baseURL: "https://api.deepseek.com",
			want:    "https://api.deepseek.com/chat/completions",
		},
		{
			name:    "openai compatible v1 base",
			baseURL: "https://api.deepseek.com/v1",
			want:    "https://api.deepseek.com/v1/chat/completions",
		},
		{
			name:    "full endpoint",
			baseURL: "https://api.deepseek.com/chat/completions",
			want:    "https://api.deepseek.com/chat/completions",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildChatCompletionsEndpoint(test.baseURL)
			if got != test.want {
				t.Fatalf("buildChatCompletionsEndpoint() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestApplyExtraHeaders(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	err = applyExtraHeaders(request, `{"X-Test":"ok","X-Trace":"trace-id"}`)
	if err != nil {
		t.Fatalf("applyExtraHeaders() error = %v", err)
	}
	if request.Header.Get("X-Test") != "ok" {
		t.Fatalf("X-Test header = %q, want ok", request.Header.Get("X-Test"))
	}
	if request.Header.Get("X-Trace") != "trace-id" {
		t.Fatalf("X-Trace header = %q, want trace-id", request.Header.Get("X-Trace"))
	}
}

func TestApplyExtraHeadersInvalidJSON(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	if err := applyExtraHeaders(request, `{bad-json}`); err == nil {
		t.Fatal("applyExtraHeaders() error = nil, want error")
	}
}
