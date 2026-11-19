# pocketsmith

A Go client and command for the parts of PocketSmith's v2 API I use.

## What

A Go package for the parts of [PocketSmith's v2 API](https://developers.pocketsmith.com/) its users need: the user's profile and time zone, listing categories, listing and reading transactions, and updating a transaction's category, review flag and labels. A small command, `pocketsmith transactions`, prints a range of transactions as JSON lines.

## Why

I categorize my PocketSmith transactions with [pocketsmith-categorizer](https://github.com/cehbz/pocketsmith-categorizer), and wanted it and my ad hoc tools to share one client. The command gives a Claude Code session my transactions without the PocketSmith MCP connector.

## How

The package is the standard library's HTTP client over PocketSmith's v2 endpoints: a `Client` holds the base URL and developer key, each method calls one or two endpoints, and a listing follows PocketSmith's `Link` header through every page. The command in `cmd/pocketsmith` reads its TOML config with [BurntSushi/toml](https://github.com/BurntSushi/toml) and the key with `op read`. Tests run the client against a local `httptest` server, so they never reach PocketSmith. Releases are tagged `v0.x` for consumers to pin. The domain model is in [DESIGN.md](DESIGN.md).

## How it differs

I know of no other client for PocketSmith's API. This one was written from scratch against the API.

## What it doesn't do

It covers only the calls above: no creating or deleting transactions, and nothing for budgets, institutions, attachments or other parts of the API. A category listing is flat and drops the parent relationship. The command only reads, and it doesn't print a transaction's closing balance.

## More

- Library: `go get github.com/cehbz/pocketsmith`; API reference: `go doc -all github.com/cehbz/pocketsmith`
- Command: [docs/command.md](docs/command.md)
- Configuration template: [config.example.toml](config.example.toml)
- Tests: `go test ./...`
