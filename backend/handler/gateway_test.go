package handler

import (
	"testing"

	"ai-chat-platform/backend/model"
)

func TestResolveUpstreamModel(t *testing.T) {
	upstream, err := resolveUpstreamModel(model.ProviderChannel{
		ModelMapping: `{"deepseek-chat":"deepseek-reasoner"}`,
	}, "deepseek-chat")
	if err != nil {
		t.Fatalf("resolveUpstreamModel() error = %v", err)
	}
	if upstream != "deepseek-reasoner" {
		t.Fatalf("resolveUpstreamModel() = %q, want deepseek-reasoner", upstream)
	}
}

func TestResolveUpstreamModelFallback(t *testing.T) {
	upstream, err := resolveUpstreamModel(model.ProviderChannel{
		ModelMapping: `{"deepseek-chat":"deepseek-reasoner"}`,
	}, "other-model")
	if err != nil {
		t.Fatalf("resolveUpstreamModel() error = %v", err)
	}
	if upstream != "other-model" {
		t.Fatalf("resolveUpstreamModel() = %q, want other-model", upstream)
	}
}

func TestResolveUpstreamModelInvalidJSON(t *testing.T) {
	_, err := resolveUpstreamModel(model.ProviderChannel{
		ModelMapping: `{bad-json}`,
	}, "deepseek-chat")
	if err == nil {
		t.Fatal("resolveUpstreamModel() error = nil, want error")
	}
}
