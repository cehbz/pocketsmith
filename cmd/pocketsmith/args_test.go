package main

import (
	"testing"
	"time"
)

func day(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse(time.DateOnly, s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestParseTransactionsDefaultsToTheLastNinetyDays(t *testing.T) {
	q, err := parseTransactions(nil, 42, day(t, "2026-03-31"))
	if err != nil {
		t.Fatal(err)
	}
	if q.User != 42 || q.Account != 0 || !q.From.Equal(day(t, "2026-01-01")) || !q.To.Equal(day(t, "2026-03-31")) || q.NeedsReview {
		t.Errorf("query = %+v, want user 42, 2026-01-01..2026-03-31", q)
	}
}

func TestParseTransactionsReadsEveryFlag(t *testing.T) {
	q, err := parseTransactions([]string{"-from", "2026-02-01", "-to", "2026-02-28", "-account", "7", "-needs-review"}, 42, day(t, "2026-03-31"))
	if err != nil {
		t.Fatal(err)
	}
	if q.User != 0 || q.Account != 7 || !q.From.Equal(day(t, "2026-02-01")) || !q.To.Equal(day(t, "2026-02-28")) || !q.NeedsReview {
		t.Errorf("query = %+v", q)
	}
}

func TestParseTransactionsRejectsBadInput(t *testing.T) {
	for _, args := range [][]string{
		{"-from", "02/01/2026"},
		{"-to", "tomorrow"},
		{"-from", "2026-03-01", "-to", "2026-02-01"},
		{"extra"},
	} {
		if _, err := parseTransactions(args, 42, day(t, "2026-03-31")); err == nil {
			t.Errorf("parseTransactions(%q) succeeded", args)
		}
	}
}
