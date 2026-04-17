package openalex

import (
	"net/http"
	"time"
)

// Option configures a Client.
type Option func(*Client)

// WithAPIKey sets the API key for higher rate limits.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

// WithMailto sets the email for the polite pool (used when no API key is set).
func WithMailto(email string) Option {
	return func(c *Client) { c.mailto = email }
}

// WithHTTPClient replaces the default http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithBaseURL sets a custom API base URL (useful for testing or proxies).
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = u }
}

// WithTimeout sets the request timeout (default 15s).
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}
