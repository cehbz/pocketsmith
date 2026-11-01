package pocketsmith

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestCategoriesFlattensTheTree(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /users/42/categories": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`[
				{"id":1,"title":"Food","children":[{"id":2,"title":"Eating Out","children":[]},{"id":3,"title":"Eating In","children":null}]},
				{"id":4,"title":"Rent","children":[]}
			]`))
		},
	})
	got, err := c.Categories(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	want := []Category{{ID: 1, Title: "Food"}, {ID: 2, Title: "Eating Out"}, {ID: 3, Title: "Eating In"}, {ID: 4, Title: "Rent"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Categories = %+v, want %+v", got, want)
	}
}
