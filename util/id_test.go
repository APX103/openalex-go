package util

import "testing"

func TestShortID(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"work ID", "https://openalex.org/W2626778328", "W2626778328"},
		{"author ID", "https://openalex.org/A5023888391", "A5023888391"},
		{"source ID", "https://openalex.org/S137773608", "S137773608"},
		{"just ID", "W12345", "W12345"},
		{"empty returns dot", "", "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShortID(tt.url)
			if got != tt.want {
				t.Errorf("ShortID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestJoinPipe(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
		want string
	}{
		{"single", []string{"W1"}, "W1"},
		{"multiple", []string{"W1", "W2", "W3"}, "W1|W2|W3"},
		{"empty slice", []string{}, ""},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JoinPipe(tt.ids)
			if got != tt.want {
				t.Errorf("JoinPipe() = %q, want %q", got, tt.want)
			}
		})
	}
}
