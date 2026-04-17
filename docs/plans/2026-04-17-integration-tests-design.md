# Integration Tests Design

## Overview

Add integration tests for all API functions in the `openalex-go` SDK, calling the real OpenAlex API (`api.openalex.org`).

## Testing Strategy

- **Real API calls** — no mocks, no httptest servers
- **Build tag** — all integration tests use `//go:build integration`, run via `go test -tags=integration`
- **Rate limiting** — API key via `OPENALEX_API_KEY` env var (100 req/s); fallback to `mailto` polite pool (10 req/s) with 200ms sleep between tests
- **Deterministic IDs** — use well-known stable OpenAlex IDs rather than search results

## Test Files

| File | Tests | Package |
|------|-------|---------|
| `util/abstract_test.go` | RestoreAbstract | util (no tag, runs by default) |
| `util/id_test.go` | ShortID, JoinPipe | util (no tag, runs by default) |
| `util/pdf_test.go` | ResolvePDF, PDFSourceName | util (no tag, runs by default) |
| `client_test.go` | New options, DoRequest error | openalex (integration) |
| `author/api_test.go` | Search, Get | author_test (integration) |
| `source/api_test.go` | Search, Get | source_test (integration) |
| `work/api_test.go` | Search, Get, GetByIDs, GetCitedBy, GetReferencedWorks, GetRelated, GetByAuthor, GetBySource | work_test (integration) |

## Known Stable IDs

- Work: `W2626778328` (used in project docs)
- Author: `A5023898321` (used in project docs)
- Source: `S137773608` (Nature)

## Verification Approach

- Validate structural correctness (non-empty fields, list lengths)
- Do not assert specific values that may change (titles, counts)
- Check error handling for edge cases (invalid IDs, empty results)
