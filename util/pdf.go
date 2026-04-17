package util

import "fmt"

// PDFSource indicates where a PDF URL was resolved from.
type PDFSource int

const (
	PDFSourceArXiv PDFSource = iota + 1
	PDFSourceOpenAlex
	PDFSourceUnpaywall // reserved
	PDFSourceDOI
	PDFSourceNone
)

// PDFSourceName returns a human-readable name for the PDF source.
func PDFSourceName(s PDFSource) string {
	switch s {
	case PDFSourceArXiv:
		return "arxiv"
	case PDFSourceOpenAlex:
		return "openalex"
	case PDFSourceUnpaywall:
		return "unpaywall"
	case PDFSourceDOI:
		return "doi"
	default:
		return ""
	}
}

// PDFResult holds a resolved PDF URL and its source.
type PDFResult struct {
	URL    string
	Source PDFSource
}

// PDFWork is the interface needed by ResolvePDF.
type PDFWork interface {
	GetArXivID() string
	GetBestOAPdfURL() string
	GetOaURL() string
	GetDoi() string
}

// ResolvePDF resolves a PDF URL by priority:
// 1. arXiv direct link (100% reliable)
// 2. OpenAlex OA (best_oa_location.pdf_url or open_access.oa_url)
// 3. Unpaywall API (reserved, not implemented)
// 4. DOI redirect (publisher page)
// 5. No access
func ResolvePDF(w PDFWork) PDFResult {
	if w.GetArXivID() != "" {
		return PDFResult{
			URL:    fmt.Sprintf("https://arxiv.org/pdf/%s", w.GetArXivID()),
			Source: PDFSourceArXiv,
		}
	}
	if w.GetBestOAPdfURL() != "" {
		return PDFResult{URL: w.GetBestOAPdfURL(), Source: PDFSourceOpenAlex}
	}
	if w.GetOaURL() != "" {
		return PDFResult{URL: w.GetOaURL(), Source: PDFSourceOpenAlex}
	}
	if w.GetDoi() != "" {
		return PDFResult{
			URL:    fmt.Sprintf("https://doi.org/%s", w.GetDoi()),
			Source: PDFSourceDOI,
		}
	}
	return PDFResult{Source: PDFSourceNone}
}
