package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/cehbz/pocketsmith"
)

func TestWriteLinesPrintsOneObjectPerTransaction(t *testing.T) {
	d, _ := time.Parse(time.DateOnly, "2026-01-12")
	txs := []pocketsmith.Transaction{
		{ID: 1, Date: d, Amount: -2.5, Payee: "CAFE ONE", Category: &pocketsmith.Category{ID: 4, Title: "Eating Out"}, NeedsReview: true, Labels: []string{"x"}, Account: pocketsmith.Account{ID: 7, Name: "Card"}},
		{ID: 2, Date: d, Amount: -100, Payee: "TO SAVINGS", Account: pocketsmith.Account{ID: 8, Name: "Checking"}, Transfer: true},
	}
	var b bytes.Buffer
	if err := writeLines(&b, txs); err != nil {
		t.Fatal(err)
	}
	want := `{"id":1,"date":"2026-01-12","amount":-2.5,"payee":"CAFE ONE","category_id":4,"category":"Eating Out","needs_review":true,"labels":["x"],"account_id":7,"account":"Card","transfer":false}
{"id":2,"date":"2026-01-12","amount":-100,"payee":"TO SAVINGS","category_id":0,"category":"","needs_review":false,"labels":[],"account_id":8,"account":"Checking","transfer":true}
`
	if b.String() != want {
		t.Errorf("got\n%s\nwant\n%s", b.String(), want)
	}
}

func TestWriteLinesPrintsCharactersUnescaped(t *testing.T) {
	var b bytes.Buffer
	if err := writeLines(&b, []pocketsmith.Transaction{{ID: 1, Payee: "A&B <C>"}}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(b.Bytes(), []byte(`"payee":"A&B <C>"`)) {
		t.Errorf("got %s, want the payee unescaped", b.String())
	}
}

func TestWriteLinesPrintsEmptyLabels(t *testing.T) {
	var b bytes.Buffer
	if err := writeLines(&b, []pocketsmith.Transaction{{ID: 1}}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(b.Bytes(), []byte(`"labels":[]`)) {
		t.Errorf("got %s, want \"labels\":[]", b.String())
	}
}

func TestWriteLinesPrintsNothingForNoTransactions(t *testing.T) {
	var b bytes.Buffer
	if err := writeLines(&b, nil); err != nil || b.Len() != 0 {
		t.Errorf("writeLines(nil) wrote %q, %v", b.String(), err)
	}
}
