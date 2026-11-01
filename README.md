# pocketsmith

A Go client for the parts of [PocketSmith's v2 API](https://developers.pocketsmith.com/) its users need: the user's profile and time zone, listing categories, listing and reading transactions, and updating a transaction's category, review flag and labels.

## Command

```sh
go install github.com/cehbz/pocketsmith/cmd/pocketsmith@latest
pocketsmith transactions [-from YYYY-MM-DD] [-to YYYY-MM-DD] [-account id] [-needs-review]
```

It prints one JSON object per transaction per line: `id`, `date`, `amount` (a debit is negative), `payee`, `category_id`, `category`, `needs_review`, `labels`, `account_id`, `account`, `transfer`. The default range is the 90 days through today in the user's PocketSmith time zone, across every account.

The developer key is read through the 1Password CLI from the reference in `~/.config/pocketsmith/config.toml` (template: `config.example.toml`).
