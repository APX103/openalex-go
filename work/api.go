package work

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/APX103/openalex-go"
	"github.com/APX103/openalex-go/util"
)

// Search searches for works matching the query.
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Work], error) {
	q := params.toQuery()
	var result openalex.ListResponse[Work]
	if err := c.DoRequest(ctx, "/works", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GroupBy returns facet buckets for works matching the given query and filters.
// The GroupBy field in params specifies which field to aggregate (e.g. "type", "publication_year").
// Note: OpenAlex uses per_page to control the number of groups returned, so we set it to 50
// to ensure all groups are included (e.g. 26 fields, ~24 types).
func GroupBy(ctx context.Context, c *openalex.Client, params SearchParams) ([]openalex.GroupBy, error) {
	q := params.toQuery()
	q.Set("per_page", "50")
	var result openalex.ListResponse[Work]
	if err := c.DoRequest(ctx, "/works", q, &result); err != nil {
		return nil, err
	}
	return result.GroupBy, nil
}

// Get retrieves a single work by ID.
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Work, error) {
	q := make(url.Values)
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result Work
	if err := c.DoRequest(ctx, "/works/"+id, q, &result); err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("work %s not found", id)
	}
	return &result, nil
}

// GetByIDs batch-fetches works by OpenAlex IDs.
// OpenAlex limits batch requests to 200 IDs; an error is returned for larger batches.
func GetByIDs(ctx context.Context, c *openalex.Client, ids []string, selectFields ...string) ([]Work, error) {
	if len(ids) > 200 {
		return nil, fmt.Errorf("GetByIDs: max 200 IDs per request, got %d", len(ids))
	}
	q := make(url.Values)
	q.Set("filter", "openalex:"+util.JoinPipe(ids))
	q.Set("per_page", fmt.Sprintf("%d", len(ids)))
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result openalex.ListResponse[Work]
	if err := c.DoRequest(ctx, "/works", q, &result); err != nil {
		return nil, err
	}
	return result.Results, nil
}

// GetCitedBy returns works that cite the given work.
func GetCitedBy(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error) {
	page = page.Apply()
	q := make(url.Values)
	q.Set("filter", fmt.Sprintf("cites:%s", workID))
	q.Set("page", fmt.Sprintf("%d", page.Page))
	q.Set("per_page", fmt.Sprintf("%d", page.PerPage))
	q.Set("sort", "cited_by_count:desc")
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result openalex.ListResponse[Work]
	if err := c.DoRequest(ctx, "/works", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetReferencedWorks returns works referenced by the given work.
func GetReferencedWorks(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error) {
	page = page.Apply()
	q := make(url.Values)
	q.Set("filter", fmt.Sprintf("cited_by:%s", workID))
	q.Set("page", fmt.Sprintf("%d", page.Page))
	q.Set("per_page", fmt.Sprintf("%d", page.PerPage))
	q.Set("sort", "cited_by_count:desc")
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result openalex.ListResponse[Work]
	if err := c.DoRequest(ctx, "/works", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRelated returns works related to the given work.
func GetRelated(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error) {
	page = page.Apply()
	q := make(url.Values)
	q.Set("filter", fmt.Sprintf("related_to:%s", workID))
	q.Set("page", fmt.Sprintf("%d", page.Page))
	q.Set("per_page", fmt.Sprintf("%d", page.PerPage))
	q.Set("sort", "cited_by_count:desc")
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result openalex.ListResponse[Work]
	if err := c.DoRequest(ctx, "/works", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetByAuthor returns works by a specific author.
// If extraFilter is non-empty, it is appended to the filter query (e.g. "concepts.id:C123").
func GetByAuthor(ctx context.Context, c *openalex.Client, authorID string, page openalex.PageParams, extraFilter string, selectFields ...string) (*openalex.ListResponse[Work], error) {
	page = page.Apply()
	q := make(url.Values)
	f := fmt.Sprintf("author.id:%s", authorID)
	if extraFilter != "" {
		f += "," + extraFilter
	}
	q.Set("filter", f)
	q.Set("sort", "cited_by_count:desc")
	q.Set("page", fmt.Sprintf("%d", page.Page))
	q.Set("per_page", fmt.Sprintf("%d", page.PerPage))
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result openalex.ListResponse[Work]
	if err := c.DoRequest(ctx, "/works", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBySource returns works published in a specific source/journal.
func GetBySource(ctx context.Context, c *openalex.Client, sourceID string, page openalex.PageParams, sort *openalex.SortOption, selectFields ...string) (*openalex.ListResponse[Work], error) {
	page = page.Apply()
	q := make(url.Values)
	q.Set("filter", fmt.Sprintf("primary_location.source.id:%s", sourceID))
	if sort != nil {
		q.Set("sort", sort.String())
	} else {
		q.Set("sort", "cited_by_count:desc")
	}
	q.Set("page", fmt.Sprintf("%d", page.Page))
	q.Set("per_page", fmt.Sprintf("%d", page.PerPage))
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result openalex.ListResponse[Work]
	if err := c.DoRequest(ctx, "/works", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
