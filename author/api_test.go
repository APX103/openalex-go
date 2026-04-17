//go:build integration

package author_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/APX103/openalex-go"
	"github.com/APX103/openalex-go/author"
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

	result, err := author.Search(ctx, c, author.SearchParams{
		Query:   "Einstein",
		PerPage: 5,
	})
	if err != nil {
		t.Fatalf("author.Search() error = %v", err)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 result for 'Einstein'")
	}
	if result.Meta.Count <= 0 {
		t.Errorf("expected Meta.Count > 0, got %d", result.Meta.Count)
	}
	found := false
	for _, a := range result.Results {
		if a.ID != "" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected at least one author with non-empty ID")
	}
}

func TestGet(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	a, err := author.Get(ctx, c, "A5023898321")
	if err != nil {
		t.Fatalf("author.Get() error = %v", err)
	}
	if a.ID == "" {
		t.Error("expected non-empty ID")
	}
	if a.DisplayName == "" {
		t.Error("expected non-empty DisplayName")
	}
	if a.WorksCount <= 0 {
		t.Errorf("expected WorksCount > 0, got %d", a.WorksCount)
	}
}

func TestGetWithSelectFields(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	a, err := author.Get(ctx, c, "A5023898321", "id", "display_name", "works_count")
	if err != nil {
		t.Fatalf("author.Get() with select error = %v", err)
	}
	if a.DisplayName == "" {
		t.Error("expected non-empty DisplayName with select fields")
	}
}
