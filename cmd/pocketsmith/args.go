package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/cehbz/pocketsmith"
)

// parseTransactions reads the transactions command's flags into a query. today is the
// current date in the user's PocketSmith time zone, as midnight UTC; user is the profile ID.
func parseTransactions(args []string, user int64, today time.Time) (pocketsmith.Query, error) {
	fs := flag.NewFlagSet("transactions", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	from := fs.String("from", "", "first date, YYYY-MM-DD (default: 89 days before -to)")
	to := fs.String("to", "", "last date, YYYY-MM-DD (default: today)")
	account := fs.Int64("account", 0, "transaction account ID (default: every account)")
	needsReview := fs.Bool("needs-review", false, "only transactions needing review")
	if err := fs.Parse(args); err != nil {
		return pocketsmith.Query{}, err
	}
	if fs.NArg() > 0 {
		return pocketsmith.Query{}, fmt.Errorf("transactions: unexpected argument %q", fs.Arg(0))
	}
	q := pocketsmith.Query{Account: *account, To: today, NeedsReview: *needsReview}
	if *account == 0 {
		q.User = user
	}
	var err error
	if *to != "" {
		if q.To, err = time.Parse(time.DateOnly, *to); err != nil {
			return pocketsmith.Query{}, fmt.Errorf("transactions -to: %w", err)
		}
	}
	q.From = q.To.AddDate(0, 0, -89)
	if *from != "" {
		if q.From, err = time.Parse(time.DateOnly, *from); err != nil {
			return pocketsmith.Query{}, fmt.Errorf("transactions -from: %w", err)
		}
	}
	if q.From.After(q.To) {
		return pocketsmith.Query{}, errors.New("transactions: -from is after -to")
	}
	return q, nil
}
