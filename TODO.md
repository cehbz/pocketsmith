# TODO

- A label containing a comma doesn't survive an update: `Update` sends a transaction's labels as one comma-joined string, so PocketSmith would read `a,b` as two labels. Read from update.go, not yet seen against PocketSmith.
