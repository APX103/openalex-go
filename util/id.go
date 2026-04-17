package util

import (
	"path"
	"strings"
)

// ShortID extracts the short ID from an OpenAlex URL.
// Example: "https://openalex.org/W2626778328" → "W2626778328"
func ShortID(openalexURL string) string {
	return path.Base(openalexURL)
}

// JoinPipe joins IDs with pipe separators for OpenAlex filter queries.
// Example: ["W1", "W2"] → "W1|W2"
func JoinPipe(ids []string) string {
	return strings.Join(ids, "|")
}
