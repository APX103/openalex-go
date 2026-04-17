package source

import "github.com/APX103/openalex-go"

// Source represents a journal or repository from OpenAlex.
type Source struct {
	ID           string                 `json:"id"`
	DisplayName  string                 `json:"display_name"`
	ISSN         []string               `json:"issn,omitempty"`
	ISSNL        string                 `json:"issn_l,omitempty"`
	IsOA         bool                   `json:"is_oa"`
	Type         string                 `json:"type"`
	WorksCount   int                    `json:"works_count"`
	CitedByCount int                    `json:"cited_by_count"`
	SummaryStats openalex.SummaryStats  `json:"summary_stats"`
	HomepageURL  *string                `json:"homepage_url,omitempty"`
	HostOrgName  *string                `json:"host_organization_name,omitempty"`
	APCUSD       *float64               `json:"apc_usd,omitempty"`
	CountryCode  string                 `json:"country_code,omitempty"`
	Topics       []openalex.AuthorTopic `json:"topics,omitempty"`
	TopicShare   []openalex.AuthorTopic `json:"topic_share,omitempty"`
	CountsByYear []openalex.CountByYear `json:"counts_by_year,omitempty"`
	WorksAPIURL  string                 `json:"works_api_url,omitempty"`
}
