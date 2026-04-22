# openalex-go API Reference

Go SDK for the [OpenAlex](https://openalex.org/) academic knowledge graph API.

## Install

```bash
go get github.com/APX103/openalex-go
```

## Packages

| Package | Import Path | Description |
|---------|-------------|-------------|
| `openalex` | `github.com/APX103/openalex-go` | Client, options, generic types |
| `work` | `github.com/APX103/openalex-go/work` | Works (papers) — search, get, citations, references, related |
| `author` | `github.com/APX103/openalex-go/author` | Authors — search, get |
| `source` | `github.com/APX103/openalex-go/source` | Sources (journals) — search, get |
| `util` | `github.com/APX103/openalex-go/util` | Utilities — PDF resolution, abstract restoration, ID helpers |
| `metrics` | `github.com/APX103/openalex-go/metrics` | Prometheus metrics recorder |

---

## Client Initialization

```go
c := openalex.New(
    openalex.WithAPIKey("your-key"),              // Higher rate limits
    openalex.WithMailto("you@example.com"),       // Polite pool fallback
    openalex.WithTimeout(30 * time.Second),       // Request timeout (default 15s)
    openalex.WithBaseURL("https://proxy.example.com"), // Custom base URL
    openalex.WithHTTPClient(customClient),        // Custom *http.Client
    openalex.WithRecorder(recorder),              // Prometheus metrics recorder
)
```

### Configuration Options

| Option | Type | Description |
|--------|------|-------------|
| `WithAPIKey(key)` | `string` | OpenAlex API key for higher rate limits |
| `WithMailto(email)` | `string` | Email for the polite pool (fallback when no API key) |
| `WithTimeout(d)` | `time.Duration` | HTTP request timeout (default 15s) |
| `WithBaseURL(u)` | `string` | Custom API base URL (for testing/proxies) |
| `WithHTTPClient(hc)` | `*http.Client` | Custom HTTP client |
| `WithRecorder(r)` | `RequestRecorder` | Metrics recorder (see [Metrics](#metrics) section) |

---

## Common Types

### Pagination

```go
type PageParams struct {
    Page    int // page number, default 1
    PerPage int // items per page, default 25, max 200
}
```

### Sorting

```go
type SortOption struct {
    Field string // e.g. "relevance_score", "cited_by_count", "publication_date"
    Order string // "desc" (default) or "asc"
}
```

### Generic Response

```go
type ListResponse[T any] struct {
    Meta    Meta      `json:"meta"`
    Results []T       `json:"results"`
    GroupBy []GroupBy `json:"group_by,omitempty"`
}

type Meta struct {
    Count int `json:"count"` // total number of results
}

type GroupBy struct {
    Key            string `json:"key"`
    KeyDisplayName string `json:"key_display_name"`
    Count          int    `json:"count"`
}
```

---

## work Package

Import: `github.com/APX103/openalex-go/work`

### SearchParams

```go
type SearchParams struct {
    Query   string                 // Full-text search query
    Page    int                    // Page number (default 1)
    PerPage int                    // Results per page (default 25, max 200)
    Sort    *openalex.SortOption   // Sort field and order
    Select  []string               // Fields to return (e.g. "id", "display_name")
    Filters map[string]string      // Key-value filters (e.g. {"publication_year": "2024"})
    GroupBy string                 // Aggregation field (e.g. "type", "publication_year")
}
```

### Functions

#### `work.Search`

Search for works matching a query.

```go
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Work], error)
```

**Example:**

```go
resp, err := work.Search(ctx, c, work.SearchParams{
    Query:   "large language model",
    PerPage: 10,
    Sort:    &openalex.SortOption{Field: "cited_by_count"},
    Filters: map[string]string{"publication_year": "2024"},
})
```

---

#### `work.GroupBy`

Return facet/aggregation buckets for works.

```go
func GroupBy(ctx context.Context, c *openalex.Client, params SearchParams) ([]openalex.GroupBy, error)
```

**Example:**

```go
groups, err := work.GroupBy(ctx, c, work.SearchParams{
    Filters: map[string]string{"publication_year": "2024"},
    GroupBy: "type",
})
```

---

#### `work.Get`

Retrieve a single work by OpenAlex ID.

```go
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Work, error)
```

**Example:**

```go
w, err := work.Get(ctx, c, "W2626778328", "id", "display_name", "abstract_inverted_index")
```

---

#### `work.GetByIDs`

Batch-fetch up to 200 works by OpenAlex IDs.

```go
func GetByIDs(ctx context.Context, c *openalex.Client, ids []string, selectFields ...string) ([]Work, error)
```

**Example:**

```go
works, err := work.GetByIDs(ctx, c, []string{"W1", "W2", "W3"})
```

---

#### `work.GetCitedBy`

Return works that cite the given work.

```go
func GetCitedBy(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error)
```

**Example:**

```go
resp, err := work.GetCitedBy(ctx, c, "W123", openalex.PageParams{Page: 1, PerPage: 20})
```

---

#### `work.GetReferencedWorks`

Return works referenced (in the bibliography) by the given work.

```go
func GetReferencedWorks(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error)
```

**Example:**

```go
resp, err := work.GetReferencedWorks(ctx, c, "W123", openalex.PageParams{PerPage: 50})
```

---

#### `work.GetRelated`

Return works related to the given work (OpenAlex similarity).

```go
func GetRelated(ctx context.Context, c *openalex.Client, workID string, page openalex.PageParams, selectFields ...string) (*openalex.ListResponse[Work], error)
```

**Example:**

```go
resp, err := work.GetRelated(ctx, c, "W123", openalex.PageParams{PerPage: 10})
```

---

#### `work.GetByAuthor`

Return works by a specific author. Optionally append `extraFilter` to narrow results.

```go
func GetByAuthor(ctx context.Context, c *openalex.Client, authorID string, page openalex.PageParams, extraFilter string, selectFields ...string) (*openalex.ListResponse[Work], error)
```

**Example:**

```go
// All works by author
resp, err := work.GetByAuthor(ctx, c, "A5023888391", openalex.PageParams{PerPage: 50}, "")

// Works by author filtered by concept
resp, err := work.GetByAuthor(ctx, c, "A5023888391", openalex.PageParams{PerPage: 50}, "concepts.id:C123")
```

---

#### `work.GetBySource`

Return works published in a specific source/journal.

```go
func GetBySource(ctx context.Context, c *openalex.Client, sourceID string, page openalex.PageParams, sort *openalex.SortOption, selectFields ...string) (*openalex.ListResponse[Work], error)
```

**Example:**

```go
resp, err := work.GetBySource(ctx, c, "S137773604", openalex.PageParams{PerPage: 25},
    &openalex.SortOption{Field: "publication_date"})
```

---

### Work Struct

```go
type Work struct {
    ID           string                 // OpenAlex ID, e.g. "https://openalex.org/W2626778328"
    Doi          string                 // DOI, e.g. "https://doi.org/10.1038/..."
    DisplayName  string                 // Paper title
    PubYear      int                    // Publication year
    PubDate      string                 // Publication date (YYYY-MM-DD)
    Type         string                 // e.g. "article", "preprint", "conference-paper"
    Language     string                 // ISO language code
    OpenAccess   OpenAccess             // Open access status and URL
    Authorships  []Authorship           // Authors and their institutions
    PrimaryLoc   *PrimaryLocation       // Where the work is published
    BestOALoc    *PrimaryLocation       // Best open access location
    Topics       []WorkTopic            // Topics (field/domain hierarchy)
    Concepts     []Concept              // Legacy concepts with scores
    Keywords     []Keyword              // Keywords
    Refs         []string               // Referenced work IDs
    Related      []string               // Related work IDs
    CountsByYear []CountByYear          // Yearly citation/works counts
    CitedByCount int                    // Total citations
    AbstractInv  map[string][]int       // Abstract as inverted index
    Biblio       *Biblio                // Volume, issue, pages
    IDs          WorkIDs                // External IDs (DOI, ArXiv, PMID, etc.)
}
```

---

## author Package

Import: `github.com/APX103/openalex-go/author`

### SearchParams

```go
type SearchParams struct {
    Query   string   // Search query
    Page    int      // Page number
    PerPage int      // Results per page
    Select  []string // Fields to return
}
```

### Functions

#### `author.Search`

Search for authors matching a query.

```go
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Author], error)
```

**Example:**

```go
resp, err := author.Search(ctx, c, author.SearchParams{Query: "Andrew Ng", PerPage: 5})
```

---

#### `author.Get`

Retrieve a single author by OpenAlex ID.

```go
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Author, error)
```

**Example:**

```go
a, err := author.Get(ctx, c, "A5023888391")
```

---

### Author Struct

```go
type Author struct {
    ID             string        // OpenAlex ID
    DisplayName    string        // Author name
    Orcid          string        // ORCID
    WorksCount     int           // Total works
    CitedByCount   int           // Total citations
    SummaryStats   SummaryStats  // h_index, i10_index, 2yr_mean_citedness
    LastKnownInsts []Institution // Affiliated institutions
    Topics         []AuthorTopic // Research topics
    XConcepts      []Concept     // Legacy concepts
    CountsByYear   []CountByYear // Yearly works/citation counts
    WorksAPIURL    string        // Direct link to author's works
}
```

---

## source Package

Import: `github.com/APX103/openalex-go/source`

### SearchParams

```go
type SearchParams struct {
    Query   string   // Search query
    Page    int      // Page number
    PerPage int      // Results per page
    Select  []string // Fields to return
}
```

### Functions

#### `source.Search`

Search for sources (journals/repositories) matching a query.

```go
func Search(ctx context.Context, c *openalex.Client, params SearchParams) (*openalex.ListResponse[Source], error)
```

**Example:**

```go
resp, err := source.Search(ctx, c, source.SearchParams{Query: "Nature", PerPage: 10})
```

---

#### `source.Get`

Retrieve a single source by OpenAlex ID.

```go
func Get(ctx context.Context, c *openalex.Client, id string, selectFields ...string) (*Source, error)
```

**Example:**

```go
s, err := source.Get(ctx, c, "S137773604")
```

---

### Source Struct

```go
type Source struct {
    ID           string        // OpenAlex ID
    DisplayName  string        // Journal/repository name
    ISSN         []string      // ISSNs
    ISSNL        string        // Linking ISSN
    IsOA         bool          // Is open access
    Type         string        // e.g. "journal", "repository", "conference"
    WorksCount   int           // Total works
    CitedByCount int           // Total citations
    SummaryStats SummaryStats  // h_index, i10_index, 2yr_mean_citedness
    HomepageURL  *string       // Homepage URL
    HostOrgName  *string       // Host organization name
    APCUSD       *float64      // Article processing charge in USD
    CountryCode  string        // ISO country code
    Topics       []AuthorTopic // Topics
    CountsByYear []CountByYear // Yearly works/citation counts
    WorksAPIURL  string        // Direct link to source's works
}
```

---

## util Package

Import: `github.com/APX103/openalex-go/util`

### Functions

#### `util.ShortID`

Extract the short ID from an OpenAlex URL.

```go
func ShortID(openalexURL string) string
```

```go
util.ShortID("https://openalex.org/W2626778328") // → "W2626778328"
```

---

#### `util.JoinPipe`

Join IDs with pipe separators for OpenAlex filter queries.

```go
func JoinPipe(ids []string) string
```

```go
util.JoinPipe([]string{"W1", "W2", "W3"}) // → "W1|W2|W3"
```

---

#### `util.RestoreAbstract`

Convert an OpenAlex inverted index back to plain text.

```go
func RestoreAbstract(idx map[string][]int) string
```

```go
text := util.RestoreAbstract(work.AbstractInv)
```

---

#### `util.ResolvePDF`

Resolve a PDF URL by priority: arXiv → OpenAlex OA → DOI → none.

```go
func ResolvePDF(w PDFWork) PDFResult
```

```go
pdf := util.ResolvePDF(w)
fmt.Println(pdf.URL)                    // PDF download URL
fmt.Println(util.PDFSourceName(pdf.Source)) // "arxiv", "openalex", "doi", or ""
```

---

#### `util.PDFSourceName`

Return a human-readable name for a PDF source.

```go
func PDFSourceName(s PDFSource) string
```

```go
util.PDFSourceName(util.PDFSourceArXiv) // → "arxiv"
```

---

## metrics Package

Import: `github.com/APX103/openalex-go/metrics`

Prometheus metrics for monitoring OpenAlex API calls.

### Usage

```go
import (
    "github.com/APX103/openalex-go"
    "github.com/APX103/openalex-go/metrics"
)

recorder := metrics.NewPrometheusRecorder()
c := openalex.New(openalex.WithRecorder(recorder))
```

### Recorded Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `openalex_request_duration_seconds` | Histogram | `endpoint` | Request duration buckets (0.1, 0.25, 0.5, 1, 2.5, 5, 10s) |
| `openalex_request_duration_seconds_summary` | Summary | `endpoint` | Request duration quantiles (p50, p90, p99) |
| `openalex_requests_total` | Counter | `endpoint`, `status` | Total requests by endpoint and HTTP status (0 = network error) |

### Configuration Options

```go
recorder := metrics.NewPrometheusRecorder(
    metrics.WithRegisterer(customRegistry),  // Custom Prometheus registry
    metrics.WithNamespace("myapp"),           // Metric prefix (default "openalex")
    metrics.WithSubsystem("api"),             // Subsystem prefix
    metrics.WithBuckets([]float64{0.05, 0.1, 0.5, 1, 5}), // Custom histogram buckets
)
```

---

## Not Implemented

The following OpenAlex entities are **not** yet wrapped by this SDK:

- Institutions (`/institutions`)
- Concepts (`/concepts`)
- Topics (`/topics`)
- Keywords (`/keywords`)
- Funders (`/funders`)
- Publishers (`/publishers`)

You can still access these via `Client.DoRequest` directly.
