package source

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/APX103/openalex-go"
)

// Search searches for sources (journals) matching the query.
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Source], error) {
	q := params.toQuery()
	var result openalex.ListResponse[Source]
	if err := c.DoRequest(ctx, "/sources", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a single source by ID.
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Source, error) {
	q := make(url.Values)
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result Source
	if err := c.DoRequest(ctx, "/sources/"+id, q, &result); err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("source %s not found", id)
	}
	return &result, nil
}
