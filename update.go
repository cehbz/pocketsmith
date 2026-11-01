package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Update is the part of a transaction the client sets.
type Update struct {
	Category    int64 // zero uncategorizes
	NeedsReview bool
	Labels      []string
}

// Update sets a transaction's category, review flag and labels. Labels go as one
// comma-joined string, and an empty one clears them.
func (c *Client) Update(ctx context.Context, id int64, u Update) error {
	var category any = u.Category
	if u.Category == 0 {
		category = ""
	}
	body := map[string]any{
		"category_id":  category,
		"needs_review": u.NeedsReview,
		"labels":       strings.Join(u.Labels, ","),
	}
	_, err := c.do(ctx, http.MethodPut, fmt.Sprintf("%s/transactions/%d", c.BaseURL, id), body, nil)
	return err
}
