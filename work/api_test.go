//go:build integration

package work_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/APX103/openalex-go"
	"github.com/APX103/openalex-go/work"
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

	result, err := work.Search(ctx, c, work.SearchParams{
		Query:   "machine learning",
		PerPage: 5,
	})
	if err != nil {
		t.Fatalf("work.Search() error = %v", err)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 result for 'machine learning'")
	}
	if result.Meta.Count <= 0 {
		t.Errorf("expected Meta.Count > 0, got %d", result.Meta.Count)
	}
	if result.Results[0].ID == "" {
		t.Error("expected non-empty ID in first result")
	}
}

func TestSearchWithFilters(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := work.Search(ctx, c, work.SearchParams{
		Query:   "deep learning",
		PerPage: 3,
		Filters: map[string]string{
			"publication_year": "2020",
			"type":             "article",
		},
	})
	if err != nil {
		t.Fatalf("work.Search() with filters error = %v", err)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 result with filters")
	}
	for _, w := range result.Results {
		if w.PubYear != 2020 {
			t.Errorf("expected PubYear=2020, got %d", w.PubYear)
		}
	}
}

func TestGet(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	w, err := work.Get(ctx, c, "W2626778328")
	if err != nil {
		t.Fatalf("work.Get() error = %v", err)
	}
	if w.ID == "" {
		t.Error("expected non-empty ID")
	}
	if w.DisplayName == "" {
		t.Error("expected non-empty DisplayName")
	}
}

func TestGetWithSelectFields(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	w, err := work.Get(ctx, c, "W2626778328", "id", "display_name")
	if err != nil {
		t.Fatalf("work.Get() with select error = %v", err)
	}
	if w.ID == "" {
		t.Error("expected non-empty ID with select fields")
	}
	if w.DisplayName == "" {
		t.Error("expected non-empty DisplayName with select fields")
	}
}

func TestGetByIDs(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ids := []string{"W2626778328"}
	works, err := work.GetByIDs(ctx, c, ids)
	if err != nil {
		t.Fatalf("work.GetByIDs() error = %v", err)
	}
	if len(works) != 1 {
		t.Errorf("expected 1 work, got %d", len(works))
	}
	for _, w := range works {
		if w.ID == "" {
			t.Error("expected non-empty ID in batch result")
		}
	}
}

func TestGetByIDsTooMany(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ids := make([]string, 201)
	for i := range ids {
		ids[i] = "W1"
	}
	_, err := work.GetByIDs(ctx, c, ids)
	if err == nil {
		t.Fatal("expected error for > 200 IDs")
	}
}

func TestGetCitedBy(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := work.GetCitedBy(ctx, c, "W2626778328", openalex.PageParams{PerPage: 5})
	if err != nil {
		t.Fatalf("work.GetCitedBy() error = %v", err)
	}
	if result.Meta.Count <= 0 {
		t.Errorf("expected Meta.Count > 0 for highly cited work, got %d", result.Meta.Count)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 citing work")
	}
}

func TestGetReferencedWorks(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := work.GetReferencedWorks(ctx, c, "W2626778328", openalex.PageParams{PerPage: 5})
	if err != nil {
		t.Fatalf("work.GetReferencedWorks() error = %v", err)
	}
	if result.Meta.Count <= 0 {
		t.Errorf("expected Meta.Count > 0, got %d", result.Meta.Count)
	}
}

func TestGetRelated(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := work.GetRelated(ctx, c, "W2626778328", openalex.PageParams{PerPage: 5})
	if err != nil {
		t.Fatalf("work.GetRelated() error = %v", err)
	}
	if result.Meta.Count <= 0 {
		t.Errorf("expected Meta.Count > 0, got %d", result.Meta.Count)
	}
}

func TestGetByAuthor(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := work.GetByAuthor(ctx, c, "A5023898321", openalex.PageParams{PerPage: 5}, "")
	if err != nil {
		t.Fatalf("work.GetByAuthor() error = %v", err)
	}
	if result.Meta.Count <= 0 {
		t.Errorf("expected Meta.Count > 0 for Einstein, got %d", result.Meta.Count)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 work by the author")
	}
}

func TestGetBySource(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := work.GetBySource(ctx, c, "S137773608", openalex.PageParams{PerPage: 5}, nil)
	if err != nil {
		t.Fatalf("work.GetBySource() error = %v", err)
	}
	if result.Meta.Count <= 0 {
		t.Errorf("expected Meta.Count > 0 for Nature, got %d", result.Meta.Count)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 work in the source")
	}
}

func TestGroupBy(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	buckets, err := work.GroupBy(ctx, c, work.SearchParams{
		Filters: map[string]string{"publication_year": "2024"},
		GroupBy: "type",
	})
	if err != nil {
		t.Fatalf("work.GroupBy() error = %v", err)
	}
	if len(buckets) == 0 {
		t.Fatal("expected at least 1 group_by bucket")
	}
	for _, b := range buckets {
		if b.Key == "" {
			t.Error("expected non-empty Key in bucket")
		}
		if b.Count <= 0 {
			t.Errorf("expected Count > 0 for bucket %q, got %d", b.Key, b.Count)
		}
	}
}

func TestGroupByPublicationYear(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	buckets, err := work.GroupBy(ctx, c, work.SearchParams{
		Query:   "transformer",
		GroupBy: "publication_year",
	})
	if err != nil {
		t.Fatalf("work.GroupBy() error = %v", err)
	}
	if len(buckets) == 0 {
		t.Fatal("expected at least 1 group_by bucket")
	}
}

func TestGetBySourceWithSort(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sortOpt := &openalex.SortOption{Field: "publication_date", Order: "desc"}
	result, err := work.GetBySource(ctx, c, "S137773608", openalex.PageParams{PerPage: 3}, sortOpt)
	if err != nil {
		t.Fatalf("work.GetBySource() with sort error = %v", err)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least 1 work with custom sort")
	}
}
