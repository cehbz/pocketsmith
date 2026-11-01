package main

import (
	"encoding/json"
	"io"
	"time"

	"github.com/cehbz/pocketsmith"
)

// line is one transaction as the command prints it.
type line struct {
	ID          int64    `json:"id"`
	Date        string   `json:"date"`
	Amount      float64  `json:"amount"`
	Payee       string   `json:"payee"`
	CategoryID  int64    `json:"category_id"`
	Category    string   `json:"category"`
	NeedsReview bool     `json:"needs_review"`
	Labels      []string `json:"labels"`
	AccountID   int64    `json:"account_id"`
	Account     string   `json:"account"`
	Transfer    bool     `json:"transfer"`
}

// writeLines prints each transaction as one JSON object per line.
func writeLines(w io.Writer, txs []pocketsmith.Transaction) error {
	enc := json.NewEncoder(w)
	for _, t := range txs {
		l := line{
			ID: t.ID, Date: t.Date.Format(time.DateOnly), Amount: t.Amount, Payee: t.Payee,
			NeedsReview: t.NeedsReview, Labels: t.Labels, AccountID: t.Account.ID, Account: t.Account.Name,
			Transfer: t.Transfer,
		}
		if l.Labels == nil {
			l.Labels = []string{}
		}
		if t.Category != nil {
			l.CategoryID, l.Category = t.Category.ID, t.Category.Title
		}
		if err := enc.Encode(l); err != nil {
			return err
		}
	}
	return nil
}
