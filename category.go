package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
)

type categoryNode struct {
	ID       int64          `json:"id"`
	Title    string         `json:"title"`
	Children []categoryNode `json:"children"`
}

// Categories lists every category of the user, parents before their children.
func (c *Client) Categories(ctx context.Context, user int64) ([]Category, error) {
	var tree []categoryNode
	if _, err := c.do(ctx, http.MethodGet, fmt.Sprintf("%s/users/%d/categories", c.BaseURL, user), nil, &tree); err != nil {
		return nil, err
	}
	var out []Category
	var walk func([]categoryNode)
	walk = func(nodes []categoryNode) {
		for _, n := range nodes {
			out = append(out, Category{ID: n.ID, Title: n.Title})
			walk(n.Children)
		}
	}
	walk(tree)
	return out, nil
}
