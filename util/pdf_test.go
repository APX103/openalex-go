package util

import "testing"

// mockPDFWork implements PDFWork for testing.
type mockPDFWork struct {
	arXivID     string
	bestOAURL   string
	oaURL       string
	doi         string
}

func (m *mockPDFWork) GetArXivID() string       { return m.arXivID }
func (m *mockPDFWork) GetBestOAPdfURL() string   { return m.bestOAURL }
func (m *mockPDFWork) GetOaURL() string          { return m.oaURL }
func (m *mockPDFWork) GetDoi() string            { return m.doi }

func TestPDFSourceName(t *testing.T) {
	tests := []struct {
		name string
		s    PDFSource
		want string
	}{
		{"arxiv", PDFSourceArXiv, "arxiv"},
		{"openalex", PDFSourceOpenAlex, "openalex"},
		{"unpaywall", PDFSourceUnpaywall, "unpaywall"},
		{"doi", PDFSourceDOI, "doi"},
		{"none", PDFSourceNone, ""},
		{"unknown", PDFSource(99), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PDFSourceName(tt.s)
			if got != tt.want {
				t.Errorf("PDFSourceName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolvePDF(t *testing.T) {
	tests := []struct {
		name      string
		work      PDFWork
		wantSrc   PDFSource
		wantURL   string
	}{
		{
			name: "arxiv takes priority",
			work: &mockPDFWork{arXivID: "2301.00001", bestOAURL: "https://example.com/pdf"},
			wantSrc: PDFSourceArXiv,
			wantURL: "https://arxiv.org/pdf/2301.00001",
		},
		{
			name: "best OA PDF",
			work: &mockPDFWork{bestOAURL: "https://example.com/best.pdf"},
			wantSrc: PDFSourceOpenAlex,
			wantURL: "https://example.com/best.pdf",
		},
		{
			name: "OA URL fallback",
			work: &mockPDFWork{oaURL: "https://example.com/oa"},
			wantSrc: PDFSourceOpenAlex,
			wantURL: "https://example.com/oa",
		},
		{
			name: "DOI redirect",
			work: &mockPDFWork{doi: "https://doi.org/10.1234/test"},
			wantSrc: PDFSourceDOI,
			wantURL: "https://doi.org/10.1234/test",
		},
		{
			name: "no access",
			work: &mockPDFWork{},
			wantSrc: PDFSourceNone,
			wantURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolvePDF(tt.work)
			if got.Source != tt.wantSrc {
				t.Errorf("ResolvePDF().Source = %v, want %v", got.Source, tt.wantSrc)
			}
			if got.URL != tt.wantURL {
				t.Errorf("ResolvePDF().URL = %q, want %q", got.URL, tt.wantURL)
			}
		})
	}
}
