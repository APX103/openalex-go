# openalex-go

Go SDK for the [OpenAlex](https://openalex.org/) academic knowledge graph API.

## Install

```bash
go get github.com/APX103/openalex-go
```

## Quick Start

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/APX103/openalex-go"
    "github.com/APX103/openalex-go/work"
    "github.com/APX103/openalex-go/util"
)

func main() {
    c := openalex.New(openalex.WithAPIKey("your-key"))

    // Search papers
    resp, err := work.Search(context.Background(), c, work.SearchParams{
        Query:   "large language model",
        PerPage: 10,
        Sort:    &openalex.SortOption{Field: "cited_by_count"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d papers\n", resp.Meta.Count)

    // Get single paper + resolve PDF
    if len(resp.Results) > 0 {
        w, _ := work.Get(context.Background(), c, util.ShortID(resp.Results[0].ID))
        pdf := util.ResolvePDF(w)
        fmt.Printf("PDF: %s (%s)\n", pdf.URL, util.PDFSourceName(pdf.Source))
    }
}
```

## Packages

| Package | Description |
|---------|-------------|
| `openalex-go` | Client, options, generic types (`ListResponse`, `PageParams`, `SortOption`) |
| `openalex-go/work` | Works (papers) — search, get, citations, references, related |
| `openalex-go/author` | Authors — search, get |
| `openalex-go/source` | Sources (journals) — search, get |
| `openalex-go/util` | Utilities — PDF resolution, abstract restoration, ID helpers |

## Configuration

```go
c := openalex.New(
    openalex.WithAPIKey("key"),          // Higher rate limits
    openalex.WithMailto("you@example.com"), // Polite pool fallback
    openalex.WithTimeout(30 * time.Second),
    openalex.WithBaseURL("https://your-proxy.example.com"),
)
```

## API Reference

### work

```go
work.Search(ctx, c, work.SearchParams{Query: "...", Filters: map[string]string{"publication_year": "2024"}})
work.Get(ctx, c, "W1234567890", "id", "display_name", "abstract_inverted_index")
work.GetByIDs(ctx, c, []string{"W1", "W2"})
work.GetCitedBy(ctx, c, "W123", openalex.PageParams{Page: 1, PerPage: 20})
work.GetReferencedWorks(ctx, c, "W123", openalex.DefaultPageParams())
work.GetRelated(ctx, c, "W123", openalex.DefaultPageParams())
work.GetByAuthor(ctx, c, "A123", openalex.PageParams{PerPage: 50})
work.GetBySource(ctx, c, "S123", openalex.DefaultPageParams(), &openalex.SortOption{Field: "publication_date"})
```

### author

```go
author.Search(ctx, c, author.SearchParams{Query: "Andrew Ng"})
author.Get(ctx, c, "A5023898321")
```

### source

```go
source.Search(ctx, c, source.SearchParams{Query: "Nature"})
source.Get(ctx, c, "S137773608")
```

### util

```go
util.ShortID("https://openalex.org/W123")  // "W123"
util.JoinPipe([]string{"W1", "W2"})         // "W1|W2"
util.RestoreAbstract(work.AbstractInvertedIndex)  // plain text
util.ResolvePDF(&w)  // PDFResult{URL, Source}
util.PDFSourceName(util.PDFSourceArXiv)     // "arxiv"
```

## License

MIT
