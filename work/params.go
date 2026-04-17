package work

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/APX103/openalex-go"
)

// SearchParams configures a work search request.
type SearchParams struct {
	Query   string
	Page    int
	PerPage int
	Sort    *openalex.SortOption
	Select  []string
	Filters map[string]string
	GroupBy string // OpenAlex group_by field, e.g. "publication_year", "type", "primary_topic.field.id"
}

// toQuery converts SearchParams into url.Values for the API request.
func (p SearchParams) toQuery() url.Values {
	q := make(url.Values)
	if p.Query != "" {
		q.Set("search", p.Query)
	}
	pp := openalex.PageParams{Page: p.Page, PerPage: p.PerPage}.Apply()
	q.Set("page", fmt.Sprintf("%d", pp.Page))
	q.Set("per_page", fmt.Sprintf("%d", pp.PerPage))
	if p.Sort != nil && p.Sort.Field != "" {
		q.Set("sort", p.Sort.String())
	}
	if len(p.Select) > 0 {
		q.Set("select", strings.Join(p.Select, ","))
	}
	if len(p.Filters) > 0 {
		parts := make([]string, 0, len(p.Filters))
		for k, v := range p.Filters {
			parts = append(parts, k+":"+v)
		}
		q.Set("filter", strings.Join(parts, ","))
	}
	if p.GroupBy != "" {
		q.Set("group_by", p.GroupBy)
	}
	return q
}
