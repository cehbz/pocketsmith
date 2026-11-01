package pocketsmith

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func updateBody(t *testing.T, u Update) map[string]any {
	t.Helper()
	var body map[string]any
	c, _ := server(t, map[string]http.HandlerFunc{
		"PUT /transactions/5": func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("content type %q", r.Header.Get("Content-Type"))
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			w.Write([]byte(`{"id":5}`))
		},
	})
	if err := c.Update(context.Background(), 5, u); err != nil {
		t.Fatal(err)
	}
	return body
}

func TestUpdateSendsCategoryReviewAndLabels(t *testing.T) {
	body := updateBody(t, Update{Category: 7, NeedsReview: false, Labels: []string{"x", "auto"}})
	if body["category_id"] != float64(7) || body["needs_review"] != false || body["labels"] != "x,auto" {
		t.Errorf("body = %v", body)
	}
}

func TestUpdateSendsEmptyLabelsAndUncategorized(t *testing.T) {
	body := updateBody(t, Update{Category: 0, NeedsReview: true})
	labels, ok := body["labels"]
	if !ok || labels != "" {
		t.Errorf("labels = %v (present %v), want an empty string", labels, ok)
	}
	if body["category_id"] != "" || body["needs_review"] != true {
		t.Errorf("body = %v", body)
	}
}
