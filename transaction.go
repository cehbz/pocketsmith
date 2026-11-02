package pocketsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Category is a PocketSmith category.
type Category struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// Account is the transaction account a transaction belongs to.
type Account struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Transaction is one PocketSmith transaction.
type Transaction struct {
	ID             int64
	Date           time.Time // calendar date, midnight UTC
	Amount         float64   // a debit is negative
	Payee          string
	Category       *Category // nil when uncategorized
	NeedsReview    bool
	Labels         []string
	Account        Account
	Transfer       bool
	ClosingBalance float64 // the account's balance after this transaction; zero when PocketSmith gives none
}

type transaction struct {
	ID             int64     `json:"id"`
	Date           string    `json:"date"`
	Amount         float64   `json:"amount"`
	Payee          string    `json:"payee"`
	Category       *Category `json:"category"`
	NeedsReview    bool      `json:"needs_review"`
	Labels         []string  `json:"labels"`
	Account        Account   `json:"transaction_account"`
	Transfer       *bool     `json:"is_transfer"`
	ClosingBalance float64   `json:"closing_balance"`
}

func (t transaction) decode() (Transaction, error) {
	d, err := time.Parse(time.DateOnly, t.Date)
	if err != nil {
		return Transaction{}, fmt.Errorf("transaction %d date %q: %w", t.ID, t.Date, err)
	}
	return Transaction{
		ID:             t.ID,
		Date:           d,
		Amount:         t.Amount,
		Payee:          t.Payee,
		Category:       t.Category,
		NeedsReview:    t.NeedsReview,
		Labels:         t.Labels,
		Account:        t.Account,
		Transfer:       t.Transfer != nil && *t.Transfer,
		ClosingBalance: t.ClosingBalance,
	}, nil
}

// Query selects the transactions of one user or one transaction account dated From to To
// inclusive, optionally only those needing review.
type Query struct {
	User, Account int64 // exactly one is set
	From, To      time.Time
	NeedsReview   bool
}

// Transactions lists the transactions q selects, following every page.
func (c *Client) Transactions(ctx context.Context, q Query) ([]Transaction, error) {
	var path string
	switch {
	case q.User != 0 && q.Account == 0:
		path = fmt.Sprintf("/users/%d/transactions", q.User)
	case q.Account != 0 && q.User == 0:
		path = fmt.Sprintf("/transaction_accounts/%d/transactions", q.Account)
	default:
		return nil, errors.New("pocketsmith: a query names exactly one of a user and an account")
	}
	v := url.Values{}
	v.Set("start_date", q.From.Format(time.DateOnly))
	v.Set("end_date", q.To.Format(time.DateOnly))
	v.Set("per_page", "100")
	if q.NeedsReview {
		v.Set("needs_review", "1")
	}
	var out []Transaction
	for next := c.BaseURL + path + "?" + v.Encode(); next != ""; {
		var page []transaction
		var err error
		if next, err = c.do(ctx, http.MethodGet, next, nil, &page); err != nil {
			return nil, err
		}
		for _, t := range page {
			tx, err := t.decode()
			if err != nil {
				return nil, err
			}
			out = append(out, tx)
		}
	}
	return out, nil
}

// Transaction reads one transaction.
func (c *Client) Transaction(ctx context.Context, id int64) (Transaction, error) {
	var t transaction
	if _, err := c.do(ctx, http.MethodGet, fmt.Sprintf("%s/transactions/%d", c.BaseURL, id), nil, &t); err != nil {
		return Transaction{}, err
	}
	return t.decode()
}
