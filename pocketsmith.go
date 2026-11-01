// Package pocketsmith is a client for the parts of PocketSmith's v2 API its users need.
package pocketsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// DefaultBaseURL is PocketSmith's v2 API.
const DefaultBaseURL = "https://api.pocketsmith.com/v2"

// Client calls PocketSmith with a developer key.
type Client struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

var nextLink = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)

// do sends a request and decodes the response into out, returning the next page's URL if any.
func (c *Client) do(ctx context.Context, method, target string, body, out any) (string, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return "", err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, target, r)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Developer-Key", c.Key)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("pocketsmith %s %s: %s: %s", method, req.URL.Path, resp.Status, bytes.TrimSpace(msg))
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return "", fmt.Errorf("pocketsmith %s %s: %w", method, req.URL.Path, err)
		}
	}
	if m := nextLink.FindStringSubmatch(resp.Header.Get("Link")); m != nil {
		return m[1], nil
	}
	return "", nil
}
