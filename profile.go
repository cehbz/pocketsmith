package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Profile is the user's identity and time zone.
type Profile struct {
	ID       int64
	TimeZone *time.Location
}

// Profile reads the user's profile and resolves its time zone to an IANA location.
func (c *Client) Profile(ctx context.Context) (Profile, error) {
	var me struct {
		ID       int64  `json:"id"`
		TimeZone string `json:"time_zone"`
	}
	if _, err := c.do(ctx, http.MethodGet, c.BaseURL+"/me", nil, &me); err != nil {
		return Profile{}, err
	}
	var zones []struct {
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
	}
	if _, err := c.do(ctx, http.MethodGet, c.BaseURL+"/time_zones", nil, &zones); err != nil {
		return Profile{}, err
	}
	for _, z := range zones {
		if z.Name == me.TimeZone {
			loc, err := time.LoadLocation(z.Identifier)
			if err != nil {
				return Profile{}, err
			}
			return Profile{ID: me.ID, TimeZone: loc}, nil
		}
	}
	return Profile{}, fmt.Errorf("PocketSmith time zone %q has no identifier", me.TimeZone)
}
