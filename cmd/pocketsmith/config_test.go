package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, name, content string, mode os.FileMode) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadConfigReadsTheKeyRef(t *testing.T) {
	c, err := loadConfig(write(t, "config.toml", "key_ref = \"op://Vault/Item/field\"\n", 0o600))
	if err != nil || c.KeyRef != "op://Vault/Item/field" {
		t.Errorf("loadConfig = %+v, %v", c, err)
	}
}

func TestLoadConfigRejectsMissingAndUnknownKeys(t *testing.T) {
	for name, content := range map[string]string{
		"missing": "",
		"unknown": "key_ref = \"op://V/I/f\"\nkey_rf = \"x\"\n",
	} {
		if _, err := loadConfig(write(t, "config.toml", content, 0o600)); err == nil {
			t.Errorf("%s: loadConfig succeeded", name)
		}
	}
}

func TestReadKeyRunsOpRead(t *testing.T) {
	op := write(t, "op", "#!/bin/sh\n[ \"$1\" = read ] && [ \"$2\" = \"op://Vault/Item/field\" ] && echo '  secret-value  '\n", 0o755)
	key, err := readKey(context.Background(), op, "op://Vault/Item/field")
	if err != nil || key != "secret-value" {
		t.Errorf("readKey = %q, %v", key, err)
	}
}

func TestReadKeyFailsOnEmptyOrFailedOutput(t *testing.T) {
	for name, script := range map[string]string{
		"empty":  "#!/bin/sh\necho\n",
		"failed": "#!/bin/sh\nexit 1\n",
	} {
		op := write(t, "op", script, 0o755)
		if _, err := readKey(context.Background(), op, "op://Vault/Item/field"); err == nil {
			t.Errorf("%s: readKey succeeded", name)
		}
	}
}

func TestReadKeyReportsOpStderrButNotStdout(t *testing.T) {
	op := write(t, "op", "#!/bin/sh\necho partial-secret\necho \"not signed in\" >&2\nexit 1\n", 0o755)
	_, err := readKey(context.Background(), op, "op://Vault/Item/field")
	if want := "op read op://Vault/Item/field: exit status 1: not signed in"; err == nil || err.Error() != want {
		t.Errorf("error = %v, want %q", err, want)
	}
}
