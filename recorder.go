package openalex

import (
	"context"
	"strings"
	"time"
)

// RequestRecorder records metrics for each API request.
type RequestRecorder interface {
	RecordRequest(ctx context.Context, endpoint string, duration time.Duration, statusCode int)
}

type noopRecorder struct{}

func (noopRecorder) RecordRequest(_ context.Context, _ string, _ time.Duration, _ int) {}

func extractEndpoint(path string) string {
	path = strings.TrimPrefix(path, "/")
	if i := strings.IndexByte(path, '/'); i >= 0 {
		return path[:i]
	}
	return path
}
