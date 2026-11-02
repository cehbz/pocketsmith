package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func date(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse(time.DateOnly, s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestTransactionsPagesThroughAUser(t *testing.T) {
	var srvURL string
	c, srv := server(t, map[string]http.HandlerFunc{
		"GET /users/42/transactions": func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("page") == "2" {
				w.Write([]byte(`[{"id":3,"date":"2026-01-13","amount":5.00,"payee":"REFUND CO","category":{"id":9,"title":"Misc"},"needs_review":false,"labels":[],"transaction_account":{"id":7,"name":"Card"},"is_transfer":false}]`))
				return
			}
			if q.Get("start_date") != "2026-01-01" || q.Get("end_date") != "2026-01-31" || q.Get("per_page") != "100" || q.Has("needs_review") {
				t.Errorf("query %s", r.URL.RawQuery)
			}
			w.Header().Set("Link", fmt.Sprintf(`<%s/users/42/transactions?page=2>; rel="next"`, srvURL))
			w.Write([]byte(`[
				{"id":1,"date":"2026-01-12","amount":-2.50,"payee":"CAFE ONE","category":{"id":4,"title":"Eating Out"},"needs_review":true,"labels":["x"],"transaction_account":{"id":7,"name":"Card"},"is_transfer":false},
				{"id":2,"date":"2026-01-12","amount":-100.00,"payee":"TO SAVINGS","category":{"id":5,"title":"Savings"},"needs_review":false,"labels":[],"transaction_account":{"id":8,"name":"Checking"},"is_transfer":true}
			]`))
		},
	})
	srvURL = srv.URL
	got, err := c.Transactions(context.Background(), Query{User: 42, From: date(t, "2026-01-01"), To: date(t, "2026-01-31")})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].ID != 1 || got[1].ID != 2 || got[2].ID != 3 {
		t.Fatalf("transactions = %+v, want IDs 1, 2, 3", got)
	}
	one := got[0]
	if !one.Date.Equal(date(t, "2026-01-12")) || one.Amount != -2.50 || one.Payee != "CAFE ONE" || !one.NeedsReview {
		t.Errorf("transaction 1 = %+v", one)
	}
	if one.Category == nil || one.Category.ID != 4 || one.Category.Title != "Eating Out" {
		t.Errorf("transaction 1 category = %+v", one.Category)
	}
	if one.Account != (Account{ID: 7, Name: "Card"}) || len(one.Labels) != 1 || one.Labels[0] != "x" || one.Transfer {
		t.Errorf("transaction 1 = %+v", one)
	}
	if !got[1].Transfer {
		t.Errorf("transaction 2 is not a transfer")
	}
}

func TestTransactionsOfAnAccountNeedingReview(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /transaction_accounts/7/transactions": func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("needs_review") != "1" {
				t.Errorf("query %s, want needs_review=1", r.URL.RawQuery)
			}
			w.Write([]byte(`[]`))
		},
	})
	got, err := c.Transactions(context.Background(), Query{Account: 7, From: date(t, "2026-01-01"), To: date(t, "2026-01-31"), NeedsReview: true})
	if err != nil || len(got) != 0 {
		t.Fatalf("Transactions = %v, %v", got, err)
	}
}

func TestTransactionsNeedsExactlyOneScope(t *testing.T) {
	c := &Client{BaseURL: "http://unused", Key: key, HTTP: http.DefaultClient}
	for _, q := range []Query{{}, {User: 1, Account: 2}} {
		if _, err := c.Transactions(context.Background(), q); err == nil {
			t.Errorf("Transactions(%+v) succeeded", q)
		}
	}
}

func TestTransactionsDecodesNulls(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /users/42/transactions": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`[{"id":1,"date":"2026-01-12","amount":-1,"payee":"P","category":null,"needs_review":true,"labels":null,"transaction_account":{"id":7,"name":"Card"},"is_transfer":null}]`))
		},
	})
	got, err := c.Transactions(context.Background(), Query{User: 42, From: date(t, "2026-01-01"), To: date(t, "2026-01-31")})
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Category != nil || got[0].Transfer || len(got[0].Labels) != 0 {
		t.Errorf("transaction = %+v, want uncategorized, not a transfer, no labels", got[0])
	}
}

func TestTransactionsRejectsABadDate(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /users/42/transactions": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`[{"id":1,"date":"12/01/2026","amount":-1,"payee":"P"}]`))
		},
	})
	if _, err := c.Transactions(context.Background(), Query{User: 42, From: date(t, "2026-01-01"), To: date(t, "2026-01-31")}); err == nil {
		t.Fatal("Transactions accepted a malformed date")
	}
}

func TestTransactionReadsOne(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /transactions/5": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"id":5,"date":"2026-01-12","amount":-2.75,"payee":"SHOP","category":{"id":7,"title":"Misc"},"needs_review":false,"labels":["auto"],"transaction_account":{"id":7,"name":"Card"},"is_transfer":false}`))
		},
	})
	got, err := c.Transaction(context.Background(), 5)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 5 || got.Amount != -2.75 || got.Category == nil || got.Category.ID != 7 || len(got.Labels) != 1 {
		t.Errorf("Transaction = %+v", got)
	}
}

func TestTransactionsDecodesTheClosingBalance(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /users/42/transactions": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`[
				{"id":1,"date":"2026-01-12","amount":-1,"payee":"P","closing_balance":1234.56,"transaction_account":{"id":7,"name":"Card"}},
				{"id":2,"date":"2026-01-12","amount":-1,"payee":"P","closing_balance":null,"transaction_account":{"id":7,"name":"Card"}}
			]`))
		},
	})
	got, err := c.Transactions(context.Background(), Query{User: 42, From: date(t, "2026-01-01"), To: date(t, "2026-01-31")})
	if err != nil {
		t.Fatal(err)
	}
	if got[0].ClosingBalance != 1234.56 || got[1].ClosingBalance != 0 {
		t.Errorf("closing balances = %v, %v; want 1234.56, 0", got[0].ClosingBalance, got[1].ClosingBalance)
	}
}
