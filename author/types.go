package author

import (
	"github.com/APX103/openalex-go"
	"github.com/APX103/openalex-go/work"
)

// Author represents a scholar from OpenAlex.
type Author struct {
	ID             string                 `json:"id"`
	DisplayName    string                 `json:"display_name"`
	Orcid          string                 `json:"orcid,omitempty"`
	WorksCount     int                    `json:"works_count"`
	CitedByCount   int                    `json:"cited_by_count"`
	SummaryStats   openalex.SummaryStats  `json:"summary_stats"`
	LastKnownInsts []work.Institution     `json:"last_known_institutions"`
	Topics         []openalex.AuthorTopic `json:"topics,omitempty"`
	XConcepts      []openalex.Concept     `json:"x_concepts,omitempty"`
	CountsByYear   []openalex.CountByYear `json:"counts_by_year,omitempty"`
	WorksAPIURL    string                 `json:"works_api_url,omitempty"`
}
