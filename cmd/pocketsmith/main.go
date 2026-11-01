// Command pocketsmith reads PocketSmith transactions for ad hoc use.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cehbz/pocketsmith"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "pocketsmith:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out io.Writer) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("pocketsmith", flag.ContinueOnError)
	configPath := fs.String("config", filepath.Join(home, ".config", "pocketsmith", "config.toml"), "config file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() == 0 || fs.Arg(0) != "transactions" {
		return errors.New("usage: pocketsmith [-config file] transactions [-from YYYY-MM-DD] [-to YYYY-MM-DD] [-account id] [-needs-review]")
	}
	cfg, err := loadConfig(*configPath)
	if err != nil {
		return err
	}
	key, err := readKey(ctx, "op", cfg.KeyRef)
	if err != nil {
		return err
	}
	c := &pocketsmith.Client{BaseURL: pocketsmith.DefaultBaseURL, Key: key, HTTP: &http.Client{Timeout: 30 * time.Second}}
	p, err := c.Profile(ctx)
	if err != nil {
		return err
	}
	y, m, d := time.Now().In(p.TimeZone).Date()
	q, err := parseTransactions(fs.Args()[1:], p.ID, time.Date(y, m, d, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	txs, err := c.Transactions(ctx, q)
	if err != nil {
		return err
	}
	return writeLines(out, txs)
}
