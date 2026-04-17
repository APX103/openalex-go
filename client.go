package openalex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const DefaultBaseURL = "https://api.openalex.org"

// Client is an OpenAlex API client.
type Client struct {
	baseURL    string
	apiKey     string
	mailto     string
	httpClient *http.Client
}

// New creates a new Client with the given options.
func New(opts ...Option) *Client {
	c := &Client{
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// DoRequest sends an authenticated GET request and decodes the JSON response into result.
// This is exported so sub-packages (work, author, source) can use it.
func (c *Client) DoRequest(ctx context.Context, path string, params url.Values, result any) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}
	q := u.Query()
	for k, vs := range params {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	if c.apiKey != "" {
		q.Set("api_key", c.apiKey)
	} else if c.mailto != "" {
		q.Set("mailto", c.mailto)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("create request %s: %w", path, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body %s: %w", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
			URL:        u.String(),
		}
	}
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

// APIError is returned when the OpenAlex API responds with a non-200 status code.
type APIError struct {
	StatusCode int
	Message    string
	URL        string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("openalex: %s returned %d: %s", e.URL, e.StatusCode, e.Message)
}
