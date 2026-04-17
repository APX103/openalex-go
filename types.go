package openalex

// Meta contains pagination metadata.
type Meta struct {
	Count int `json:"count"`
}

// GroupBy represents a group-by result.
type GroupBy struct {
	Key            string `json:"key"`
	KeyDisplayName string `json:"key_display_name"`
	Count          int    `json:"count"`
}

// ListResponse is a generic paginated response from the OpenAlex API.
type ListResponse[T any] struct {
	Meta    Meta      `json:"meta"`
	Results []T       `json:"results"`
	GroupBy []GroupBy `json:"group_by,omitempty"`
}

// PageParams controls pagination.
type PageParams struct {
	Page    int // page number, default 1
	PerPage int // items per page, default 25, max 200
}

// Apply sets defaults to PageParams.
func (p PageParams) Apply() PageParams {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 25
	}
	if p.PerPage > 200 {
		p.PerPage = 200
	}
	return p
}

// SortOption controls result ordering.
type SortOption struct {
	Field string // e.g. "relevance_score", "cited_by_count", "publication_date"
	Order string // "desc" or "asc"
}

// String returns "field:order" format for the API.
func (s SortOption) String() string {
	if s.Order == "" {
		s.Order = "desc"
	}
	return s.Field + ":" + s.Order
}

// --- Shared types used by multiple sub-packages ---

// SummaryStats contains citation metrics (shared by Author and Source).
type SummaryStats struct {
	HIndex             int     `json:"h_index"`
	I10Index           int     `json:"i10_index"`
	TwoYrMeanCitedness float64 `json:"2yr_mean_citedness"`
}

// TopicRef is a reference to a taxonomy topic (shared by Work and Author).
type TopicRef struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// AuthorTopic describes a topic associated with an author (shared by Author and Source).
type AuthorTopic struct {
	ID          string   `json:"id"`
	DisplayName string   `json:"display_name"`
	Count       int      `json:"count"`
	Value       float64  `json:"value,omitempty"`
	Subfield    TopicRef `json:"subfield"`
	Field       TopicRef `json:"field"`
	Domain      TopicRef `json:"domain"`
}

// CountByYear contains yearly citation stats (shared by Work, Author, Source).
type CountByYear struct {
	Year         int `json:"year"`
	WorksCount   int `json:"works_count"`
	CitedByCount int `json:"cited_by_count"`
}

// Concept is a tag associated with a work.
type Concept struct {
	ID          string  `json:"id"`
	DisplayName string  `json:"display_name"`
	Score       float64 `json:"score"`
}
