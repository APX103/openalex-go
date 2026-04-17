package util

// ShortID extracts the short ID from an OpenAlex URL.
// Example: "https://openalex.org/W2626778328" → "W2626778328"
func ShortID(openalexURL string) string {
	for i := len(openalexURL) - 1; i >= 0; i-- {
		if openalexURL[i] == '/' {
			return openalexURL[i+1:]
		}
	}
	return openalexURL
}

// JoinPipe joins IDs with pipe separators for OpenAlex filter queries.
// Example: ["W1", "W2"] → "W1|W2"
func JoinPipe(ids []string) string {
	result := ""
	for i, id := range ids {
		if i > 0 {
			result += "|"
		}
		result += id
	}
	return result
}
