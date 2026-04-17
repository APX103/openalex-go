//go:build integration

package openalex

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()
	opts := []Option{WithMailto("test@example.com")}
	if key := os.Getenv("OPENALEX_API_KEY"); key != "" {
		opts = []Option{WithAPIKey(key)}
	}
	return New(opts...)
}

func TestNewWithOptions(t *testing.T) {
	c := New(
		WithAPIKey("test-key"),
		WithMailto("test@example.com"),
		WithBaseURL("https://example.com"),
		WithTimeout(30*time.Second),
	)
	if c.apiKey != "test-key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "test-key")
	}
	if c.mailto != "test@example.com" {
		t.Errorf("mailto = %q, want %q", c.mailto, "test@example.com")
	}
	if c.baseURL != "https://example.com" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "https://example.com")
	}
	if c.httpClient.Timeout != 30*time.Second {
		t.Errorf("timeout = %v, want %v", c.httpClient.Timeout, 30*time.Second)
	}
}

func TestNewDefaults(t *testing.T) {
	c := New()
	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, DefaultBaseURL)
	}
	if c.httpClient.Timeout != 15*time.Second {
		t.Errorf("timeout = %v, want %v", c.httpClient.Timeout, 15*time.Second)
	}
	if c.apiKey != "" || c.mailto != "" {
		t.Error("apiKey and mailto should be empty by default")
	}
}

func TestDoRequest(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result map[string]any
	err := c.DoRequest(ctx, "/works/W2626778328", url.Values{}, &result)
	if err != nil {
		t.Fatalf("DoRequest() error = %v", err)
	}
	if id, ok := result["id"].(string); !ok || id == "" {
		t.Error("expected non-empty id in response")
	}
}

func TestDoRequestError(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result map[string]any
	err := c.DoRequest(ctx, "/works/W9999999999", url.Values{}, &result)
	if err == nil {
		t.Fatal("expected error for invalid work ID")
	}
	if apiErr, ok := err.(*APIError); ok {
		if apiErr.StatusCode == 200 {
			t.Errorf("expected non-200 status, got %d", apiErr.StatusCode)
		}
	} else {
		t.Errorf("expected *APIError, got %T", err)
	}
}
