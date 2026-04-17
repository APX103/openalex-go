//go:build integration

package source_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/APX103/openalex-go"
	"github.com/APX103/openalex-go/source"
)

func newTestClient(t *testing.T) *openalex.Client {
	t.Helper()
	opts := []openalex.Option{openalex.WithMailto("test@example.com")}
	if key := os.Getenv("OPENALEX_API_KEY"); key != "" {
		opts = []openalex.Option{openalex.WithAPIKey(key)}
	}
	return openalex.New(opts...)
}

func TestSearch(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := source.Search(ctx, c, source.SearchParams{
		Query:   "Nature",
		PerPage: 5,
	})
	if err != nil {
		t.Fatalf("source.Search() error = %v", err)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 result for 'Nature'")
	}
	if result.Meta.Count <= 0 {
		t.Errorf("expected Meta.Count > 0, got %d", result.Meta.Count)
	}
}

func TestGet(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s, err := source.Get(ctx, c, "S137773608")
	if err != nil {
		t.Fatalf("source.Get() error = %v", err)
	}
	if s.ID == "" {
		t.Error("expected non-empty ID")
	}
	if s.DisplayName == "" {
		t.Error("expected non-empty DisplayName")
	}
	if s.WorksCount <= 0 {
		t.Errorf("expected WorksCount > 0, got %d", s.WorksCount)
	}
}

func TestGetWithSelectFields(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s, err := source.Get(ctx, c, "S137773608", "id", "display_name", "issn")
	if err != nil {
		t.Fatalf("source.Get() with select error = %v", err)
	}
	if s.DisplayName == "" {
		t.Error("expected non-empty DisplayName with select fields")
	}
}
