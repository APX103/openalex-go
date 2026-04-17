package source

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/APX103/openalex-go"
)

// SearchParams configures a source search request.
type SearchParams struct {
	Query   string
	Page    int
	PerPage int
	Select  []string
}

func (p SearchParams) toQuery() url.Values {
	q := make(url.Values)
	if p.Query != "" {
		q.Set("search", p.Query)
	}
	pp := openalex.PageParams{Page: p.Page, PerPage: p.PerPage}.Apply()
	q.Set("page", fmt.Sprintf("%d", pp.Page))
	q.Set("per_page", fmt.Sprintf("%d", pp.PerPage))
	q.Set("sort", "relevance_score:desc")
	if len(p.Select) > 0 {
		q.Set("select", strings.Join(p.Select, ","))
	}
	return q
}
