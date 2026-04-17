package author

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/xxx/openalex-go"
)

// Search searches for authors matching the query.
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Author], error) {
	q := params.toQuery()
	var result openalex.ListResponse[Author]
	if err := c.DoRequest(ctx, "/authors", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a single author by ID.
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Author, error) {
	q := make(url.Values)
	if len(selectFields) > 0 {
		q.Set("select", strings.Join(selectFields, ","))
	}
	var result Author
	if err := c.DoRequest(ctx, "/authors/"+id, q, &result); err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("author %s not found", id)
	}
	return &result, nil
}
