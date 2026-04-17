package work

import (
	"github.com/APX103/openalex-go"
	"github.com/APX103/openalex-go/util"
)

// Work represents a scholarly work (paper) from OpenAlex.
type Work struct {
	ID           string                 `json:"id"`
	Doi          string                 `json:"doi"`
	DisplayName  string                 `json:"display_name"`
	PubYear      int                    `json:"publication_year"`
	PubDate      string                 `json:"publication_date"`
	Type         string                 `json:"type"`
	Language     string                 `json:"language,omitempty"`
	IndexedIn    []string               `json:"indexed_in,omitempty"`
	OpenAccess   OpenAccess             `json:"open_access"`
	Authorships  []Authorship           `json:"authorships"`
	PrimaryLoc   *PrimaryLocation       `json:"primary_location"`
	BestOALoc    *PrimaryLocation       `json:"best_oa_location,omitempty"`
	Topics       []WorkTopic            `json:"topics,omitempty"`
	Concepts     []openalex.Concept     `json:"concepts,omitempty"`
	Keywords     []Keyword              `json:"keywords,omitempty"`
	Refs         []string               `json:"referenced_works,omitempty"`
	Related      []string               `json:"related_works,omitempty"`
	CountsByYear []openalex.CountByYear `json:"counts_by_year,omitempty"`
	CitedByCount int                    `json:"cited_by_count"`
	AbstractInv  map[string][]int       `json:"abstract_inverted_index,omitempty"`
	Biblio       *Biblio                `json:"biblio,omitempty"`
	IDs          WorkIDs                `json:"ids,omitempty"`
}

// PDFWork interface implementation — enables util.ResolvePDF(w).
func (w *Work) GetArXivID() string { return w.IDs.ArXiv }

func (w *Work) GetBestOAPdfURL() string {
	if w.BestOALoc != nil && w.BestOALoc.PdfURL != nil {
		return *w.BestOALoc.PdfURL
	}
	return ""
}

func (w *Work) GetOaURL() string {
	if w.OpenAccess.OaURL != nil {
		return *w.OpenAccess.OaURL
	}
	return ""
}

func (w *Work) GetDoi() string { return w.Doi }

// Compile-time check: Work implements util.PDFWork.
var _ util.PDFWork = (*Work)(nil)

// WorkIDs holds external identifiers for a work.
type WorkIDs struct {
	OpenAlex string `json:"openalex"`
	Doi      string `json:"doi"`
	Mag      string `json:"mag"`
	PMID     string `json:"pmid,omitempty"`
	ArXiv    string `json:"arxiv,omitempty"`
}

// OpenAccess describes the open-access status of a work.
type OpenAccess struct {
	IsOA               bool    `json:"is_oa"`
	OaStatus           string  `json:"oa_status"`
	OaURL              *string `json:"oa_url"`
	AnyRepoHasFulltext bool    `json:"any_repository_has_fulltext"`
}

// Authorship describes an author's contribution to a work.
type Authorship struct {
	AuthorPosition  string        `json:"author_position"`
	Author          AuthorRef     `json:"author"`
	Institutions    []Institution `json:"institutions"`
	Countries       []string      `json:"countries,omitempty"`
	IsCorresponding bool          `json:"is_corresponding"`
}

// AuthorRef is a lightweight reference to an author.
type AuthorRef struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Orcid       string `json:"orcid,omitempty"`
}

// Institution is a lightweight reference to an institution.
type Institution struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Ror         string `json:"ror,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Type        string `json:"type,omitempty"`
}

// PrimaryLocation describes where a work is published (journal, repository, etc.).
type PrimaryLocation struct {
	IsOA           bool            `json:"is_oa"`
	LandingPageURL string          `json:"landing_page_url"`
	PdfURL         *string         `json:"pdf_url,omitempty"`
	Source         *LocationSource `json:"source,omitempty"`
	License        string          `json:"license,omitempty"`
	Version        string          `json:"version,omitempty"`
}

// LocationSource describes the source (journal/repository) of a location.
type LocationSource struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	ISSN        any    `json:"issn,omitempty"`
	Type        string `json:"type"`
	IsOA        bool   `json:"is_oa"`
}

// WorkTopic describes a topic associated with a work.
type WorkTopic struct {
	ID          string            `json:"id"`
	DisplayName string            `json:"display_name"`
	Count       int               `json:"count"`
	Subfield    openalex.TopicRef `json:"subfield"`
	Field       openalex.TopicRef `json:"field"`
	Domain      openalex.TopicRef `json:"domain"`
}

// Keyword is a keyword attached to a work.
type Keyword struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// Biblio holds bibliographic details (volume, issue, pages).
type Biblio struct {
	Volume    string `json:"volume,omitempty"`
	Issue     string `json:"issue,omitempty"`
	FirstPage string `json:"first_page,omitempty"`
	LastPage  string `json:"last_page,omitempty"`
}
